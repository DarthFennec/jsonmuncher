Benchmarks
==========

These benchmarks compare a number of different Go JSON parsers. Four different
JSON files are parsed:

- `fixture_small.json`: 190 B
- `fixture_large.json`: 2.3 KB
- `fixture_medium.json`: 41 KB
- `fixture_huge.json`: 333 MB

To run the benchmarks yourself, use `make build` to build the Docker image, then
`make bench` to run it. Alternatively, use `go get` to install all the parsers
being benchmarked, then run:

``` shell
go test -bench=. -benchtime=5s -benchmem ./benchmark/
```

These benchmark tests are originally from
[`buger/jsonparser`](https://github.com/buger/jsonparser).

Results
-------

These results were taken on a MacBook Pro 11,2. Your mileage may vary slightly.

```
BenchmarkJsonMuncherHuge-8               	      10	 859714936 ns/op	    4342 B/op	       6 allocs/op
BenchmarkJsonParserHuge-8                	       5	1147621332 ns/op	360846811 B/op	 1000012 allocs/op
BenchmarkEncodingJsonStructHuge-8        	       1	5452741452 ns/op	1442500544 B/op	10502039 allocs/op
BenchmarkEncodingJsonInterfaceHuge-8     	       1	5700589742 ns/op	1489375432 B/op	28034766 allocs/op
BenchmarkJsonIteratorHuge-8              	       2	3322275736 ns/op	1839352848 B/op	17002144 allocs/op
BenchmarkGabsHuge-8                      	       1	5822377683 ns/op	1517439336 B/op	29534833 allocs/op
BenchmarkGoSimpleJsonHuge-8              	       1	5972472324 ns/op	2563132440 B/op	28034845 allocs/op
BenchmarkFFJsonHuge-8                    	       3	2475434876 ns/op	1442499509 B/op	10502029 allocs/op
BenchmarkJasonHuge-8                     	       1	9590892930 ns/op	4191279992 B/op	49634929 allocs/op
BenchmarkUjsonHuge-8                     	       2	3945387376 ns/op	1593325728 B/op	34534686 allocs/op
BenchmarkDjsonHuge-8                     	       2	2857508424 ns/op	1489454784 B/op	28035037 allocs/op
BenchmarkUgorjiHuge-8                    	       1	7792939090 ns/op	2667890904 B/op	18503643 allocs/op
BenchmarkEasyJsonHuge-8                  	       3	2041085219 ns/op	1510498482 B/op	 9502021 allocs/op
BenchmarkJsonMuncherLarge-8              	   50000	    108270 ns/op	    4336 B/op	       6 allocs/op
BenchmarkJsonParserLarge-8               	  100000	    111169 ns/op	   49616 B/op	       7 allocs/op
BenchmarkEncodingJsonStructLarge-8       	   10000	    612850 ns/op	   56265 B/op	     248 allocs/op
BenchmarkEncodingJsonInterfaceLarge-8    	   10000	    841988 ns/op	  261826 B/op	    2881 allocs/op
BenchmarkJsonIteratorLarge-8             	   20000	    410673 ns/op	  118215 B/op	    1379 allocs/op
BenchmarkGabsLarge-8                     	   10000	    855958 ns/op	  265119 B/op	    3041 allocs/op
BenchmarkGoSimpleJsonLarge-8             	   10000	    991333 ns/op	  392650 B/op	    2845 allocs/op
BenchmarkFFJsonLarge-8                   	   30000	    250692 ns/op	   55976 B/op	     243 allocs/op
BenchmarkJasonLarge-8                    	   10000	   1022959 ns/op	  421090 B/op	    3284 allocs/op
BenchmarkUjsonLarge-8                    	   10000	    666410 ns/op	  288516 B/op	    4021 allocs/op
BenchmarkDjsonLarge-8                    	   10000	    523898 ns/op	  261149 B/op	    2746 allocs/op
BenchmarkUgorjiLarge-8                   	   10000	    798633 ns/op	   57458 B/op	     254 allocs/op
BenchmarkEasyJsonLarge-8                 	   50000	    174619 ns/op	   55096 B/op	     232 allocs/op
BenchmarkJsonMuncherMedium-8             	  500000	     14731 ns/op	    1264 B/op	       6 allocs/op
BenchmarkJsonParserMedium-8              	  500000	     16866 ns/op	    3536 B/op	       7 allocs/op
BenchmarkEncodingJsonStructMedium-8      	  200000	     43188 ns/op	    4626 B/op	      30 allocs/op
BenchmarkEncodingJsonInterfaceMedium-8   	  200000	     52360 ns/op	   13964 B/op	     213 allocs/op
BenchmarkJsonIteratorMedium-8            	  200000	     31346 ns/op	    7614 B/op	     101 allocs/op
BenchmarkGabsMedium-8                    	  200000	     54059 ns/op	   14440 B/op	     232 allocs/op
BenchmarkGoSimpleJsonMedium-8            	  100000	     59932 ns/op	   20603 B/op	     220 allocs/op
BenchmarkFFJsonMedium-8                  	  300000	     22878 ns/op	    4346 B/op	      25 allocs/op
BenchmarkJasonMedium-8                   	  100000	     63614 ns/op	   22444 B/op	     248 allocs/op
BenchmarkUjsonMedium-8                   	  200000	     42901 ns/op	   15203 B/op	     284 allocs/op
BenchmarkDjsonMedium-8                   	  200000	     33364 ns/op	   13659 B/op	     201 allocs/op
BenchmarkUgorjiMedium-8                  	  200000	     57036 ns/op	    5789 B/op	      36 allocs/op
BenchmarkEasyJsonMedium-8                	  500000	     16699 ns/op	    3952 B/op	      19 allocs/op
BenchmarkJsonMuncherSmall-8              	 1000000	      6339 ns/op	     496 B/op	       6 allocs/op
BenchmarkJsonParserSmall-8               	 1000000	      8126 ns/op	    1168 B/op	       7 allocs/op
BenchmarkEncodingJsonStructSmall-8       	  500000	     12587 ns/op	    1912 B/op	      23 allocs/op
BenchmarkEncodingJsonInterfaceSmall-8    	  500000	     12756 ns/op	    2521 B/op	      39 allocs/op
BenchmarkJsonIteratorSmall-8             	 1000000	     11449 ns/op	    2001 B/op	      32 allocs/op
BenchmarkGabsSmall-8                     	  500000	     13063 ns/op	    2649 B/op	      47 allocs/op
BenchmarkGoSimplejsonSmall-8             	  500000	     13217 ns/op	    3337 B/op	      39 allocs/op
BenchmarkFFJsonSmall-8                   	 1000000	     10190 ns/op	    1752 B/op	      21 allocs/op
BenchmarkJasonSmall-8                    	  300000	     25239 ns/op	    8333 B/op	     104 allocs/op
BenchmarkUjsonSmall-8                    	 1000000	     11632 ns/op	    2633 B/op	      46 allocs/op
BenchmarkDjsonSmall-8                    	 1000000	     10051 ns/op	    2345 B/op	      34 allocs/op
BenchmarkUgorjiSmall-8                   	 1000000	      9877 ns/op	    2304 B/op	      12 allocs/op
BenchmarkEasyJsonSmall-8                 	 1000000	      8815 ns/op	    1304 B/op	      15 allocs/op
```

Projects
--------

The benchmarks compare the following projects:

- [`github.com/darthfennec/jsonmuncher`](https://github.com/darthfennec/jsonmuncher)
- [`github.com/buger/jsonparser`](https://github.com/buger/jsonparser)
- [`encoding/json`](https://golang.org/pkg/encoding/json)
- [`github.com/json-iterator/go`](https://github.com/json-iterator/go)
- [`github.com/Jeffail/gabs`](https://github.com/Jeffail/gabs)
- [`github.com/bitly/go-simplejson`](https://github.com/bitly/go-simplejson)
- [`github.com/pquerna/ffjson`](https://github.com/pquerna/ffjson)
- [`github.com/antonholmquist/jason`](https://github.com/antonholmquist/jason)
- [`github.com/mreiferson/go-ujson`](https://github.com/mreiferson/go-ujson)
- [`github.com/a8m/djson`](https://github.com/a8m/djson)
- [`github.com/ugorji/go/codec`](https://github.com/ugorji/go/tree/master/codec)
- [`github.com/mailru/easyjson`](https://github.com/mailru/easyjson)
