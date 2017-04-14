package utils

import (
	"bytes"
	"fmt"
)

type Set struct {
	m map[string]struct{}
}

func NewSet() *Set {
	return &Set{m: map[string]struct{}{}}
}

func (s *Set) Add(key string) {
	s.m[key] = struct{}{}
}

func (s *Set) Exist(key string) bool {
	_, ok := s.m[key]
	return ok
}

func (s *Set) Delete(key string) {
	delete(s.m, key)
}

func (s *Set) String() string {
	if len(s.m) == 0 {
		return ""
	}
	buf := bytes.Buffer{}
	buf.WriteString("[")

	first := true
	for key, _ := range s.m {
		if first {
			first = false
		} else {
			buf.WriteString(",")
		}
		buf.WriteString(fmt.Sprintf("%q", key))
	}

	buf.WriteString("]")
	return buf.String()
}
