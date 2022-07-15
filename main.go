package main

import (
	"fmt"
	"github.com/cilium/cilium/pkg/policy/api"
	"github.com/odinuge/cilium-fqdn-regex-benchmarking/benchmarks"
)

func main() {

	slect := []*api.FQDNSelector{
		{MatchName: "uged.al"},
		{MatchName: "cilium.io"},
		{MatchPattern: "wil*dc.ard"},
		{MatchPattern: "*.s3.io"},
	}
	fmt.Println("FQDN patterns:")
	fmt.Printf("%+v\n\n", slect)
	fmt.Println("Regexes:")
	for _, sol := range benchmarks.AllProposedTechniques {
		maap, rei := sol.RegexGenerator(slect)
		fmt.Printf("%-25s: %q  +  %+v\n", sol.Name, rei, maap)
	}
}
