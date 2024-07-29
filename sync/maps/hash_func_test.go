package maps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashFunc(t *testing.T) {
	t.Run("string hash", func(t *testing.T) {
		v := "some string"
		h := hash(v)

		for i := 0; i < 10; i++ {
			assert.Equal(t, h, hash(v))
		}
	})
}
