package parse

import (
	"golang.org/x/tools/container/intsets"
)

// Index Op will index the left and right nodes
func (o *Op) Index(start int) Node {

	var left Node
	var right Node

	if o.Left != nil {
		left = o.Left.Index(start)
	}

	if o.Right != nil {
		right = o.Right.Index(start)
	}

	return &Op{
		Left:  left,
		Val:   o.Val,
		Right: right,
	}
}

// Index Leaf will add start to the current value
func (l *Leaf) Index(start int) Node {
	return &Leaf{
		Val: uint(int(l.Val) + start),
	}
}

// Sequence will walk all of the leaves and re-sequence the values, starting
// at 0 and removing sparseness
func Sequence(n Node) Node {
	// collect numbers from all leafs into an array

	leafs := &intsets.Sparse{}
	addSet := initAddSet(leafs)
	WalkLeaves(n, addSet)

	larr := make([]int, 0, leafs.Len())
	larr = leafs.AppendTo(larr)

	cmap := map[int]int{}
	for i, v := range larr {
		cmap[v] = i
	}

	setVal := initSetVal(cmap)
	return WalkLeaves(n, setVal)
}

// Visitor visits a node, allowing you to take action on a node
type Visitor func(Node) Node

// WalkLeaves will visit every leaf and run the Visitor action on each
func WalkLeaves(n Node, visit Visitor) Node {
	if n == nil {
		return n
	}
	switch node := n.(type) {
	case *Leaf:
		return visit(node)
	case *Op:
		return &Op{
			Left:  WalkLeaves(node.Left, visit),
			Val:   node.Val,
			Right: WalkLeaves(node.Right, visit),
		}
	}
	return n
}

func initAddSet(s *intsets.Sparse) Visitor {
	return func(n Node) Node {
		if l, ok := n.(*Leaf); ok {
			s.Insert(int(l.Val))
		}
		return n
	}
}

func initSetVal(cmap map[int]int) Visitor {
	return func(n Node) Node {
		if l, ok := n.(*Leaf); ok {
			if v, ok := cmap[int(l.Val)]; ok {
				return &Leaf{uint(v)}
			}
		}
		return n
	}
}
