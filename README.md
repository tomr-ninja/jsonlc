The more data you avoid to unmarshal, the faster it will be.

Doesn't make sense to use for small types like `int` or `bool`, but for big strings or structs it can be useful.

```
goos: linux
goarch: amd64
pkg: github.com/tomr-ninja/jsonlc
cpu: AMD Ryzen 5 5600H with Radeon Graphics         
BenchmarkUnmarshalJSON
BenchmarkUnmarshalJSON/string
BenchmarkUnmarshalJSON/string/short
BenchmarkUnmarshalJSON/string/short/standard
BenchmarkUnmarshalJSON/string/short/standard-12         	 1000000	      1222 ns/op	     248 B/op	       6 allocs/op
BenchmarkUnmarshalJSON/string/short/optimized
BenchmarkUnmarshalJSON/string/short/optimized-12        	 1239019	      1021 ns/op	     232 B/op	       5 allocs/op
BenchmarkUnmarshalJSON/string/long
BenchmarkUnmarshalJSON/string/long/standard
BenchmarkUnmarshalJSON/string/long/standard-12          	  134064	     10976 ns/op	    2288 B/op	       6 allocs/op
BenchmarkUnmarshalJSON/string/long/optimized
BenchmarkUnmarshalJSON/string/long/optimized-12         	  198180	      7162 ns/op	     232 B/op	       5 allocs/op
BenchmarkUnmarshalJSON/string/extremely_long
BenchmarkUnmarshalJSON/string/extremely_long/standard
BenchmarkUnmarshalJSON/string/extremely_long/standard-12         	    1092	   1150572 ns/op	  180464 B/op	       6 allocs/op
BenchmarkUnmarshalJSON/string/extremely_long/optimized
BenchmarkUnmarshalJSON/string/extremely_long/optimized-12        	    2257	    553487 ns/op	     232 B/op	       5 allocs/op
BenchmarkUnmarshalJSON/struct
BenchmarkUnmarshalJSON/struct/standard
BenchmarkUnmarshalJSON/struct/standard-12                        	  723272	      1699 ns/op	     312 B/op	       7 allocs/op
BenchmarkUnmarshalJSON/struct/optimized
BenchmarkUnmarshalJSON/struct/optimized-12                       	  887960	      1405 ns/op	     248 B/op	       6 allocs/op
```

### See also

- [easyjson](https://github.com/mailru/easyjson) - just a faster json marshaller/unmarshaller in general
(but requires code generation). Also it has `intern` and `nocopy` tags, which work similar, but for string only.
