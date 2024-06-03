package generators

import (
	"fmt"
	"testing"
)

func TestMap(t *testing.T) {
	// Вместо примера использования
	t.Run("success", func(t *testing.T) {
		inStream := make(chan int, 1)
		go func() {
			defer close(inStream)

			for i := 10; i < 100; i += 10 {
				inStream <- i
			}
		}()

		doneCh := make(chan struct{})
		defer close(doneCh)

		mappedStream := Map(doneCh, inStream, func(v int) string {
			return fmt.Sprintf("Int as string: %d", v)
		})

		for v := range mappedStream {
			fmt.Println(v)
		}
	})
}
