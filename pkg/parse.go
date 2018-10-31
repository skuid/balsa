package condparse

import (
	"bytes"
	"fmt"
	"strconv"
	"unicode"
)

type Token int

const (
	// Special tokens
	NIL Token = iota
	WS
	OP
	LEAF
)

func isWhitespace(c rune) bool {
	return c == ' ' || c == '\t'
}

func isChar(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

type parser struct {
	buffer bytes.Buffer
	kind   Token
	tree   Node
}

func (p *parser) Reset() {
	p.kind = NIL
	p.buffer.Reset()
}

func (p *parser) Eval(pos int) error {

	if p.kind == LEAF {
		i, err := strconv.ParseUint(p.buffer.String(), 10, 0)
		if err != nil {
			return fmt.Errorf("Not able to convert %s at position %d to unit", p.buffer.String(), pos)
		}
		current := &Leaf{uint(i)}

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
	}

	if p.kind == OP {
		op := p.buffer.String()

		if op != "AND" && op != "OR" {
			return fmt.Errorf("Operation '%s' was not an acceptable operation at location %d", op, pos)
		}

		if p.tree == nil {
			return fmt.Errorf("Found an operation in a bad location at %d", pos)
		}

		if t, ok := p.tree.(*Op); ok {
			if t.Left == nil {
				return fmt.Errorf("Found an operation in a bad location at %d", pos)
			} else if t.Val != "" && t.Right == nil {
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
	}

	return nil
}

func Parse(logic string) (Node, error) {

	var err error
	var p parser

	for i, r := range logic {

		if unicode.IsSpace(r) {
			// Store off buffer into tree we're building, and reset the buffer
			err = p.Eval(i)
			if err != nil {
				return nil, err
			}
			p.Reset()
		} else if unicode.IsNumber(r) {
			// first check to make sure we're not started or we're on a number
			if !(p.kind == NIL || p.kind == LEAF) {
				return nil, fmt.Errorf("Found an unexpected number at position %d", i)
			}
			// start buffering a number
			p.kind = LEAF
			p.buffer.WriteRune(r)
		} else if unicode.IsLetter(r) {
			if !(p.kind == NIL || p.kind == OP) {
				return nil, fmt.Errorf("Found an unexpected character at position %d", i)
			}
			// start buffering a string
			p.kind = OP
			p.buffer.WriteRune(r)
			// } else if c == '(' {
			// 	// Start an expression, but we may need to write out last buffer.

			// } else if c == ')' {
			// 	// end an expression, attach it to parent

		} else {
			// throw exception
			return nil, fmt.Errorf("There was a general parsing error at position %d", i)
		}

	}

	err = p.Eval(len(logic))
	if err != nil {
		return nil, err
	}

	return p.tree, nil
}
