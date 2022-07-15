# cilium-fqdn-regex-benchmarking

Benchmarking code used for comparing various techniques for improving the dnsproxy performance in cilium.

More info here; https://docs.cilium.io/en/v1.10/gettingstarted/dns/


### Displaying the different implementations
```
FQDN patterns:
[MatchName: uged.al, MatchPattern:  MatchName: cilium.io, MatchPattern:  MatchName: , MatchPattern: wil*dc.ard MatchName: , MatchPattern: *.s3.io]

Regexes:
baseline                 : "(^uged[.]al[.]$)|(^cilium[.]io[.]$)|(^wil[-a-zA-Z0-9_]*dc[.]ard[.]$)|(^[-a-zA-Z0-9_]*[.]s3[.]io[.]$)"  +  map[]
baseline-optimized       : "^(?:[-a-zA-Z0-9_]*[.]s3[.]io[.]|cilium[.]io[.]|uged[.]al[.]|wil[-a-zA-Z0-9_]*dc[.]ard[.])$"  +  map[]
reverse                  : "^[.](?:la[.]degu|oi[.]muilic|dra[.]cd[-a-zA-Z0-9_]*liw|oi[.]3s[.][-a-zA-Z0-9_]*)$"  +  map[]
reverse+sort             : "^[.](?:dra[.]cd[-a-zA-Z0-9_]*liw|la[.]degu|oi[.]3s[.][-a-zA-Z0-9_]*|oi[.]muilic)$"  +  map[]
map+baseline             : "(^wil[-a-zA-Z0-9_]*dc[.]ard[.]$)|(^[-a-zA-Z0-9_]*[.]s3[.]io[.]$)"  +  map[cilium.io.: uged.al.:]
map+optimized            : "^(?:wil[-a-zA-Z0-9_]*dc[.]ard|[-a-zA-Z0-9_]*[.]s3[.]io)[.]$"  +  map[cilium.io.: uged.al.:]
map+reverse+sort         : "^[.](?:dra[.]cd[-a-zA-Z0-9_]*liw|oi[.]3s[.][-a-zA-Z0-9_]*)$"  +  map[.la.degu: .oi.muilic:]
```

As you see above, we have 7 different "solutions":
- `baseline` - the currently implemented solution in cilium
- `baseline-optimized` - basline where we extract the regex anchors and stop using groups.
- `reverse` - where we also do as `baseline-optimized`, but also extracts the final dot in the fqdn, and reverses the domains when matching
- `reverse+sort` - where we do the same as `reverse`, just that after we reverse the matchNames/patterns, we sort them before creating the regex
- `map+baseline` - adds matchNames to a map and matchPatterns into the regex
- `map+optimized` - same as `map+baseline` plus the optimizations from `baseline-optimized`
- `map+reverse+sort` - All the optimizations

## Running the benchmarks

Tests can be configured inside the `benchmark` folder, together with the proposed techniques.

_note_: All benchmarks have been executed on my intel macbook. Will rerun on a linux machine with less noisy neighbours and reserved CPU cores.

### Generating the regex and compiling it

One very interesting thing is benchmarking the process of extracting the matchNames and matchPatterns, and then creating the regexp pattern (and possibly extracting the matchNames to a map),
and at the end compiling the regexp. This example is for 500 domains in the policy.

The `B(heap)/op` part is the increase in heap usage after the test iteration. In this test, that means the heap space used by the compiled regexp and the map with matchNames.

```
$ go test -bench="BenchmarkAll/combined-compile" -test.cpu=1 -benchtime=1000x ./benchmarks
goos: darwin
goarch: amd64
pkg: github.com/odinuge/cilium-fqdn-regex-benchmarking/benchmarks
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkAll/combined-compile_:_baseline_______________                     1000           4746479 ns/op           1261568 B(heap)/op           500.0 domains     5219801 B/op      11175 allocs/op
BenchmarkAll/combined-compile_:_baseline-optimized_____                     1000           8053903 ns/op           1048576 B(heap)/op           500.0 domains     6149577 B/op       8185 allocs/op
BenchmarkAll/combined-compile_:_reverse________________                     1000           3907888 ns/op            614400 B(heap)/op           500.0 domains     3375586 B/op       8542 allocs/op
BenchmarkAll/combined-compile_:_reverse+sort___________                     1000           2126606 ns/op            409600 B(heap)/op           500.0 domains     2317998 B/op       8561 allocs/op
BenchmarkAll/combined-compile_:_map+baseline___________                     1000           1161384 ns/op            385024 B(heap)/op           500.0 domains     1582837 B/op       5757 allocs/op
BenchmarkAll/combined-compile_:_map+optimized__________                     1000           1354995 ns/op            319488 B(heap)/op           500.0 domains     1509868 B/op       4324 allocs/op
BenchmarkAll/combined-compile_:_map+reverse+sort_______                     1000           1116912 ns/op            229376 B(heap)/op           500.0 domains     1073699 B/op       5532 allocs/op
PASS
ok      github.com/odinuge/cilium-fqdn-regex-benchmarking/benchmarks    74.920s
```

