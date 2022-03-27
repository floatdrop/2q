package twoqueue_test

import (
	"fmt"

	twoqueue "github.com/floatdrop/2q"
)

func ExampleTwoQueue() {
	cache := twoqueue.New[string, int](64, 128, 192)

	cache.Set("Hello", 5)

	if e := cache.Get("Hello"); e != nil {
		fmt.Println(*e)
		// Output: 5
	}
}
