package utils

type Stack struct {
	stack []interface{}
}

func (s *Stack) Len() int {
	return len(s.stack)
}

func (s *Stack) Push(v interface{}) {
	s.stack = append(s.stack, v)
}

func (s *Stack) Pop() (interface{}, bool) {
	len := len(s.stack)
	if len == 0 {
		return nil, false
	}
	v := s.stack[len-1]
	s.stack = s.stack[:len-1]
	return v, true
}

type StackConfig struct {
	Capacity int
}

var (
	DefaultStackConfig = &StackConfig{
		Capacity: 0,
	}
)

func NewStack() *Stack {
	return NewStackWithConfig(DefaultStackConfig)
}

func NewStackWithConfig(config *StackConfig) *Stack {
	if config == nil {
		return newStackWithConfig(DefaultStackConfig)
	}
	if config.Capacity == 0 {
		config.Capacity = DefaultStackConfig.Capacity
	}
	return newStackWithConfig(config)
}

func newStackWithConfig(config *StackConfig) *Stack {
	return &Stack{
		stack: make([]interface{}, 0, config.Capacity),
	}
}
