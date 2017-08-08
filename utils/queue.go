package utils

type Queue struct {
	head *node
	tail *node
}
type node struct {
	value interface{}
	next  *node
}

func (q *Queue) IsEmpty() bool {
	return q.head == nil
}

func (q *Queue) Front() (interface{}, bool) {
	if q.IsEmpty() {
		return nil, false
	}
	return q.head.value, true
}

func (q *Queue) Back() (interface{}, bool) {
	if q.IsEmpty() {
		return nil, false
	}
	return q.tail.value, true
}

func (q *Queue) Pop() (interface{}, bool) {
	if q.IsEmpty() {
		return nil, false
	}
	out := q.head
	q.head = q.head.next
	if q.IsEmpty() {
		q.tail = nil
	}
	return out.value, true
}

func (q *Queue) Push(value interface{}) {
	if q.IsEmpty() {
		q.head = &node{value: value}
		q.tail = q.head
		return
	}
	newNode := &node{value: value}
	q.tail.next = newNode
	q.tail = newNode
}

func (q *Queue) For(fn func(interface{})) {
	for !q.IsEmpty() {
		value, _ := q.Pop()
		fn(value)
	}
}

func NewQueue() *Queue {
	return &Queue{}
}
