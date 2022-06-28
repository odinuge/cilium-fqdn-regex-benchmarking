# cilium-fqdn-regex-benchmarking

Benchmarking code used for comparing various techniques for improving the dnsproxy performance in cilium.

More info here; https://docs.cilium.io/en/v1.10/gettingstarted/dns/


### Displaying the different implementations
```
$ go run main.go
FQDN patterns:
[MatchName: uged.al, MatchPattern:  MatchName: cilium.io, MatchPattern:  MatchName: , MatchPattern: wil*dc.ard MatchName: , MatchPattern: *.s3.io]

Regexes:
baseline                 : "(^uged[.]al[.]$)|(^cilium[.]io[.]$)|(^wil[-a-zA-Z0-9_]*dc[.]ard[.]$)|(^[-a-zA-Z0-9_]*[.]s3[.]io[.]$)"  +  map[]
reverse                  : "^[.](?:la[.]degu|oi[.]muilic|dra[.]cd[-a-zA-Z0-9_]*liw|oi[.]3s[.][-a-zA-Z0-9_]*)$"  +  map[]
reverse+sort             : "^[.](?:dra[.]cd[-a-zA-Z0-9_]*liw|la[.]degu|oi[.]3s[.][-a-zA-Z0-9_]*|oi[.]muilic)$"  +  map[]
map+baseline             : "(^wil[-a-zA-Z0-9_]*dc[.]ard[.]$)|(^[-a-zA-Z0-9_]*[.]s3[.]io[.]$)"  +  map[cilium.io.: uged.al.:]
map+optimized            : "^(?:wil[-a-zA-Z0-9_]*dc[.]ard|[-a-zA-Z0-9_]*[.]s3[.]io)[.]$"  +  map[cilium.io.: uged.al.:]
map+reverse+no+onepass   : "^[.](?:dra[.]cd[-a-zA-Z0-9_]*liw|oi[.]3s[.][-a-zA-Z0-9_]*)($)"  +  map[.la.degu: .oi.muilic:]
map+reverse              : "^[.](?:dra[.]cd[-a-zA-Z0-9_]*liw|oi[.]3s[.][-a-zA-Z0-9_]*)$"  +  map[.la.degu: .oi.muilic:]
```


### Running the benchmarks

Tests can be configured inside the `benchmark` folder, together with the proposed techniques.

