# 2q
[![Go Reference](https://pkg.go.dev/badge/github.com/floatdrop/2q.svg)](https://pkg.go.dev/github.com/floatdrop/2q)
[![CI](https://github.com/floatdrop/2q/actions/workflows/ci.yml/badge.svg)](https://github.com/floatdrop/2q/actions/workflows/ci.yml)
![Coverage](https://img.shields.io/badge/Coverage-70.6%25-yellow)
[![Go Report Card](https://goreportcard.com/badge/github.com/floatdrop/2q)](https://goreportcard.com/report/github.com/floatdrop/2q)

Thread safe GoLang [2Q](http://www.vldb.org/conf/1994/P439.PDF) cache.

## Example

```go
import (
	"fmt"

	twoqueue "github.com/floatdrop/2q"
)

func main() {
	cache := twoqueue.New[string, int](64, 128, 192)

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
	Benchmark2Q_Rand-8   	 4375202	       272.2 ns/op	      21 B/op	       1 allocs/op
	Benchmark2Q_Freq-8   	 4882906	       250.9 ns/op	      20 B/op	       1 allocs/op

hashicorp/golang-lru:
	Benchmark2Q_Rand-8    	 2847627	       411.9 ns/op	     135 B/op	       5 allocs/op
	Benchmark2Q_Freq-8    	 3323764	       354.2 ns/op	     122 B/op	       5 allocs/op
```
