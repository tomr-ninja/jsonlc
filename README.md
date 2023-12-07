**WORK IN PROGRESS**

Avoid unmarshalling parts of your JSON messages that don't change.

The more data you avoid to unmarshal, the faster it will be: long strings, big structs, etc. Doesn't make sense to use
for small types like `int` or `bool`.

Move rarely changing fields to a separate struct to maximize the effect, e.g. server name/host, app version,
environment variables, etc. - everything that doesn't produce a lot of entropy when logged together.

This technique is similar to [string interning](https://en.wikipedia.org/wiki/String_interning), but for any type.
Mostly inspired by [LowCardinality](https://clickhouse.com/docs/en/sql-reference/data-types/lowcardinality) column
type in ClickHouse.

## Usage

```go
type (
    MyStructOptimized struct {
        FieldOne int                           `json:"field_one"`
        // just use jsonlc.LowCardinality[T] instead of T (almost a drop-in replacement)
        FieldTwo jsonlc.LowCardinality[string] `json:"field_two"`
    }
)

func myFunc() error {
    var s MyStructOptimized
    if err := json.Unmarshal(data, &s); err != nil { // unmarshal as usual
        return err
    }
    if s.FieldTwo.Value() == "hello" { // immutable, readonly. Use Pointer() for no-copy access.
        // ...	
    }
    // ...
}
```

Some benchmarks (edited for readability):

```    
BenchmarkUnmarshalJSON

string
string/short
string/short/standard             1000000       1222 ns/op       248 B/op    6 allocs/op
string/short/optimized            1239019       1021 ns/op       232 B/op    5 allocs/op
string/long/standard               134064      10976 ns/op      2288 B/op    6 allocs/op
string/long/optimized              198180       7162 ns/op       232 B/op    5 allocs/op
string/extremely_long/standard       1092    1150572 ns/op    180464 B/op    6 allocs/op
string/extremely_long/optimized      2257     553487 ns/op       232 B/op    5 allocs/op
struct/standard                    723272       1699 ns/op       312 B/op    7 allocs/op
struct/optimized                   887960       1405 ns/op       248 B/op    6 allocs/op
```

You can see that there is always 1 less allocation. "struct" bench uses relatively small struct with 4 integers,
and it still shows ~18% improvement.

Besides that, keep in mind that by allocating less memory you also reduce GC pressure, which is always good.

## Why only JSON?

- You might want the same functionality for other formats, but JSON is the most popular one, so it's a good start.
- At the moment it's ~70 LOC, so if you understand how it works, you can easily implement it for other formats.
- I've also thought about supporting `msgpack`, but it's not present in the `encoding` package in the standard library. Sorry. 

## See also

- [easyjson](https://github.com/mailru/easyjson) - just a *way* faster json marshaller/unmarshaller in general
(but requires code generation). Also, it has `intern` and `nocopy` tags, which work similar, but for strings only.
- [msgp](https://github.com/tinylib/msgp) - similar to `easyjson`, but for msgpack format.
