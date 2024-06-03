package generators

import (
	"fmt"
	"testing"
	"time"
)

func TestRepeat(t *testing.T) {
	// Вместо примера использования
	t.Run("success", func(t *testing.T) {
		inValues := []string{"one", "two", "three"}

		done := make(chan struct{})

		genCh := Repeat(done, inValues...)

		time.AfterFunc(time.Microsecond*100, func() { close(done) })

		for v := range genCh {
			fmt.Println(v)
		}
	})
}
