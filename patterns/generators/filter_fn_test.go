package generators

import (
	"fmt"
	"testing"
)

func TestFilter(t *testing.T) {
	// Вместо примера использования
	t.Run("success", func(t *testing.T) {
		inStream := make(chan int, 1)
		go func() {
			defer close(inStream)

			for i := 0; i < 10; i++ {
				inStream <- i
			}
		}()

		doneCh := make(chan struct{})
		defer close(doneCh)

		filteredStream := Filter(doneCh, inStream, func(v int) bool {
			return v%2 == 0
		})

		for v := range filteredStream {
			fmt.Println(v)
		}
	})
}
