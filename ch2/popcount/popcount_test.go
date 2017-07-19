package popcount

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPopcount(t *testing.T) {
	assert := require.New(t)
	var n uint64 = 12312312445123
	c1 := PopCount(n)
	c2 := PopCountSlow(n)
	c3 := PopCountSlow1(n)
	c4 := PopCountSlow2(n)
	assert.Equal(c1, c2)
	assert.Equal(c1, c3)
	assert.Equal(c1, c4)
}

func BenchmarkPopcount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCount(123123123132)
	}
}

func BenchmarkPopcountSlow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCountSlow(123123123132)
	}
}

func BenchmarkPopcountSlow1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCountSlow1(123123123132)
	}
}

func BenchmarkPopcountSlow2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCountSlow2(123123123132)
	}
}
