package main

import (
	"fmt"
	"strings"
)

type treeNode struct {
	value  int
	parent *treeNode
	left   *treeNode
	right  *treeNode
}

func (n *treeNode) print() string {
	if n == nil {
		return ""
	}

	str := fmt.Sprintf("%d", n.value)
	children := ""
	children += "("
	children += n.left.print()
	children += ","
	children += n.right.print()
	children += ")"
	if children != "(,)" {
		str += children
	}
	return str
}

func (n *treeNode) add(parent *treeNode, value int) *treeNode {
	if n == nil {
		return &treeNode{
			parent: parent,
			value:  value,
		}
	}
	if value < n.value {
		n.left = n.left.add(n, value)
	} else {
		n.right = n.right.add(n, value)
	}
	return n
}

func (n *treeNode) hasChildren() bool {
	if n.left != nil {
		return true
	}
	if n.right != nil {
		return true
	}
	return false
}

func (n *treeNode) hasParent() bool {
	return n.parent != nil
}

func (n *treeNode) format(prefix string) (string, int) {
	str := prefix
	if n.hasParent() {
		str += "|-"
	}
	str += fmt.Sprintf("%d", n.value)
	if n.hasChildren() {
		str += "-|"
	}
	length := len(str) - len(prefix)
	if n.hasChildren() {
		length--
	}
	return str, length
}

func (n *treeNode) horizonPrint(prefix string) string {
	if n == nil {
		return ""
	}

	cur, length := n.format(prefix)
	str := ""
	str += n.right.horizonPrint(prefix + strings.Repeat(".", length))
	str += cur + "\n"
	str += n.left.horizonPrint(prefix + strings.Repeat(".", length))
	return str
}

type Tree struct {
	node *treeNode
}

func (t *Tree) Add(value int) {
	t.node = t.node.add(nil, value)
}

func (t *Tree) Print() string {
	return t.node.print()
}

func (t *Tree) HorizonPrint() string {
	return t.node.horizonPrint("")
}

func NewTree() *Tree {
	return &Tree{}
}

func main() {
	//10 8 5 7 12 4
	tree := NewTree()
	tree.Add(10)
	tree.Add(8)
	tree.Add(5)
	tree.Add(7)
	tree.Add(12)
	tree.Add(4)
	fmt.Println(tree.Print())
	fmt.Println(tree.HorizonPrint())
}
