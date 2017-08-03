package main

import (
	"fmt"
	"strings"
)

type nodeStatut int

const (
	unvisited nodeStatut = iota
	visiting
	visited
)

func toposort(m map[string][]string) (result []string, e error) {
	defer func() {
		if r := recover(); r != nil {
			switch r := r.(type) {
			case error:
				e = r
			default:
				e = fmt.Errorf("%v", r)
			}
		}
	}()
	order := []string{}
	resolved := map[string]nodeStatut{}
	var visitAll func([]string, []string)
	visitAll = func(items []string, visitings []string) {
		for _, item := range items {
			if resolved[item] == visiting {
				start := findIndexInSlice(visitings, item)
				if start == -1 {
					panic("impossible")
				}
				paths := strings.Join(append(visitings[start:], item), " -> ")
				panic(fmt.Errorf("cycle: %s", paths))
			}
			if resolved[item] == unvisited {
				resolved[item] = visiting
				visitAll(m[item], append(visitings, item))
				resolved[item] = visited
				order = append(order, item)
			}
		}
	}
	for key := range m {
		visitAll([]string{key}, []string{})
	}
	return order, nil
}

func findIndexInSlice(slice []string, value string) int {
	for i, v := range slice {
		if v == value {
			return i
		}
	}
	return -1
}

// prereqs maps computer science courses to their prerequisites.
var prereqs = map[string][]string{
	"algorithms":     {"data structures"},
	"calculus":       {"linear algebra"},
	"linear algebra": {"calculus"},
	// uncomment line below to introduce a dependency cycle.
	"intro to programming": {"data structures"},

	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},

	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}

func main() {
	order, err := toposort(prereqs)
	fmt.Println(order, err)
}
