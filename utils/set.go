package utils

import (
	"bytes"
	"fmt"
)

type Set struct {
	m map[string]struct{}
}

func (s Set) Add(key string) {
	s.m[key] = struct{}{}
}

func (s Set) Exist(key string) bool {
	_, ok := s.m[key]
	return ok
}

func (s Set) Delete(key string) {
	delete(s.m, key)
}

func (s Set) String() string {
	if len(s.m) == 0 {
		return "[]"
	}
	buf := bytes.Buffer{}
	buf.WriteString("[")

	first := true
	for key := range s.m {
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

type SetConfig struct {
	Capacity int
}

var (
	DefaultSetConfig = &SetConfig{
		Capacity: 0,
	}
)

func NewSetWithConfig(config *SetConfig) *Set {
	if config.Capacity == 0 {
		config.Capacity = DefaultSetConfig.Capacity
	}
	return &Set{
		m: make(map[string]struct{}, config.Capacity),
	}
}

func NewSet() *Set {
	return NewSetWithConfig(DefaultSetConfig)
}
