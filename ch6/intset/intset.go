package intset

import (
	"bytes"
	"fmt"

	"github.com/insisthzr/gopl/ch2/popcount"
)

type IntSet struct {
	words []uint64
}

func (s *IntSet) Has(x int) bool {
	word := x / 64
	bit := uint64(x % 64)
	if word < len(s.words) {
		if s.words[word]&(1<<bit) != 0 {
			return true
		}
	}
	return false
}

func (s *IntSet) Add(x int) {
	word := x / 64
	bit := uint64(x % 64)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= (1 << bit)
}

func (s *IntSet) AddAll(xs ...int) {
	for _, x := range xs {
		s.Add(x)
	}
}

func (s *IntSet) Remove(x int) {
	word := x / 64
	bit := uint64(x % 64)
	if word < len(s.words) {
		y := uint64(1 << bit)
		z := ^y
		s.words[word] &= z
	}
}

func (s *IntSet) Clear() {
	s.words = nil
}

func (s *IntSet) Copy() *IntSet {
	words := make([]uint64, len(s.words))
	copy(words, s.words) // depend on slice's length
	return &IntSet{words: words}
}

func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

func (s *IntSet) IntersectWith(t *IntSet) {
	for i, tword := range t.words {
		if i >= len(s.words) {
			break
		}
		s.words[i] &= tword
	}
}

func (s *IntSet) DifferenceWith(t *IntSet) {
	for i, tword := range t.words {
		if i >= len(s.words) {
			break
		}
		s.words[i] &^= tword
	}
}

func (s *IntSet) SymmetricDifference(t *IntSet) {
	sCopy := s.Copy()
	tCopy := t.Copy()
	s.DifferenceWith(t)
	tCopy.DifferenceWith(sCopy)
	s.UnionWith(tCopy)
}

func (s *IntSet) Len() int {
	sum := 0
	for _, word := range s.words {
		sum += popcount.PopCount(word)
	}
	return sum
}

func (s *IntSet) Elems() []int {
	elems := []int{}
	for i, word := range s.words {
		for j := 0; j < 64; j++ {
			ju := uint64(j)
			if word&uint64(1<<ju) != 0 {
				elems = append(elems, 64*i+j)
			}
		}
	}
	return elems
}

func (s *IntSet) String() string {
	buf := &bytes.Buffer{}
	buf.WriteByte('{')
	first := true
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if word&(1<<uint(j)) != 0 {
				if first {
					first = false
				} else {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(buf, "%d", 64*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}
