package main

import (
	"bytes"
	"fmt"
	"github.com/insisthzr/gopl/utils/queue"
	"strings"
)

var (
	Size = 1
	space = strings.Repeat(" ", Size)
)

type tree struct {
	value int
	left  *tree
	right *tree
}

func (t *tree) display() {
	if t == nil {
		return
	}
	fmt.Println(t.value)
	t.left.display()
	t.right.display()
}

func (t *tree) Height() int {
	if t == nil {
		return 0
	}
	lheight := t.left.Height()
	rheight := t.right.Height()
	if lheight > rheight {
		return lheight + 1
	}
	return rheight + 1
}

func (t *tree) Stringy() string {

	height := t.Height()
	//width, _ := twoPowN(height)

	m := genHalfMatrix(height)

	t.bfs(func(node *withPosition) {
		m[node.Level][node.Index] = node
	})

	return m.format()
}

type matrix [][]*withPosition

func genHalfMatrix(height int) matrix {
	m := make([][]*withPosition, height)
	width := 1
	for i, _ := range m {
		m[i] = make([]*withPosition, width)
		width *= 2
	}
	return m
}

func (m matrix) format() string {
	var buf bytes.Buffer
	height := len(m)
	width := 1

	for i := 0; i < height; i++ {
		prefix, _ := twoPowN(height - 1 - i)
		jump := 2 * prefix - 1
		for j := 0; j < width; j++ {
			if j == 0 {
				buf.WriteString(strings.Repeat(space, prefix))
			} else {
				buf.WriteString(strings.Repeat(space, jump))
			}
			e := m[i][j]
			buf.WriteString(e.format())
		}
		buf.WriteString("\n")
		width *= 2
	}

	return buf.String()
}

func twoPowN(n int) (int, error) {
	if n < 0 {
		return 0, fmt.Errorf("n < 0")
	}
	result := 1
	for i := 0; i < n; i++ {
		result *= 2
	}
	return result, nil
}

func (t *tree) bfs(f func(node *withPosition)) error {
	q := queue.NewQueue()
	q.Push(&withPosition{t, 0, 0})
	for !q.Empty() {
		elem, _ := q.Pop()
		node := elem.(*withPosition)
		f(node)
		if node.Value.left != nil {
			q.Push(&withPosition{node.Value.left, node.Level + 1, 2 * node.Index})
		}
		if node.Value.right != nil {
			q.Push(&withPosition{node.Value.right, node.Level + 1, 2 * node.Index + 1})
		}
	}
	return nil
}

type withPosition struct {
	Value *tree
	Level int
	Index int
}

func (p *withPosition) format() string {
	if p == nil {
		return space
	}
	format := fmt.Sprintf("%%%d.d", Size)
	return fmt.Sprintf(format, p.Value.value)
}

func Sort(values []int) {
	var root *tree
	for _, v := range values {
		root = add(root, v)
	}
	appendValues(values[:0], root)
}

func appendValues(values []int, t *tree) []int {
	if t != nil {
		values = appendValues(values, t.left)
		values = append(values, t.value)
		values = appendValues(values, t.right)
	}
	return values
}

func add(t *tree, value int) *tree {
	if t == nil {
		t = new(tree)
		t.value = value
		return t
	}
	if value < t.value {
		t.left = add(t.left, value)
	} else {
		t.right = add(t.right, value)
	}
	return t
}

func main() {
	arr := []int{4, 1, 2, 3}
	var t *tree
	for _, v := range arr {
		t = add(t, v)
	}
	fmt.Println(t.Stringy())
}
