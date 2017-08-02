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
	Length   int
	Capacity int
}

var (
	defaultStackConfig = &StackConfig{
		Length:   0,
		Capacity: 0,
	}
)

func NewStack() *Stack {
	return NewStackWithConfig(defaultStackConfig)
}

func NewStackWithConfig(config *StackConfig) *Stack {
	if config == nil {
		return newStackWithConfig(defaultStackConfig)
	}
	if config.Length == 0 {
		config.Length = defaultStackConfig.Length
	}
	if config.Capacity == 0 {
		config.Capacity = defaultStackConfig.Capacity
	}
	return newStackWithConfig(config)
}

func newStackWithConfig(config *StackConfig) *Stack {
	return &Stack{
		stack: make([]interface{}, config.Length, config.Capacity),
	}
}
