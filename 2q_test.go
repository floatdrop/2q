package twoqueue

import (
	"math/rand"
	"testing"
)

func Benchmark2Q_Rand(b *testing.B) {
	l := New[int64, int64](8196)

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
	l := New[int64, int64](8196)

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
	l := New[int, int](128)

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

func Test2Q_Add_RecentEvict(t *testing.T) {
	l := NewParams[int, int](4, 4, 4)

	// Set 1,2,3,4,5 -> Evict 1
	l.Set(1, 1)
	l.Set(2, 2)
	l.Set(3, 3)
	l.Set(4, 4)

	if e := l.Peek(2); e == nil || *e != 2 {
		t.Fatalf("peek should return 2 from recent: %+v", e)
	}

	if e := l.Set(5, 5); e == nil || e.Key != 1 {
		t.Fatalf("1 should be evicted to ghosted elements: %+v", e)
	}

	if n := l.recent.Len(); n != 4 {
		t.Fatalf("recent should have size of 4: %d", n)
	}

	if n := l.recentEvict.Len(); n != 1 {
		t.Fatalf("one element should be in recentEvict: %d", n)
	}

	if n := l.frequent.Len(); n != 0 {
		t.Fatalf("no elemenets should be in frequent: %d", n)
	}

	// Pull in the recently evicted
	if e := l.Set(1, 1); e != nil {
		t.Fatalf("frequent item should go straight to frequent: %+v", e)
	}

	if l.recent.Len() != 4 || l.frequent.Len() != 1 || l.recentEvict.Len() != 0 {
		t.Fatalf("recently evicted should go to frequent")
	}

	// Local brust to recent should not change frequent
	l.Set(3, 3)
	l.Set(3, 3)
	l.Set(3, 3)
	l.Set(3, 3)

	if l.recent.Len() != 4 || l.frequent.Len() != 1 || l.recentEvict.Len() != 0 {
		t.Fatalf("recently evicted should go to frequent")
	}

	if l.Len() != 5 {
		t.Fatalf("bad length: %v != 5", l.Len())
	}

	if e := l.Get(1); e == nil || *e != 1 {
		t.Fatalf("get should return 1 from frequent: %+v", e)
	}

	if e := l.Peek(1); e == nil || *e != 1 {
		t.Fatalf("peek should return 1 from frequent: %+v", e)
	}

	if e := l.Remove(3); e == nil {
		t.Fatalf("recent item should be removed: %+v", e)
	}

	if e := l.Remove(1); e == nil {
		t.Fatalf("frequent item should be removed: %+v", e)
	}
}
