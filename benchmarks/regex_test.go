package benchmarks

import (
	"bufio"
	"fmt"
	"github.com/cilium/cilium/pkg/policy/api"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

// Number of domains from list to add to regexp/matching
// zero means all
const numberOfDomainsInCNP = 500

// Number of domains to run checks against.
// zero means all domains.
// Keeping this number static when changing the numberOfDomainsInCNP
// is useful for comparing performance impact of the CNP size; since number of
// domain in matching tests will stay the same
const domainCheckSize = numberOfDomainsInCNP

var filename = "./10000-domains.txt"

// If you want to test a single technique to compare with eg. benchstat
func BenchmarkSingle(b *testing.B) {
	b.Skip()
	benchmarkTechniques(b, []ProposedTechnique{MapAndBaselineTechnique})
}
func BenchmarkAll(b *testing.B) {
	benchmarkTechniques(b, AllProposedTechniques)
}
func reverse[S ~[]E, E any](s S) S {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
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
			{"match-name", true, domainsFqdn},
			// Domains allowlisted by a matchPattern with a wildcard
			{"match-pattern", true, reverse(domainsWildcard)},
			// Domains not allowlisted, has a prefix found in no domain
			{"bad-match-pattern", false, badDomains},
			// Domains not allowlisted, has a tld found in no domain
			{"bad-match-tld", false, badDomainsTld},
		}
		testName := func(name string) string {
			var techniqueName string
			if len(techniques) != 1 {
				techniqueName = bm.Name
			}
			return fmt.Sprintf("%-17s: %-23s", name, techniqueName)
		}
		// the generate bench is for benchmarking the generation of the regex string, plus
		// the extra map in case its used
		b.Run(testName("generate"), func(b *testing.B) {
			b.StopTimer()
			b.ResetTimer()
			runtime.GC()
			var sum float64
			b.ReportAllocs()
			for a := 0; a < b.N; a++ {
				runtime.GC()
				runtime.GC()
				s := getMemStats()
				b.StartTimer()
				fqs, regex := bm.RegexGenerator(selector)
				b.StopTimer()
				runtime.GC()
				runtime.GC()
				s2 := getMemStats()
				sum += float64(s2.HeapInuse - s.HeapInuse)
				if regex == "" || (fqs != nil && len(fqs) == 0) {
					b.Fail()
				}
				if !(strings.Contains(regex, "testing-domain") || strings.Contains(regex, Reverse("testing-domain"))) {
					b.Fail()
				}
				for k := range fqs {
					delete(fqs, k)
				}
			}
			b.ReportMetric(sum/float64(b.N), "B(heap)/op")
			b.ReportMetric(float64(len(domains)), "domains")
		})
		b.Run(testName("compile"), func(b *testing.B) {
			b.StopTimer()
			_, regex := bm.RegexGenerator(selector)
			runtime.GC()
			s := getMemStats()
			var sum float64
			b.StartTimer()
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				pattern := regexp.MustCompile(regex)
				b.StopTimer()
				runtime.GC()
				runtime.GC()
				s2 := getMemStats()
				sum += float64(s2.HeapInuse - s.HeapInuse)
				if pattern.String() == "" || pattern.MatchString("asd") || pattern.MatchString(badDomains[0]) || !pattern.MatchString(domainsWildcard[0]) {
					b.Fail()
				}
				b.StartTimer()
			}
			b.StopTimer()
			b.ReportMetric(sum/float64(b.N), "B(heap)/op")
			b.ReportMetric(float64(len(domains)), "domains")
		})

		b.Run(testName("combined-compile"), func(b *testing.B) {
			b.StopTimer()
			runtime.GC()
			s := getMemStats()
			var sum float64
			b.ResetTimer()
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				runtime.GC()
				runtime.GC()
				b.StartTimer()
				fqs, regex := bm.RegexGenerator(selector)
				pattern := regexp.MustCompile(regex)
				b.StopTimer()
				runtime.GC()
				runtime.GC()
				s2 := getMemStats()
				sum += float64(s2.HeapInuse - s.HeapInuse)
				if pattern.String() == "" || pattern.MatchString("asd") || pattern.MatchString(badDomains[0]) || !pattern.MatchString(domainsWildcard[0]) {
					b.Fail()
				}
				for k := range fqs {
					delete(fqs, k)
				}
			}
			b.StopTimer()
			b.ReportMetric(sum/float64(b.N), "B(heap)/op")
			b.ReportMetric(float64(len(domains)), "domains")
		})
		mapping, regex := bm.RegexGenerator(selector)
		pattern := regexp.MustCompile(regex)

		numberOfDomains := domainCheckSize
		if numberOfDomains == 0 {
			numberOfDomains = len(domains)
		}
		for _, matchingTest := range matchingTests {
			b.Run(testName(matchingTest.name), func(b *testing.B) {
				b.ReportAllocs()
				b.StopTimer()
				b.ResetTimer()
				for n := 0; n < b.N; n++ {
					runtime.GC()
					runtime.GC()
					runtime.GC()
					runtime.GC()
					b.StartTimer()
					var isInFqdnMap bool
					testDomain := matchingTest.inputDomains[n%len(matchingTest.inputDomains)]
					if mapping != nil {
						_, isInFqdnMap = mapping[testDomain]
					}
					if (isInFqdnMap || pattern.MatchString(testDomain)) != matchingTest.matchExpected {
						fmt.Println(testDomain)
						b.Fail()
					}
					b.StopTimer()
				}
				b.ReportMetric(float64(len(domains)), "domains")
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
	lines = append([]string{"*.testing-domains.aws.baa"}, lines...)
	for i, line := range lines {
		if dataSize != 0 && i >= dataSize {
			break
		}
		domains = append(domains,
			domainRewriter(strings.ReplaceAll(line, "*", "bar")))
		badDomains = append(badDomains,
			domainRewriter("nope.nop.nop.nop."+strings.ReplaceAll(line, "*", "bar.hmm")))
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
func getMemStats() runtime.MemStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m
}