Its also possible to check the regex compile procedure isolated, as well as the "generate" pattern+map part;

```

# For compile
$ go test -bench="BenchmarkAll/^compile" -test.cpu=1 -benchtime=100x ./benchmarks
goos: darwin
goarch: amd64
pkg: github.com/odinuge/cilium-fqdn-regex-benchmarking/benchmarks
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkAll/compile__________:_baseline_______________                     1000           7164441 ns/op           1257357 B(heap)/op          500.0 domains     5060027 B/op       8971 allocs/op
BenchmarkAll/compile__________:_baseline-optimized_____                     1000           8174871 ns/op           1028047 B(heap)/op          500.0 domains     6025042 B/op       6979 allocs/op
BenchmarkAll/compile__________:_reverse________________                     1000           4206657 ns/op            602046 B(heap)/op          500.0 domains     3232215 B/op       6837 allocs/op
BenchmarkAll/compile__________:_reverse+sort___________                     1000           3605649 ns/op            393142 B(heap)/op          500.0 domains     2174604 B/op       6855 allocs/op
BenchmarkAll/compile__________:_map+baseline___________                     1000           2570874 ns/op            323584 B(heap)/op          500.0 domains     1456915 B/op       4443 allocs/op
BenchmarkAll/compile__________:_map+optimized__________                     1000            818229 ns/op            233472 B(heap)/op          500.0 domains     1399397 B/op       3414 allocs/op
BenchmarkAll/compile__________:_map+reverse+sort_______                     1000            778986 ns/op            143360 B(heap)/op          500.0 domains      925957 B/op       3824 allocs/op
PASS
ok      github.com/odinuge/cilium-fqdn-regex-benchmarking/benchmarks    63.334s
```

```
# For generate
$ go test -bench="BenchmarkAll/^generate" -test.cpu=1 -benchtime=100x ./benchmarks
goos: darwin
goarch: amd64
pkg: github.com/odinuge/cilium-fqdn-regex-benchmarking/benchmarks
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkAll/generate_________:_baseline_______________                     1000            209974 ns/op             57344 B(heap)/op          500.0 domains      159808 B/op       2204 allocs/op
BenchmarkAll/generate_________:_baseline-optimized_____                     1000            221192 ns/op             81920 B(heap)/op          500.0 domains      124504 B/op       1206 allocs/op
BenchmarkAll/generate_________:_reverse________________                     1000            231715 ns/op             81920 B(heap)/op          500.0 domains      143344 B/op       1705 allocs/op
BenchmarkAll/generate_________:_reverse+sort___________                     1000            305824 ns/op             81920 B(heap)/op          500.0 domains      143368 B/op       1706 allocs/op
BenchmarkAll/generate_________:_map+baseline___________                     1000            147165 ns/op             73728 B(heap)/op          500.0 domains      125864 B/op       1313 allocs/op
BenchmarkAll/generate_________:_map+optimized__________                     1000            123091 ns/op             98304 B(heap)/op          500.0 domains      110536 B/op        910 allocs/op
BenchmarkAll/generate_________:_map+reverse+sort_______                     1000            231630 ns/op             98304 B(heap)/op          500.0 domains      147760 B/op       1709 allocs/op
PASS
ok      github.com/odinuge/cilium-fqdn-regex-benchmarking/benchmarks    24.186s
```

### Matching performance for domains matching the policy

The most important aspect is that positive matches are fast and allocate as little mem as possible;


This is matching for matchNames;
```
go test -bench="BenchmarkAll/^match-name"   -test.cpu=1 -benchtime=1000x ./benchmarks
goos: darwin
goarch: amd64
pkg: github.com/odinuge/cilium-fqdn-regex-benchmarking/benchmarks
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkAll/match-name_______:_baseline_______________                     1000           4108928 ns/op               500.0 domains     6369144 B/op       1349 allocs/op
BenchmarkAll/match-name_______:_baseline-optimized_____                     1000            388407 ns/op               500.0 domains      812175 B/op        733 allocs/op
BenchmarkAll/match-name_______:_reverse________________                     1000            232877 ns/op               500.0 domains      677476 B/op        582 allocs/op
BenchmarkAll/match-name_______:_reverse+sort___________                     1000            118463 ns/op               500.0 domains      657772 B/op         73 allocs/op
BenchmarkAll/match-name_______:_map+baseline___________                     1000               568.7 ns/op             500.0 domains           0 B/op          0 allocs/op
BenchmarkAll/match-name_______:_map+optimized__________                     1000               619.1 ns/op             500.0 domains           0 B/op          0 allocs/op
BenchmarkAll/match-name_______:_map+reverse+sort_______                     1000               539.4 ns/op             500.0 domains           0 B/op          0 allocs/op
PASS
ok      github.com/odinuge/cilium-fqdn-regex-benchmarking/benchmarks    37.198s
```

