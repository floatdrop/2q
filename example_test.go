package twoqueue_test

import (
	"fmt"

	twoqueue "github.com/floatdrop/2q"
)

func ExampleTwoQueue() {
	cache := twoqueue.New[string, int](256)

	cache.Set("Hello", 5)

	if e := cache.Get("Hello"); e != nil {
		fmt.Println(*e)
		// Output: 5
	}
}
