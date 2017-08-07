package intset

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLenZeroInitially(t *testing.T) {
	assert := require.New(t)

	s := &IntSet{}
	assert.Equal(0, s.Len())
}

func TestLenAfterAddingElements(t *testing.T) {
	assert := require.New(t)

	s := &IntSet{}
	s.Add(0)
	s.Add(2000)
	assert.Equal(2, s.Len())
}

func TestRemove(t *testing.T) {
	assert := require.New(t)

	s := &IntSet{}
	s.Add(0)
	s.Remove(0)
	assert.Equal(0, s.Len())
}

func TestClear(t *testing.T) {
	assert := require.New(t)

	s := &IntSet{}
	s.Add(0)
	s.Add(1000)
	s.Clear()

	assert.False(s.Has(0))
	assert.False(s.Has(1000))
}

func TestCopy(t *testing.T) {
	assert := require.New(t)

	orig := &IntSet{}
	orig.Add(1)
	copy := orig.Copy()
	copy.Add(2)

	assert.True(copy.Has(1))
	assert.False(orig.Has(2))
}

func TestAddAll(t *testing.T) {
	assert := require.New(t)

	s := &IntSet{}
	s.AddAll(0, 2, 4)
	assert.True(s.Has(0))
	assert.True(s.Has(2))
	assert.True(s.Has(4))
}

func TestIntersectWith(t *testing.T) {
	assert := require.New(t)

	s := &IntSet{}
	s.AddAll(0, 2, 4)
	u := &IntSet{}
	u.AddAll(1, 2, 3)
	s.IntersectWith(u)
	assert.True(s.Has(2))
	assert.Equal(1, s.Len())
}

func TestDifferenceWith(t *testing.T) {
	assert := require.New(t)

	s := &IntSet{}
	s.AddAll(0, 2, 4)
	u := &IntSet{}
	u.AddAll(1, 2, 3)
	s.DifferenceWith(u)
	expected := &IntSet{}
	expected.AddAll(0, 4)
	assert.Equal(s.String(), expected.String())
}

func TestSymmetricDifference(t *testing.T) {
	assert := require.New(t)

	s := &IntSet{}
	s.AddAll(0, 2, 4)
	u := &IntSet{}
	u.AddAll(1, 2, 3)
	s.SymmetricDifference(u)
	expected := &IntSet{}
	expected.AddAll(0, 1, 3, 4)
	assert.Equal(s.String(), expected.String())
}
