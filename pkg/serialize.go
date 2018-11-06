package parse

import (
	"fmt"
	"io"
)

// Eval will print the leaf's value to a writer
func (l *Leaf) Eval(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%d", l.Val); err != nil {
		return err
	}
	return nil
}

// Eval will print the left node, the operation, and then the right node to a writer
func (o *Op) Eval(w io.Writer) error {
	if o.Left == nil {
		return &SerializeError{
			Op:     o.Val,
			Reason: "nil left node",
		}
	}

	if o.Right == nil {
		return &SerializeError{
			Op:     o.Val,
			Reason: "nil right node",
		}
	}

	if !(o.Val == "AND" || o.Val == "OR") {
		return &SerializeError{
			Op:     o.Val,
			Reason: "bad operation",
		}
	}

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

// SerializeError holds information about syntactic errors when trying to eval
type SerializeError struct {
	Op     string
	Reason string
}

func (e *SerializeError) Error() string {
	return fmt.Sprintf("Could not serialize operation '%s'. Reason: %s", e.Op, e.Reason)
}
