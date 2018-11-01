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

// Eval will print the leaf's value to a writer
func (l *Leaf) Eval(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%d", l.Val); err != nil {
		return err
	}
	return nil
}

// Eval will print the left node, the operation, and then the right node to a writer
func (o *Op) Eval(w io.Writer) error {
	if err := o.Left.Eval(w); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(w, " %s ", o.Val); err != nil {
		return err
	}

	_, parens := o.Right.(*Op)

	if parens {
		fmt.Fprint(w, "(")
	}

	if err := o.Right.Eval(w); err != nil {
		return err
	}

	if parens {
		fmt.Fprint(w, ")")
	}

	return nil
}

// Remove a node by value
//	n.Remove(1)
// Removes any leaf that has a value of 1, shifting up the tree where needed
func (l *Leaf) Remove(v uint) Node {
	if l.Val == v {
		return nil
	}
	return l
}

// Remove a node by value from an operation
//	n.Remove(1)
// Removes any leaf that has a value of 1, shifting up the tree where needed
func (o *Op) Remove(v uint) Node {

	l := o.Left.Remove(v)
	r := o.Right.Remove(v)

	if l == nil && r == nil {
		return nil
	}

	if l == nil {
		return r
	}

	if r == nil {
		return l
	}

	return &Op{
		Left:  l,
		Val:   o.Val,
		Right: r,
	}
}
