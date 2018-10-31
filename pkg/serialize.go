package condparse

import (
	"fmt"
	"io"
)

type Node interface {
	Remove(uint) Node
	Eval(io.Writer) error
}

type Leaf struct {
	Val uint
}

func (l *Leaf) Eval(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%d", l.Val); err != nil {
		return err
	}
	return nil
}

type Op struct {
	Left  Node
	Val   string
	Right Node
}

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

func (l *Leaf) Remove(v uint) Node {
	if l.Val == v {
		return nil
	}
	return l
}

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
