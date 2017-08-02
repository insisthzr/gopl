package utils

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type StackOp int

const (
	push StackOp = iota
	pop
)

const (
	stackOpStart = 0
	stackOpEnd   = 1
)

type StackOpRand struct {
	r *rand.Rand
}

func (r *StackOpRand) Rand() StackOp {
	len := stackOpEnd - stackOpStart + 1
	n := r.r.Intn(len)
	return StackOp(n - stackOpStart)
}

func NewStackOpRand(seed int64) *StackOpRand {
	return &StackOpRand{
		r: rand.New(rand.NewSource(seed)),
	}
}

func TestStack(t *testing.T) {
	const numRound = 1024 * 1024
	assert := require.New(t)
	now := time.Now().Unix()
	opRand := NewStackOpRand(now)
	nRand := rand.New(rand.NewSource(now))
	trustStack := []int{}
	stack := NewStack()

	for i := 0; i < numRound; i++ {
		op := opRand.Rand()
		switch op {
		case push:
			elem := nRand.Int()
			stack.Push(elem)
			trustStack = append(trustStack, elem)
		case pop:
			elem1, ok1 := stack.Pop()
			var elem2 interface{}
			ok2 := false
			length := len(trustStack)
			if length > 0 {
				ok2 = true
				elem2 = trustStack[length-1]
				trustStack = trustStack[:length-1]
			}
			assert.Equal(ok2, ok1)
			assert.Equal(elem2, elem1)
		default:
			panic("unreachable")
		}
	}
}
