package stacks

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTreiberStack(t *testing.T) {
	t.Run("primitive", func(t *testing.T) {
		stack := NewTreiberStack[int]()
		assert.Equal(t, int64(0), stack.Size())

		stack.Push(1)
		stack.Push(2)
		stack.Push(3)
		assert.Equal(t, int64(3), stack.Size())

		value, ok := stack.Pop()
		assert.True(t, ok)
		assert.Equal(t, value, 3)

		_, _ = stack.Pop()
		_, _ = stack.Pop()
		assert.Equal(t, int64(0), stack.Size())

		_, ok = stack.Pop()
		assert.False(t, ok)
	})

	t.Run("user type", func(t *testing.T) {
		type uType struct {
			v int
		}

		stack := NewTreiberStack[uType]()
		assert.Equal(t, int64(0), stack.Size())

		stack.Push(uType{v: 1})
		stack.Push(uType{v: 2})
		stack.Push(uType{v: 3})
		assert.Equal(t, int64(3), stack.Size())

		value, ok := stack.Pop()
		assert.True(t, ok)
		assert.Equal(t, uType{v: 3}, value)

		_, _ = stack.Pop()
		_, _ = stack.Pop()
		assert.Equal(t, int64(0), stack.Size())

		_, ok = stack.Pop()
		assert.False(t, ok)
	})
}

func TestTreiberStack_Concurrent(t *testing.T) {
	t.Run("load x100", func(t *testing.T) {
		type uType struct {
			v int
		}

		for i := 0; i < 100; i++ {
			stack := NewTreiberStack[uType]()

			wg := sync.WaitGroup{}
			wg.Add(10)

			for i := 0; i < 10; i++ {
				go func() {
					defer wg.Done()

					for i := 0; i < 10000; i++ {
						stack.Push(uType{v: 1})
					}
					for i := 0; i < 2000; i++ {
						_, ok := stack.Pop()
						assert.True(t, ok)
					}
				}()
			}

			wg.Wait()

			expected := int64(10000*10 - 2000*10)
			assert.Equal(t, expected, stack.Size())

			for i := int64(0); i < expected; i++ {
				_, _ = stack.Pop()
			}

			assert.Equal(t, int64(0), stack.Size())
			_, ok := stack.Pop()
			assert.False(t, ok)
		}
	})
}
