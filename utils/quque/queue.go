package queue

import "fmt"

var (
	QueueEmpty = fmt.Errorf("queue is empty")
)

type Queue struct {
	head *node
	tail *node
}

func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) Empty() bool {
	return q.head == nil
}

func (q *Queue) Pop() (interface{}, error) {
	if q.head == nil {
		return nil, QueueEmpty
	}
	out := q.head
	q.head = q.head.Next
	if q.head == nil {
		q.tail = nil
	}
	return out.Value, nil
}

func (q *Queue) Push(value interface{}) {
	if q.head == nil {
		q.head = &node{Value: value}
		q.tail = q.head
		return
	}
	newNode := &node{Value: value}
	q.tail.Next = newNode
	q.tail = newNode
}

type node struct {
	Value interface{}
	Next  *node
}
