package main

import (
	"fmt"
	v2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
	"github.com/cilium/cilium/pkg/policy/api"
	"io/ioutil"
	"log"
	"os"
	"sigs.k8s.io/yaml"
	"sort"
)

// Add your CNP as the firs arg
// $ go run utils/sort-cnps-by-size/main.go my-cnps.yml
// It will then sort them in ascending order by number of dns rules
func main() {
	bb, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	var cnpList v2.CiliumNetworkPolicyList
	err = yaml.Unmarshal(bb, &cnpList)

	if err != nil {
		log.Fatal(err)
	}

	countDNSRules := func(e []api.EgressRule) int {
		v := 0
		for _, egress := range e {
			if egress.ToPorts != nil {
				for _, ports := range egress.ToPorts {
					if ports.Rules != nil {
						v += len(ports.Rules.DNS)
					}
				}
			}
		}
		return v
	}

	var foo []struct {
		cnp   v2.CiliumNetworkPolicy
		count int
	}
	for _, cnp := range cnpList.Items {
		if cnp.Specs != nil {
			for _, spec := range cnp.Specs {
				if spec.Egress != nil {
					foo = append(foo, struct {
						cnp   v2.CiliumNetworkPolicy
						count int
					}{
						cnp:   cnp,
						count: countDNSRules(spec.Egress),
					})
					//TODO figure out this one
					break
				}
			}
		} else if cnp.Spec != nil {
			if cnp.Spec.Egress != nil {
				foo = append(foo, struct {
					cnp   v2.CiliumNetworkPolicy
					count int
				}{
					cnp:   cnp,
					count: countDNSRules(cnp.Spec.Egress),
				})
			}
		}
	}
	sort.Slice(foo, func(i, j int) bool {
		return foo[i].count < foo[j].count
	})

	cn := cnpList.DeepCopy()
	cn.Items = []v2.CiliumNetworkPolicy{}
	for _, v := range foo {
		cn.Items = append(cn.Items, v.cnp)
	}
	output, err := yaml.Marshal(cn)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(output))

}
