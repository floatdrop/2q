package twoqueue

import (
	"math/rand"
	"testing"
)

func Benchmark2Q_Rand(b *testing.B) {
	l := New[int64, int64](2048, 4096, 6144)

	trace := make([]int64, b.N*2)
	for i := 0; i < b.N*2; i++ {
		trace[i] = rand.Int63() % 32768
	}

	b.ResetTimer()

	var hit, miss int
	for i := 0; i < 2*b.N; i++ {
		if i%2 == 0 {
			l.Set(trace[i], trace[i])
		} else {
			if l.Get(trace[i]) == nil {
				miss++
			} else {
				hit++
			}
		}
	}
	b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))
}

func Benchmark2Q_Freq(b *testing.B) {
	l := New[int64, int64](2048, 4096, 6144)

	trace := make([]int64, b.N*2)
	for i := 0; i < b.N*2; i++ {
		if i%2 == 0 {
			trace[i] = rand.Int63() % 16384
		} else {
			trace[i] = rand.Int63() % 32768
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Set(trace[i], trace[i])
	}
	var hit, miss int
	for i := 0; i < b.N; i++ {
		if l.Get(trace[i]) == nil {
			miss++
		} else {
			hit++
		}
	}
	b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))
}

func Test2Q(t *testing.T) {
	l := New[int, int](32, 64, 96)

	for i := 0; i < 256; i++ {
		l.Set(i, i)
	}

	for i := 0; i < 128; i++ {
		if e := l.Get(i); e != nil {
			t.Fatalf("should be evicted")
		}
	}
	// for i := 128; i < 256; i++ {
	// 	if e := l.Get(i); e == nil {
	// 		t.Fatalf("should not be evicted")
	// 	}
	// }
}
