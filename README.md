# 2q
[![Go Reference](https://pkg.go.dev/badge/github.com/floatdrop/2q.svg)](https://pkg.go.dev/github.com/floatdrop/2q)
[![CI](https://github.com/floatdrop/2q/actions/workflows/ci.yml/badge.svg)](https://github.com/floatdrop/2q/actions/workflows/ci.yml)
![Coverage](https://img.shields.io/badge/Coverage-88.9%25-brightgreen)
[![Go Report Card](https://goreportcard.com/badge/github.com/floatdrop/2q)](https://goreportcard.com/report/github.com/floatdrop/2q)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)

Thread safe GoLang [2Q](http://www.vldb.org/conf/1994/P439.PDF) cache.

## Example

```go
import (
	"fmt"

	twoqueue "github.com/floatdrop/2q"
)

func main() {
	cache := twoqueue.New[string, int](256)

	cache.Set("Hello", 5)

	if e := cache.Get("Hello"); e != nil {
		fmt.Println(*e)
		// Output: 5
	}
}
```

## TTL

See [LRU TTL example](https://github.com/floatdrop/lru#ttl).

## Benchmarks

```
floatdrop/twoqueue:
	Benchmark2Q_Rand-8   	 4384994	       264.5 ns/op	      46 B/op	       3 allocs/op
	Benchmark2Q_Freq-8   	 4862632	       243.9 ns/op	      44 B/op	       3 allocs/op

hashicorp/golang-lru:
	Benchmark2Q_Rand-8    	 2847627	       411.9 ns/op	     135 B/op	       5 allocs/op
	Benchmark2Q_Freq-8    	 3323764	       354.2 ns/op	     122 B/op	       5 allocs/op
```
