package generators

import (
	"fmt"
	"math/rand/v2"
	"testing"
	"time"
)

func TestRepeatFn(t *testing.T) {
	// Вместо примера использования
	t.Run("success", func(t *testing.T) {
		fn := func() int {
			return rand.N(5000)
		}

		done := make(chan struct{})

		genCh := RepeatFn(done, fn)

		time.AfterFunc(time.Microsecond*100, func() { close(done) })

		for v := range genCh {
			fmt.Println(v)
		}
	})
}
