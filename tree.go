package main

import (
	"fmt"
	"sort"
	"strings"
)

type Tree map[string]Tree
type Type int

func (tree Tree) Add(path string) {
	fragments := strings.Split(path, "/")

	if fragments[len(fragments)-1] == "" {
		fragments = fragments[:len(fragments)-1]
	}

	fragments = append([]string{"/"}, fragments...)
	tree.add(fragments)
}

func (tree Tree) Print() {
	tree.print(true, "")
}

func (tree Tree) add(fragments []string) {
	if len(fragments) == 0 {
		return
	}

	nextTree, ok := tree[fragments[0]]

	if !ok {
		nextTree = Tree{}
		tree[fragments[0]] = nextTree
	}

	nextTree.add(fragments[1:])
}

func (tree Tree) print(root bool, padding string) {
	var index = 0
	var keys []string

	for k := range tree {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("%s%s\n", padding+getPadding(root, getType(index, len(tree))), k)

		tree[k].print(false, padding+getPadding(root, getTypeExternal(index, len(tree))))

		index++
	}
}

const (
	regular Type = iota
	last
	afterLast
	between
)

func (t Type) String() string {
	switch t {
	case regular:
		return CyanString("├── ")
	case last:
		return CyanString("└── ")
	case afterLast:
		return CyanString("    ")
	case between:
		return CyanString("│   ")
	default:
		panic("Invalid type")
	}
}

func getType(index, length int) Type {
	if index+1 == length {
		return last
	} else if index+1 > length {
		return afterLast
	}

	return regular
}

func getTypeExternal(index, length int) Type {
	if index+1 == length {
		return afterLast
	}

	return between
}

func getPadding(root bool, t Type) string {
	if root {
		return ""
	}

	return t.String()
}
