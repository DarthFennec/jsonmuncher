Benchmarks
==========

These benchmarks compare a number of different Go JSON parsers. Four different
JSON files are parsed:

- `fixture_small.json`: 190 B
- `fixture_medium.json`: 2.3 KB
- `fixture_large.json`: 41 KB
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
BenchmarkJsonMuncherHuge-8                     	      10	 761287513 ns/op	    4336 B/op	       6 allocs/op
BenchmarkJsonParserHuge-8                      	       5	1135070002 ns/op	360846475 B/op	 1000012 allocs/op
BenchmarkEncodingJsonStructHuge-8              	       1	5435946960 ns/op	1442501104 B/op	10502039 allocs/op
BenchmarkEncodingJsonInterfaceHuge-8           	       1	5701351770 ns/op	1489402024 B/op	28034854 allocs/op
BenchmarkEncodingJsonStreamStructHuge-8        	       1	5735771777 ns/op	2167391392 B/op	10502057 allocs/op
BenchmarkEncodingJsonStreamInterfaceHuge-8     	       1	5966943523 ns/op	2214184568 B/op	28034497 allocs/op
BenchmarkJstreamHuge-8                         	       1	5866067676 ns/op	1129458008 B/op	28533146 allocs/op
BenchmarkGojayHuge-8                           	       3	2090216626 ns/op	1911409520 B/op	10002027 allocs/op
BenchmarkJsonIteratorHuge-8                    	       2	3225396334 ns/op	1839351112 B/op	17002141 allocs/op
BenchmarkGabsHuge-8                            	       1	5661634428 ns/op	1517427480 B/op	29534794 allocs/op
BenchmarkGoSimpleJsonHuge-8                    	       1	5876533891 ns/op	2563080408 B/op	28034667 allocs/op
BenchmarkFFJsonHuge-8                          	       3	2415598919 ns/op	1442499717 B/op	10502030 allocs/op
BenchmarkJasonHuge-8                           	       1	9854035703 ns/op	4191166648 B/op	49634480 allocs/op
BenchmarkUjsonHuge-8                           	       2	4026916938 ns/op	1593388936 B/op	34534906 allocs/op
BenchmarkDjsonHuge-8                           	       2	3048507867 ns/op	1489389136 B/op	28034807 allocs/op
BenchmarkUgorjiHuge-8                          	       1	7739377297 ns/op	2667890632 B/op	18503643 allocs/op
BenchmarkEasyJsonHuge-8                        	       3	2051049192 ns/op	1510499090 B/op	 9502023 allocs/op
BenchmarkJsonMuncherLarge-8                    	  100000	     94460 ns/op	    4336 B/op	       6 allocs/op
BenchmarkJsonParserLarge-8                     	  100000	    111023 ns/op	   49616 B/op	       7 allocs/op
BenchmarkEncodingJsonStructLarge-8             	   10000	    606856 ns/op	   56264 B/op	     248 allocs/op
BenchmarkEncodingJsonInterfaceLarge-8          	   10000	    842066 ns/op	  261799 B/op	    2881 allocs/op
BenchmarkEncodingJsonStreamStructLarge-8       	   10000	    655320 ns/op	  136168 B/op	     256 allocs/op
BenchmarkEncodingJsonStreamInterfaceLarge-8    	   10000	    931944 ns/op	  341692 B/op	    2889 allocs/op
BenchmarkJstreamLarge-8                        	    3000	   1969405 ns/op	  438465 B/op	    5484 allocs/op
BenchmarkGojayLarge-8                          	   50000	    153690 ns/op	  102668 B/op	     178 allocs/op
BenchmarkJsonIteratorLarge-8                   	   20000	    410814 ns/op	  118218 B/op	    1379 allocs/op
BenchmarkGabsLarge-8                           	   10000	    870205 ns/op	  265079 B/op	    3041 allocs/op
BenchmarkGoSimpleJsonLarge-8                   	   10000	   1000596 ns/op	  392635 B/op	    2845 allocs/op
BenchmarkFFJsonLarge-8                         	   30000	    248859 ns/op	   55977 B/op	     243 allocs/op
BenchmarkJasonLarge-8                          	   10000	   1044236 ns/op	  421071 B/op	    3284 allocs/op
BenchmarkUjsonLarge-8                          	   10000	    679211 ns/op	  288540 B/op	    4021 allocs/op
BenchmarkDjsonLarge-8                          	   10000	    531208 ns/op	  261144 B/op	    2746 allocs/op
BenchmarkUgorjiLarge-8                         	   10000	    798903 ns/op	   57458 B/op	     254 allocs/op
BenchmarkEasyJsonLarge-8                       	   50000	    175398 ns/op	   55096 B/op	     232 allocs/op
BenchmarkJsonMuncherMedium-8                   	  500000	     13783 ns/op	    1264 B/op	       6 allocs/op
BenchmarkJsonParserMedium-8                    	  500000	     16174 ns/op	    3536 B/op	       7 allocs/op
BenchmarkEncodingJsonStructMedium-8            	  200000	     42425 ns/op	    4626 B/op	      30 allocs/op
BenchmarkEncodingJsonInterfaceMedium-8         	  200000	     55449 ns/op	   13964 B/op	     213 allocs/op
BenchmarkEncodingJsonStreamStructMedium-8      	  200000	     44838 ns/op	    7692 B/op	      34 allocs/op
BenchmarkEncodingJsonStreamInterfaceMedium-8   	  200000	     55407 ns/op	   17036 B/op	     217 allocs/op
BenchmarkJstreamMedium-8                       	  100000	     84099 ns/op	   14713 B/op	     172 allocs/op
BenchmarkGojayMedium-8                         	  500000	     15310 ns/op	    6474 B/op	      20 allocs/op
BenchmarkJsonIteratorMedium-8                  	  200000	     30852 ns/op	    7615 B/op	     101 allocs/op
BenchmarkGabsMedium-8                          	  200000	     54700 ns/op	   14440 B/op	     232 allocs/op
BenchmarkGoSimpleJsonMedium-8                  	  100000	     59926 ns/op	   20603 B/op	     220 allocs/op
BenchmarkFFJsonMedium-8                        	  300000	     21394 ns/op	    4346 B/op	      25 allocs/op
BenchmarkJasonMedium-8                         	  100000	     63887 ns/op	   22443 B/op	     248 allocs/op
BenchmarkUjsonMedium-8                         	  200000	     41953 ns/op	   15203 B/op	     284 allocs/op
BenchmarkDjsonMedium-8                         	  200000	     33466 ns/op	   13659 B/op	     201 allocs/op
BenchmarkUgorjiMedium-8                        	  200000	     56543 ns/op	    5789 B/op	      36 allocs/op
BenchmarkEasyJsonMedium-8                      	  500000	     15691 ns/op	    3952 B/op	      19 allocs/op
BenchmarkJsonMuncherSmall-8                    	 1000000	      5937 ns/op	     496 B/op	       6 allocs/op
BenchmarkJsonParserSmall-8                     	 1000000	      7322 ns/op	    1168 B/op	       7 allocs/op
BenchmarkEncodingJsonStructSmall-8             	 1000000	     11285 ns/op	    1912 B/op	      23 allocs/op
BenchmarkEncodingJsonInterfaceSmall-8          	 1000000	     12018 ns/op	    2521 B/op	      39 allocs/op
BenchmarkEncodingJsonStreamStructSmall-8       	 1000000	      9257 ns/op	    1608 B/op	      22 allocs/op
BenchmarkEncodingJsonStreamInterfaceSmall-8    	 1000000	      9651 ns/op	    2217 B/op	      38 allocs/op
BenchmarkJstreamSmall-8                        	  200000	     40813 ns/op	   13289 B/op	      40 allocs/op
BenchmarkGojaySmall-8                          	 1000000	      7603 ns/op	    1520 B/op	      13 allocs/op
BenchmarkJsonIteratorSmall-8                   	 1000000	     10582 ns/op	    2001 B/op	      32 allocs/op
BenchmarkGabsSmall-8                           	  500000	     12494 ns/op	    2649 B/op	      47 allocs/op
BenchmarkGoSimplejsonSmall-8                   	  500000	     12923 ns/op	    3337 B/op	      39 allocs/op
BenchmarkFFJsonSmall-8                         	 1000000	      9163 ns/op	    1752 B/op	      21 allocs/op
BenchmarkJasonSmall-8                          	  300000	     24120 ns/op	    8333 B/op	     104 allocs/op
BenchmarkUjsonSmall-8                          	 1000000	     10737 ns/op	    2633 B/op	      46 allocs/op
BenchmarkDjsonSmall-8                          	 1000000	      9475 ns/op	    2345 B/op	      34 allocs/op
BenchmarkUgorjiSmall-8                         	 1000000	      9467 ns/op	    2304 B/op	      12 allocs/op
BenchmarkEasyJsonSmall-8                       	 1000000	      7948 ns/op	    1304 B/op	      15 allocs/op
```

Projects
--------

The benchmarks compare the following projects:

- [`github.com/darthfennec/jsonmuncher`](https://github.com/darthfennec/jsonmuncher)
- [`github.com/buger/jsonparser`](https://github.com/buger/jsonparser)
- [`encoding/json`](https://golang.org/pkg/encoding/json)
- [`github.com/bcicen/jstream`](https://github.com/bcicen/jstream)
- [`github.com/francoispqt/gojay`](https://github.com/francoispqt/gojay)
- [`github.com/json-iterator/go`](https://github.com/json-iterator/go)
- [`github.com/Jeffail/gabs`](https://github.com/Jeffail/gabs)
- [`github.com/bitly/go-simplejson`](https://github.com/bitly/go-simplejson)
- [`github.com/pquerna/ffjson`](https://github.com/pquerna/ffjson)
- [`github.com/antonholmquist/jason`](https://github.com/antonholmquist/jason)
- [`github.com/mreiferson/go-ujson`](https://github.com/mreiferson/go-ujson)
- [`github.com/a8m/djson`](https://github.com/a8m/djson)
- [`github.com/ugorji/go/codec`](https://github.com/ugorji/go/tree/master/codec)
- [`github.com/mailru/easyjson`](https://github.com/mailru/easyjson)
