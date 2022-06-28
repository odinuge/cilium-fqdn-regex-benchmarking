package benchmarks

import (
	"github.com/cilium/cilium/pkg/fqdn/matchpattern"
	"github.com/cilium/cilium/pkg/policy/api"
	"github.com/miekg/dns"
	"sort"
	"strings"
	"unicode/utf8"
)

func Reverse(s string) string {
	size := len(s)
	buf := make([]byte, size)
	for start := 0; start < size; {
		r, n := utf8.DecodeRuneInString(s[start:])
		start += n
		utf8.EncodeRune(buf[size-start:], r)
	}
	return string(buf)
}

func Baseline(DNSRules []*api.FQDNSelector) (map[string]string, string) {
	reStrings := make([]string, 0, len(DNSRules))
	for _, dnsRule := range DNSRules {
		if len(dnsRule.MatchName) > 0 {
			dnsRuleName := strings.ToLower(dns.Fqdn(dnsRule.MatchName))
			dnsPatternAsRE := matchpattern.ToRegexp(dnsRuleName)
			reStrings = append(reStrings, "("+dnsPatternAsRE+")")
		}
		if len(dnsRule.MatchPattern) > 0 {
			dnsPattern := matchpattern.Sanitize(dnsRule.MatchPattern)
			dnsPatternAsRE := matchpattern.ToRegexp(dnsPattern)
			reStrings = append(reStrings, "("+dnsPatternAsRE+")")
		}
	}
	return make(map[string]string), strings.Join(reStrings, "|")
}

func ReverseAndSortRegex(DNSRules []*api.FQDNSelector) (map[string]string, string) {
	reStrings := make([]string, 0, len(DNSRules))
	for _, dnsRule := range DNSRules {
		if len(dnsRule.MatchName) > 0 {
			dnsRuleName := FromFqdn(dnsRule.MatchName)
			dnsPatternAsRE := CustomToRegexp(Reverse(dnsRuleName))
			reStrings = append(reStrings, dnsPatternAsRE)
		}
		if len(dnsRule.MatchPattern) > 0 {
			dnsPattern := FromFqdn(dnsRule.MatchPattern)
			dnsPatternAsRE := CustomToRegexp(Reverse(dnsPattern))
			reStrings = append(reStrings, dnsPatternAsRE)
		}
	}
	sort.Strings(reStrings)
	return make(map[string]string), "^[.](?:" + strings.Join(reStrings, "|") + ")$"
}

func MapAndBaseline(DNSRules []*api.FQDNSelector) (map[string]string, string) {
	reStrings := make([]string, 0, len(DNSRules))
	nice := make(map[string]string, len(DNSRules))
	for _, dnsRule := range DNSRules {
		if len(dnsRule.MatchName) > 0 {
			dnsRuleName := strings.ToLower(dns.Fqdn(dnsRule.MatchName))
			nice[dnsRuleName] = ""
		}
		if len(dnsRule.MatchPattern) > 0 {
			dnsPattern := matchpattern.Sanitize(dnsRule.MatchPattern)
			dnsPatternAsRE := matchpattern.ToRegexp(dnsPattern)
			reStrings = append(reStrings, "("+dnsPatternAsRE+")")
		}
	}
	return nice, strings.Join(reStrings, "|")
}

func MapAndOptimized(DNSRules []*api.FQDNSelector) (map[string]string, string) {
	reStrings := make([]string, 0, len(DNSRules))
	nice := make(map[string]string, len(DNSRules))
	for _, dnsRule := range DNSRules {
		if len(dnsRule.MatchName) > 0 {
			dnsRuleName := strings.ToLower(dns.Fqdn(dnsRule.MatchName))
			nice[dnsRuleName] = ""
		}
		if len(dnsRule.MatchPattern) > 0 {
			dnsPattern := matchpattern.Sanitize(dnsRule.MatchPattern)
			dnsPatternAsRE := CustomToRegexp(FromFqdn(dnsPattern))
			reStrings = append(reStrings, dnsPatternAsRE)
		}
	}
	return nice, "^(?:" + strings.Join(reStrings, "|") + ")[.]$"
}

