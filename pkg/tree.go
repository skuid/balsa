package condparse

import (
	"fmt"
	"io"
)

// Node is our interface and represents a node in a binary expression tree.
type Node interface {
	Remove(uint) Node
	Eval(io.Writer) error
}

// Leaf is a concrete Node that will hold a single value
type Leaf struct {
	Val uint
}

// Op is a concrete Node that will hold an operation with pointers to a left
// and right node, both of which can be a Leaf or another Op
type Op struct {
	Left  Node
	Val   string
	Right Node
}

// String is our stringer for pretty printing the tree. Well, ok, it is ugly
// printing, but it is printing. It will take a tree and print something like:
// 1 <- AND -> 2 <- OR -> 3
func (l *Leaf) String() string {
	return fmt.Sprintf("%d", l.Val)
}

// String is our stringer for pretty printing the tree. Well, ok, it is ugly
// printing, but it is printing. It will take a tree and print something like:
// 1 <- AND -> 2 <- OR -> 3
func (o *Op) String() string {
	return fmt.Sprintf("%v <- %s -> %v", o.Left, o.Val, o.Right)
}