```
$ go test ./benchmarks -bench=. -test.cpu=1 -benchtime=1x
goos: darwin
goarch: amd64
pkg: github.com/odinuge/cilium-fqdn-regex-benchmarking/benchmarks
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkAllTechniques/generate_________:_baseline_______________:_domains:3000                 1           1291227 ns/op          956616 B/op      13216 allocs/op
BenchmarkAllTechniques/regex-compile____:_baseline_______________:_domains:3000                 1          54914181 ns/op        34159856 B/op      53441 allocs/op
BenchmarkAllTechniques/match-fqdn_______:_baseline_______________:_domains:3000                 1        20954871024 ns/op       197631864 B/op      7842 allocs/op
BenchmarkAllTechniques/match-wild_______:_baseline_______________:_domains:3000                 1        14120048832 ns/op            136 B/op          2 allocs/op
BenchmarkAllTechniques/match-not-wild___:_baseline_______________:_domains:3000                 1        19075265535 ns/op            136 B/op          2 allocs/op
BenchmarkAllTechniques/match-not-tld____:_baseline_______________:_domains:3000                 1        24439911127 ns/op            136 B/op          2 allocs/op
BenchmarkAllTechniques/generate_________:_reverse________________:_domains:3000                 1           2826143 ns/op          842992 B/op      10217 allocs/op
BenchmarkAllTechniques/regex-compile____:_reverse________________:_domains:3000                 1          64226635 ns/op        21959872 B/op      40421 allocs/op
BenchmarkAllTechniques/match-fqdn_______:_reverse________________:_domains:3000                 1        2887679397 ns/op         2498512 B/op       3475 allocs/op
BenchmarkAllTechniques/match-wild_______:_reverse________________:_domains:3000                 1        1029416459 ns/op             136 B/op          2 allocs/op
BenchmarkAllTechniques/match-not-wild___:_reverse________________:_domains:3000                 1        1355863700 ns/op             136 B/op          2 allocs/op
BenchmarkAllTechniques/match-not-tld____:_reverse________________:_domains:3000                 1           4738081 ns/op             136 B/op          2 allocs/op
BenchmarkAllTechniques/generate_________:_reverse+sort___________:_domains:3000                 1           1912227 ns/op          843016 B/op      10218 allocs/op
BenchmarkAllTechniques/regex-compile____:_reverse+sort___________:_domains:3000                 1          22415234 ns/op        13373216 B/op      40942 allocs/op
BenchmarkAllTechniques/match-fqdn_______:_reverse+sort___________:_domains:3000                 1          10957081 ns/op         1429168 B/op        113 allocs/op
BenchmarkAllTechniques/match-wild_______:_reverse+sort___________:_domains:3000                 1         588660980 ns/op           45144 B/op       1300 allocs/op
BenchmarkAllTechniques/match-not-wild___:_reverse+sort___________:_domains:3000                 1         210428783 ns/op             136 B/op          2 allocs/op
BenchmarkAllTechniques/match-not-tld____:_reverse+sort___________:_domains:3000                 1           1841122 ns/op             136 B/op          2 allocs/op
BenchmarkAllTechniques/generate_________:_map+baseline___________:_domains:3000                 1            956793 ns/op          662944 B/op       7857 allocs/op
BenchmarkAllTechniques/regex-compile____:_map+baseline___________:_domains:3000                 1          17004102 ns/op         9514880 B/op      26302 allocs/op
BenchmarkAllTechniques/match-fqdn_______:_map+baseline___________:_domains:3000                 1            348499 ns/op               0 B/op          0 allocs/op
BenchmarkAllTechniques/match-wild_______:_map+baseline___________:_domains:3000                 1        6449178052 ns/op        45061528 B/op       4266 allocs/op
BenchmarkAllTechniques/match-not-wild___:_map+baseline___________:_domains:3000                 1        9894274693 ns/op             136 B/op          2 allocs/op
BenchmarkAllTechniques/match-not-tld____:_map+baseline___________:_domains:3000                 1        9083675351 ns/op             136 B/op          2 allocs/op
BenchmarkAllTechniques/generate_________:_map+optimized__________:_domains:3000                 1            834319 ns/op          572720 B/op       5432 allocs/op
BenchmarkAllTechniques/regex-compile____:_map+optimized__________:_domains:3000                 1          12419005 ns/op        10509232 B/op      20166 allocs/op
BenchmarkAllTechniques/match-fqdn_______:_map+optimized__________:_domains:3000                 1            311109 ns/op               0 B/op          0 allocs/op
BenchmarkAllTechniques/match-wild_______:_map+optimized__________:_domains:3000                 1        2624224718 ns/op         1390688 B/op       4250 allocs/op
BenchmarkAllTechniques/match-not-wild___:_map+optimized__________:_domains:3000                 1        1034604871 ns/op             136 B/op          2 allocs/op
BenchmarkAllTechniques/match-not-tld____:_map+optimized__________:_domains:3000                 1        2615534246 ns/op             136 B/op          2 allocs/op
BenchmarkAllTechniques/generate_________:_map+reverse+no+onepass_:_domains:3000                 1           1425738 ns/op          797424 B/op      10220 allocs/op
BenchmarkAllTechniques/regex-compile____:_map+reverse+no+onepass_:_domains:3000                 1           7641553 ns/op         5194592 B/op      22744 allocs/op
BenchmarkAllTechniques/match-fqdn_______:_map+reverse+no+onepass_:_domains:3000                 1            305600 ns/op               0 B/op          0 allocs/op
BenchmarkAllTechniques/match-wild_______:_map+reverse+no+onepass_:_domains:3000                 1         573638698 ns/op          764432 B/op       1411 allocs/op
BenchmarkAllTechniques/match-not-wild___:_map+reverse+no+onepass_:_domains:3000                 1         198792586 ns/op             136 B/op          2 allocs/op
BenchmarkAllTechniques/match-not-tld____:_map+reverse+no+onepass_:_domains:3000                 1           1546179 ns/op             136 B/op          2 allocs/op
BenchmarkAllTechniques/generate_________:_map+reverse____________:_domains:3000                 1           1880893 ns/op          797424 B/op      10220 allocs/op
BenchmarkAllTechniques/regex-compile____:_map+reverse____________:_domains:3000                 1           8466408 ns/op         6285632 B/op      22756 allocs/op
BenchmarkAllTechniques/match-fqdn_______:_map+reverse____________:_domains:3000                 1            953713 ns/op               0 B/op          0 allocs/op
BenchmarkAllTechniques/match-wild_______:_map+reverse____________:_domains:3000                 1         599137110 ns/op          753280 B/op       1411 allocs/op
BenchmarkAllTechniques/match-not-wild___:_map+reverse____________:_domains:3000                 1         205130379 ns/op             136 B/op          2 allocs/op
BenchmarkAllTechniques/match-not-tld____:_map+reverse____________:_domains:3000                 1           1447156 ns/op             136 B/op          2 allocs/op
PASS
ok      github.com/odinuge/cilium-fqdn-regex-benchmarking/benchmarks    119.019s
```
