package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueue(t *testing.T) {
	assert := require.New(t)

	queue := NewQueue()
	queue.Push(1)
	queue.Push(2)
	queue.Push(3)
	assert.False(queue.IsEmpty())

	value, ok := queue.Front()
	assert.True(ok)
	assert.Equal(1, value)

	value, ok = queue.Back()
	assert.True(ok)
	assert.Equal(3, value)

	value, ok = queue.Pop()
	assert.True(ok)
	assert.Equal(1, value)
	assert.False(queue.IsEmpty())

	queue.Push(4)
	assert.False(queue.IsEmpty())
	value, ok = queue.Back()
	assert.True(ok)
	assert.Equal(4, value)

	wants := []int{2, 3, 4}
	i := 0
	queue.For(func(value interface{}) {
		assert.Equal(wants[i], value)
		i++
	})
}
