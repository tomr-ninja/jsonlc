The more data you avoid to unmarshal, the faster it will be.

Doesn't make sense to use for small types like `int` or `bool`, but for big strings or structs it can be useful.

```
goos: linux
goarch: amd64
pkg: github.com/tomr-ninja/jsonlc
cpu: AMD Ryzen 5 5600H with Radeon Graphics         
BenchmarkUnmarshalJSON
BenchmarkUnmarshalJSON/short_string
BenchmarkUnmarshalJSON/short_string/standard
BenchmarkUnmarshalJSON/short_string/standard-12         	 1630040	       779.0 ns/op
BenchmarkUnmarshalJSON/short_string/optimized
BenchmarkUnmarshalJSON/short_string/optimized-12        	 1549267	       898.1 ns/op
BenchmarkUnmarshalJSON/long_string
BenchmarkUnmarshalJSON/long_string/standard
BenchmarkUnmarshalJSON/long_string/standard-12          	  136579	     10699 ns/op
BenchmarkUnmarshalJSON/long_string/optimized
BenchmarkUnmarshalJSON/long_string/optimized-12         	  196988	      7484 ns/op
BenchmarkUnmarshalJSON/extremely_long_string
BenchmarkUnmarshalJSON/extremely_long_string/standard
BenchmarkUnmarshalJSON/extremely_long_string/standard-12         	    1479	    950571 ns/op
BenchmarkUnmarshalJSON/extremely_long_string/optimized
BenchmarkUnmarshalJSON/extremely_long_string/optimized-12        	    2275	    536340 ns/op
```
