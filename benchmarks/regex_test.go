package benchmarks

import (
	"bufio"
	"fmt"
	"github.com/cilium/cilium/pkg/policy/api"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// Number of domains from list to add to regexp/matching
// zero means all
const numberOfDomainsInCNP = 3000

// Number of domains to run checks against.
// zero means all domains.
// Keeping this number static when changing the numberOfDomainsInCNP
// is useful for comparing performance impact of the CNP size; since number of
// domain in matching tests will stay the same
const domainCheckSize = numberOfDomainsInCNP

var filename = "./10000-domains.txt"

// If you want to test a single technique to compare with eg. benchstat
func BenchmarkSingleTechnique(b *testing.B) {
	b.Skip()
	benchmarkTechniques(b, []ProposedTechnique{MapAndReverseTechniqueNoOnepass})
}
func BenchmarkAllTechnique(b *testing.B) {
	benchmarkTechniques(b, AllProposedTechniques)
}

func benchmarkTechniques(b *testing.B, techniques []ProposedTechnique) {
	for _, bm := range techniques {
		selector, domains, domainsFqdn, domainsWildcard, badDomains, badDomainsTld := getTestData(filename, numberOfDomainsInCNP, bm.DomainRewriter)
		matchingTests := []struct {
			name          string
			matchExpected bool
			inputDomains  []string
		}{
			// Domains allowlisted by a matchName
			{"match-fqdn", true, domainsFqdn},
			// Domains allowlisted by a matchPattern with a wildcard
			{"match-wild", true, domainsWildcard},
			// Domains not allowlisted, has a prefix found in no domain
			{"match-not-wild", false, badDomains},
			// Domains not allowlisted, has a tld found in no domain
			{"match-not-tld", false, badDomainsTld},
		}
		testName := func(name string) string {
			var techniqueName string
			if len(techniques) != 1 {
				techniqueName = bm.Name
			}
			return fmt.Sprintf("%-17s: %-23s: domains:%4d", name, techniqueName, len(domains))
		}
		// the generate bench is for benchmarking the generation of the regex string, plus
		// the extra map in case its used
		b.Run(testName("generate"), func(b *testing.B) {
			b.ReportAllocs()
			for a := 0; a < b.N; a++ {
				_, regex := bm.RegexGenerator(selector)
				if regex == "" {
					b.Fail()
				}
			}
		})
		_, regex := bm.RegexGenerator(selector)
		b.Run(testName("regex-compile"), func(b *testing.B) {
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				pattern := regexp.MustCompile(regex)
				if pattern.String() == "" {
					b.Fail()
				}
			}

		})
		var mapping map[string]string
		mapping, regex = bm.RegexGenerator(selector)
		pattern := regexp.MustCompile(regex)

		numberOfDomains := domainCheckSize
		if numberOfDomains == 0 {
			numberOfDomains = len(domains)
		}
		for _, matchingTest := range matchingTests {
			b.Run(testName(matchingTest.name), func(b *testing.B) {
				b.ReportAllocs()
				for n := 0; n < b.N; n++ {
					for i := 0; i < numberOfDomains; i++ {
						testDomain := matchingTest.inputDomains[i%len(matchingTest.inputDomains)]
						_, isInFqdnMap := mapping[testDomain]
						if (isInFqdnMap || pattern.MatchString(testDomain)) != matchingTest.matchExpected {
							fmt.Println(testDomain)
							b.Fail()
						}
					}
				}
			})
		}
	}
}

func getTestData(filename string, dataSize int, domainRewriter func(string) string) (selector []*api.FQDNSelector, domains []string, domainsFqdn []string, domainsWildcard []string, badDomains []string, badDomainsTld []string) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	// Always append a wildcard domain like this to ensure we have one, since its required for the test
	lines = append([]string{"*.sbsdasdasd3.aws.baa"}, lines...)
	for i, line := range lines {
		if dataSize != 0 && i >= dataSize {
			break
		}
		domains = append(domains,
			domainRewriter(strings.ReplaceAll(line, "*", "bar")))
		badDomains = append(badDomains,
			domainRewriter("nope.nop.nop.nop."+strings.ReplaceAll(line, "*", "bar")))
		badDomainsTld = append(badDomainsTld,
			domainRewriter(strings.ReplaceAll(line, "*", "bar")+".bar"+strconv.Itoa(i)+"s"))

		if strings.Contains(line, "*") {
			domainsWildcard = append(domainsWildcard, domainRewriter(strings.ReplaceAll(line, "*", "abc"+strconv.Itoa(i)+"c")))
		} else {
			domainsFqdn = append(domainsFqdn, domainRewriter(line))
		}

		if strings.Contains(line, "*") {
			selector = append(selector, &api.FQDNSelector{
				MatchPattern: strings.TrimSpace(line),
			})
		} else {
			selector = append(selector, &api.FQDNSelector{
				MatchName: strings.TrimSpace(line),
			})
		}
	}
	return
}
