package maps

import (
	"fmt"
	"sync"
	"testing"
)

func originalMapConcurrent() {
	m := make(map[string]int)
	mu := sync.RWMutex{}

	wg := sync.WaitGroup{}
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func(ind int) {
			defer wg.Done()

			mu.Lock()
			for j := 0; j < 1000; j++ {
				k := fmt.Sprintf("key%d_%d", ind, j)

				m[k] = j
			}
			mu.Unlock()

			mu.RLock()
			for j := 0; j < 1000; j++ {
				k := fmt.Sprintf("key%d_%d", ind, j)

				if _, ok := m[k]; !ok {
					panic(k)
				}
			}
			mu.RUnlock()

			mu.Lock()
			for j := 0; j < 1000; j++ {
				k := fmt.Sprintf("key%d_%d", ind, j)

				delete(m, k)
			}
			mu.Unlock()
		}(i)
	}

	wg.Wait()
}

func shardedMapConcurrent() {
	m := NewShardedMap[string, int](256)

	wg := sync.WaitGroup{}
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func(ind int) {
			defer wg.Done()

			for j := 0; j < 1000; j++ {
				k := fmt.Sprintf("key%d_%d", ind, j)

				m.Set(k, j)
			}

			for j := 0; j < 1000; j++ {
				k := fmt.Sprintf("key%d_%d", ind, j)

				if _, ok := m.Get(k); !ok {
					panic(k)
				}
			}

			for j := 0; j < 1000; j++ {
				k := fmt.Sprintf("key%d_%d", ind, j)

				m.Delete(k)
			}
		}(i)
	}

	wg.Wait()
}

func BenchmarkOriginalMap(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		originalMapConcurrent()
	}
}

func BenchmarkShardedMap(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		shardedMapConcurrent()
	}
}