func MapAndRegexNoOnepass(DNSRules []*api.FQDNSelector) (map[string]string, string) {
	reStrings := make([]string, 0, len(DNSRules))
	nice := make(map[string]string, len(DNSRules))
	for _, dnsRule := range DNSRules {
		if len(dnsRule.MatchName) > 0 {
			dnsRuleName := dns.Fqdn(dnsRule.MatchName)
			nice[Reverse(dnsRuleName)] = ""
		}
		if len(dnsRule.MatchPattern) > 0 {
			dnsPattern := FromFqdn(dnsRule.MatchPattern)
			dnsPatternAsRE := CustomToRegexp(Reverse(dnsPattern))
			reStrings = append(reStrings, dnsPatternAsRE)
		}
	}
	sort.Strings(reStrings)
	return nice, "^[.](?:" + strings.Join(reStrings, "|") + ")($)"
}
func MapAndRegex(DNSRules []*api.FQDNSelector) (map[string]string, string) {
	reStrings := make([]string, 0, len(DNSRules))
	nice := make(map[string]string, len(DNSRules))
	for _, dnsRule := range DNSRules {
		if len(dnsRule.MatchName) > 0 {
			dnsRuleName := dns.Fqdn(dnsRule.MatchName)
			nice[Reverse(dnsRuleName)] = ""
		}
		if len(dnsRule.MatchPattern) > 0 {
			dnsPattern := FromFqdn(dnsRule.MatchPattern)
			dnsPatternAsRE := CustomToRegexp(Reverse(dnsPattern))
			reStrings = append(reStrings, dnsPatternAsRE)
		}
	}
	sort.Strings(reStrings)
	return nice, "^[.](?:" + strings.Join(reStrings, "|") + ")$"
}
func ReverseRegex(DNSRules []*api.FQDNSelector) (map[string]string, string) {
	reStrings := make([]string, 0, len(DNSRules))
	for _, dnsRule := range DNSRules {
		if len(dnsRule.MatchName) > 0 {
			dnsRuleName := FromFqdn(dnsRule.MatchName)
			dnsPatternAsRE := CustomToRegexp(Reverse(dnsRuleName))
			reStrings = append(reStrings, dnsPatternAsRE)
		}
		if len(dnsRule.MatchPattern) > 0 {
			dnsPattern := FromFqdn(dnsRule.MatchPattern)
			dnsPatternAsRE := CustomToRegexp(Reverse(dnsPattern))
			reStrings = append(reStrings, dnsPatternAsRE)
		}
	}
	return make(map[string]string), "^[.](?:" + strings.Join(reStrings, "|") + ")$"
}

const allowedDNSCharsREGroup = "[-a-zA-Z0-9_]"

func CustomToRegexp(pattern string) string {
	pattern = strings.TrimSpace(pattern)
	pattern = strings.ToLower(pattern)

	// handle the * match-all case. This will filter down to the end.
	if pattern == "*" {
		return "(^(" + allowedDNSCharsREGroup + "+[.])+$)|(^[.]$)"
	}

	// base case. * becomes .*, but only for DNS valid characters
	// NOTE: this only works because the case above does not leave the *
	pattern = strings.Replace(pattern, "*", allowedDNSCharsREGroup+"*", -1)

	// base case. "." becomes a literal .
	pattern = strings.Replace(pattern, ".", "[.]", -1)

	// Anchor the match to require the whole string to match this expression
	return pattern
}
func FromFqdn(s string) string {
	if !dns.IsFqdn(s) {
		return s
	}
	return s[:len(s)-1]
}

func identityFqdn(s string) string {
	return dns.Fqdn(s)
}
func reverseFqdn(s string) string {
	return Reverse(dns.Fqdn(s))
}

type ProposedTechnique struct {
	// The name of the technique
	Name string
	// The mapping function from list of fqdns/patterns to regexp and/or map of fqdns
	RegexGenerator func([]*api.FQDNSelector) (map[string]string, string)
	// Mapping from fqdn to the input to the regexp and/or map
	DomainRewriter func(string) string
}

var (
	BaselineTechnique        = ProposedTechnique{"baseline", Baseline, identityFqdn}
	ReverseTechnique         = ProposedTechnique{"reverse", ReverseRegex, reverseFqdn}
	ReverseAndSortTechnique  = ProposedTechnique{"reverse+sort", ReverseAndSortRegex, reverseFqdn}
	MapAndBaselineTechnique  = ProposedTechnique{"map+baseline", MapAndBaseline, identityFqdn}
	MapAndOptimizedTechnique = ProposedTechnique{"map+optimized", MapAndOptimized, identityFqdn}
	MapAndReverseTechnique   = ProposedTechnique{"map+reverse", MapAndRegex, reverseFqdn}
	// This should force the golang regex implementation to avoid trying to create a onepass program. Since we won't end up using it anyways,
	// and the compilation will "fail" anyways, we just short circuit to avoid the extra allocations when copying the program
	MapAndReverseTechniqueNoOnepass = ProposedTechnique{"map+reverse+no+onepass", MapAndRegexNoOnepass, reverseFqdn}

	AllProposedTechniques = []ProposedTechnique{
		BaselineTechnique,
		ReverseTechnique,
		ReverseAndSortTechnique,
		MapAndBaselineTechnique,
		MapAndOptimizedTechnique,
		MapAndReverseTechniqueNoOnepass,
		MapAndReverseTechnique,
	}
)
