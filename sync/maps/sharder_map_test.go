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

	m.Set("value2", 200)
	m.Set("value3", 300)
	m.Set("value4", 400)
	m.Set("value5", 500)

	keys := m.Keys()
	assert.Equal(t, 4, len(keys))
}