And for matchPatterns;

```
go test -bench="BenchmarkAll/^match-pattern"   -test.cpu=1 -benchtime=1000x ./benchmarks
goos: darwin
goarch: amd64
pkg: github.com/odinuge/cilium-fqdn-regex-benchmarking/benchmarks
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkAll/match-pattern____:_baseline_______________                     1000           3425486 ns/op               500.0 domains     6369144 B/op       1349 allocs/op
BenchmarkAll/match-pattern____:_baseline-optimized_____                     1000            344051 ns/op               500.0 domains      812422 B/op        743 allocs/op
BenchmarkAll/match-pattern____:_reverse________________                     1000            291342 ns/op               500.0 domains      672062 B/op        526 allocs/op
BenchmarkAll/match-pattern____:_reverse+sort___________                     1000            181816 ns/op               500.0 domains      662955 B/op        231 allocs/op
BenchmarkAll/match-pattern____:_map+baseline___________                     1000            694274 ns/op               500.0 domains     1943984 B/op        750 allocs/op
BenchmarkAll/match-pattern____:_map+optimized__________                     1000            312033 ns/op               500.0 domains      681504 B/op        750 allocs/op
BenchmarkAll/match-pattern____:_map+reverse+sort_______                     1000            155229 ns/op               500.0 domains      662955 B/op        231 allocs/op
PASS
ok      github.com/odinuge/cilium-fqdn-regex-benchmarking/benchmarks    43.581s
```


### Matching performance for domains not matching the policy

Also, in case a port on an endpoint has multiple policies, we might have a domain that doesn't match the first policy, but number two. In that case, its good that performance of domains not matching are also fast.

When the domain has some similarities to one or more matchPatterns, but still doesn't match;

```
$ go test -bench="BenchmarkAll/^bad-match-pattern"   -test.cpu=1 -benchtime=1000x ./benchmarks
goos: darwin
goarch: amd64
pkg: github.com/odinuge/cilium-fqdn-regex-benchmarking/benchmarks
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkAll/bad-match-pattern:_baseline_______________                     1000           4520088 ns/op               500.0 domains     6369144 B/op       1349 allocs/op
BenchmarkAll/bad-match-pattern:_baseline-optimized_____                     1000            256548 ns/op               500.0 domains      811952 B/op        724 allocs/op
BenchmarkAll/bad-match-pattern:_reverse________________                     1000            245519 ns/op               500.0 domains      675283 B/op        559 allocs/op
BenchmarkAll/bad-match-pattern:_reverse+sort___________                     1000            153423 ns/op               500.0 domains      659843 B/op        136 allocs/op
BenchmarkAll/bad-match-pattern:_map+baseline___________                     1000            788536 ns/op               500.0 domains     1943984 B/op        750 allocs/op
BenchmarkAll/bad-match-pattern:_map+optimized__________                     1000            223319 ns/op               500.0 domains      681504 B/op        750 allocs/op
BenchmarkAll/bad-match-pattern:_map+reverse+sort_______                     1000            114129 ns/op               500.0 domains      658712 B/op        100 allocs/op
PASS
ok      github.com/odinuge/cilium-fqdn-regex-benchmarking/benchmarks    43.897s
```


Also, when the tld of the domain tested doesn't match any tld in any matchName or matchPattern;

```
go test -bench="BenchmarkAll/^bad-match-tld"   -test.cpu=1 -benchtime=1000x ./benchmarks
goos: darwin
goarch: amd64
pkg: github.com/odinuge/cilium-fqdn-regex-benchmarking/benchmarks
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkAll/bad-match-tld____:_baseline_______________                     1000           4602285 ns/op               500.0 domains     6369144 B/op       1349 allocs/op
BenchmarkAll/bad-match-tld____:_baseline-optimized_____                     1000            415886 ns/op               500.0 domains      812213 B/op        734 allocs/op
BenchmarkAll/bad-match-tld____:_reverse________________                     1000            137327 ns/op               500.0 domains      655872 B/op         14 allocs/op
BenchmarkAll/bad-match-tld____:_reverse+sort___________                     1000            134202 ns/op               500.0 domains      655872 B/op         14 allocs/op
BenchmarkAll/bad-match-tld____:_map+baseline___________                     1000            777455 ns/op               500.0 domains     1943984 B/op        750 allocs/op
BenchmarkAll/bad-match-tld____:_map+optimized__________                     1000            295559 ns/op               500.0 domains      681504 B/op        750 allocs/op
BenchmarkAll/bad-match-tld____:_map+reverse+sort_______                     1000             92962 ns/op               500.0 domains      655872 B/op         14 allocs/op
PASS
ok      github.com/odinuge/cilium-fqdn-regex-benchmarking/benchmarks    44.673s
```
