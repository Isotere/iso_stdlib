package generators

import (
	"fmt"
	"math/rand/v2"
	"testing"
)

func TestTake(t *testing.T) {
	// Вместо примера использования
	t.Run("success", func(t *testing.T) {
		fn := func() int {
			return rand.N(5000)
		}

		done := make(chan struct{})

		for v := range Take(done, RepeatFn(done, fn), 10) {
			fmt.Println(v)
		}
	})
}
