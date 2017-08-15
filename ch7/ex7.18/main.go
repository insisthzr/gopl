package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/insisthzr/happytool/container/stack"
)

var (
	ErrInvalidSyntax = errors.New("invalid syntax")
)

type Node interface{} // CharData or *Element

type CharData string

type Element struct {
	Type     xml.Name
	Attr     []xml.Attr
	Children []Node
}

func (e *Element) String() string {
	return visit(e, 0)
}

func visit(n Node, depth int) string {
	str := ""
	switch n := n.(type) {
	case *Element:
		prefix := strings.Repeat(" ", 2*depth)
		str += fmt.Sprintf("%s%s\n", prefix, n.Type.Local)
		for _, child := range n.Children {
			str += visit(child, depth+1)
		}
	case CharData:
		prefix := strings.Repeat(" ", 2*depth)
		str += fmt.Sprintf("%s%q\n", prefix, n)
	default:
		panic("impossible")
	}
	return str
}

func parse(reader io.Reader) (interface{}, error) {
	dec := xml.NewDecoder(reader)
	stac := stack.NewStack()
	var root Node

	for {
		token, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch token := token.(type) {
		case xml.StartElement:
			elem := &Element{Type: token.Name, Attr: token.Attr}
			if stac.Len() == 0 {
				root = elem
			} else {
				e, _ := stac.Top()
				parent := e.(*Element)
				parent.Children = append(parent.Children, elem)
			}
			stac.Push(elem)
		case xml.EndElement:
			stac.Pop()
		case xml.CharData:
			//skip lf space
			if strings.TrimSpace(string(token)) == "" {
				continue
			}
			e, _ := stac.Top()
			parent := e.(*Element)
			parent.Children = append(parent.Children, CharData(token))
		}
	}

	return root, nil
}

func main() {
	root, err := parse(os.Stdin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(root)
}
