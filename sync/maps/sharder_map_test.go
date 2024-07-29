package maps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShardedMap(t *testing.T) {
	m := NewShardedMap[string, int](256)

	m.Set("value1", 12)
	assert.Equal(t, int64(1), m.Len())

	m.Set("value1", 100)
	assert.Equal(t, int64(1), m.Len())

	val, ok := m.Get("value1")
	assert.True(t, ok)
	assert.Equal(t, 100, val)

	m.Delete("value1")
	assert.Equal(t, int64(0), m.Len())
}
