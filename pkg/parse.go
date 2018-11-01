package condparse

import (
	"bytes"
	"fmt"
	"strconv"
	"unicode"
)

type token int

const (
	NIL token = iota
	OP
	LEAF
)

type parser struct {
	buffer bytes.Buffer
	kind   token
	tree   Node
	expr   []Node
}

func (p *parser) reset() {
	p.kind = NIL
	p.buffer.Reset()
}

func (p *parser) open(pos int) error {

	// We need to eval/flush. Do that now.
	if p.kind != NIL {
		p.eval(pos)
	}

	// tree should be an Op. If it is a LEAF, this is a bad state
	_, ok := p.tree.(*Op)
	if !ok {
		return fmt.Errorf("There shouldn't be an opening ( after a non-operation at poisition %d", pos)
	}

	// Store off the current tree into the expression stack. We'll pop it back out on close
	p.expr = append(p.expr, p.tree)
	p.tree = nil

	return nil
}

func (p *parser) close(pos int) error {
	if len(p.expr) <= 0 {
		return fmt.Errorf("Unbalanced parens at position %d", pos)
	}

	// Process what is already in the buffer
	if p.kind != NIL {
		p.eval(pos)
	}

	// Pop off the top expression
	var e Node
	e, p.expr = pop(p.expr)

	t, ok := e.(*Op)
	if !ok {
		return fmt.Errorf("Expression was not an operation, so not able to add the sub expression at position %d", pos)
	}

	// Put the current tree onto our popped expression
	if t.Left == nil {
		t.Left = p.tree
	} else if t.Right == nil {
		t.Right = p.tree
	} else {
		return fmt.Errorf("The parent operation already had both a left and a right at position %d", pos)
	}

	// Make our current tree our popped expression plus the last tree evaluated
	p.tree = t
	return nil
}

func (p *parser) procLeaf(pos int) error {
	i, err := strconv.ParseUint(p.buffer.String(), 10, 0)
	if err != nil {
		return fmt.Errorf("Not able to convert %s at position %d to unit", p.buffer.String(), pos)
	}

	// Create the current leaf from what was in the buffer
	current := &Leaf{uint(i)}

	// If we don't have a tree yet, start it with this leaf.
	// Otherwise, figure out where it needs to go, shifting around as needed
	if p.tree == nil {
		p.tree = current
	} else {
		if t, ok := p.tree.(*Op); ok {
			if t.Val == "" && t.Left == nil {
				t.Left = current
			} else if t.Val != "" && t.Right == nil {
				t.Right = current
			} else {
				return fmt.Errorf("Got a leaf at position %d, but other conditions were not met", pos)
			}
		} else if _, ok := p.tree.(*Leaf); ok {
			return fmt.Errorf("We're no a leaf and got another leaf at position %d", pos)
		}
	}

	return nil
}

func (p *parser) procOp(pos int) error {
	op := p.buffer.String()

	// TODO: We probably want a better way to do this that's easier to expand
	if op != "AND" && op != "OR" {
		return fmt.Errorf("Operation '%s' was not an acceptable operation at location %d", op, pos)
	}

	// This could happen if the first characters scanned were an op and not a number
	if p.tree == nil {
		return fmt.Errorf("Found an operation in a bad location at %d", pos)
	}

	// If the current tree is holding an operation already, we just need to set
	// its value to the current operation. There should already be a left node,
	// as we're scanning left to right.
	// If the tree is a Leaf, we need to put the leaf on our left and set the
	// value to the current operation
	if t, ok := p.tree.(*Op); ok {
		if t.Left == nil {
			return fmt.Errorf("Found an operation in a bad location at %d", pos)
		} else if t.Val != "" && t.Right == nil {
			// Left has a value, right does not, just set the operation
			t.Val = op
		} else {
			p.tree = &Op{
				Left: p.tree,
				Val:  op,
			}
		}
	} else if l, ok := p.tree.(*Leaf); ok {
		p.tree = &Op{
			Left: l,
			Val:  op,
		}
	}

	return nil
}

func (p *parser) eval(pos int) error {

	if p.kind == LEAF {
		if err := p.procLeaf(pos); err != nil {
			return err
		}
	}

	if p.kind == OP {
		if err := p.procOp(pos); err != nil {
			return err
		}
	}

	p.reset()
	return nil
}

// Parse will take a logic string (e.g. "1 AND 2 OR (3 AND 4)"), parse it and
// turn it into a binary expression tree.
func Parse(logic string) (Node, error) {

	var p parser

	for i, r := range logic {

		if unicode.IsSpace(r) {
			// Store off buffer into tree we're building, and reset the buffer
			if err := p.eval(i); err != nil {
				return nil, err
			}
		} else if unicode.IsNumber(r) {
			// first check to make sure we're not started or we're on a number
			if !(p.kind == NIL || p.kind == LEAF) {
				return nil, fmt.Errorf("Found an unexpected number at position %d", i)
			}
			// start buffering a number
			p.kind = LEAF
			p.buffer.WriteRune(r)
		} else if unicode.IsLetter(r) {
			// first check to make sure we're not started or we're already working on a word
			if !(p.kind == NIL || p.kind == OP) {
				return nil, fmt.Errorf("Found an unexpected character at position %d", i)
			}
			// start buffering a string
			p.kind = OP
			p.buffer.WriteRune(r)
		} else if r == '(' {
			// Start an expression, but we may need to write out last buffer.
			if err := p.open(i); err != nil {
				return nil, err
			}
		} else if r == ')' {
			// end an expression, attach it to parent
			if err := p.close(i); err != nil {
				return nil, err
			}
		} else {
			// throw exception
			return nil, fmt.Errorf("There was a general parsing error at position %d", i)
		}

	}

	if err := p.eval(len(logic)); err != nil {
		return nil, err
	}

	return p.tree, nil
}

func pop(slice []Node) (Node, []Node) {
	return slice[len(slice)-1], slice[:len(slice)-1]
}
