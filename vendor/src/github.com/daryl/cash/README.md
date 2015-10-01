# Cash

Cash is a super-simple, in-memory key/value cache for Go.

#### Install

    go get github.com/daryl/cash

#### Usage

```go
package main

import (
    "fmt"
    "github.com/daryl/cash"
    "time"
)

var c *cash.Cash

func main() {
    c = cash.New(cash.Conf{
        // Default expiration.
        10 * time.Minute,
        // Clean interval.
        30 * time.Minute,
    })

    // Cache forever (cash.Forever or -1).
    c.Set("foo", "bar", cash.Forever)

    // Default expiration (cash.Default or 0).
    c.Set("abcd", 1234, cash.Default)

    // Custom expiration time.
    c.Set("efgh", 5678, 3 * time.Minute)

    var foo string
    // Since you can store anything in the
    // cache, when you retrieve the value
    // you must use type assertion.
    if v, ok := c.Get("foo"); ok {
        foo = v.(string)
    }

    fmt.Println(foo) // bar
}
```

### Benchmarks

Some simple benchmarks. To test, just run: `go test -bench=.`:

```
BenchmarkGet	50000000	        51.3 ns/op
BenchmarkSet	10000000	       277 ns/op
BenchmarkHas	50000000	        44.0 ns/op
BenchmarkDel	50000000	        51.8 ns/op
BenchmarkClean	20000000	        93.2 ns/op
BenchmarkFlush	10000000	       152 ns/op
```

#### Documentation

For further documentation, check out [GoDoc](http://godoc.org/github.com/daryl/cash).
