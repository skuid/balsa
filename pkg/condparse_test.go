package parse

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	desc     string
	exprTree Node
	logic    string
}

var cases []testCase

func init() {
	cases = []testCase{
		{
			"Should serialize a single leaf",
			&Leaf{1},
			"1",
		},
		{
			"Should parse a very simple tree with just one operation and two leafs",
			&Op{
				Left:  &Leaf{1},
				Val:   "AND",
				Right: &Leaf{2},
			},
			"1 AND 2",
		},
		{
			"Should parse a more complex tree with left only operations (no parens)",
			&Op{
				Left: &Op{
					Left:  &Leaf{1},
					Val:   "OR",
					Right: &Leaf{2},
				},
				Val:   "AND",
				Right: &Leaf{3},
			},
			"1 OR 2 AND 3",
		},
		{
			"Should parse a more complex tree with left and right operations, including parens",
			&Op{
				Left: &Op{
					Left:  &Leaf{1},
					Val:   "OR",
					Right: &Leaf{2},
				},
				Val: "AND",
				Right: &Op{
					Left:  &Leaf{3},
					Val:   "OR",
					Right: &Leaf{4},
				},
			},
			"1 OR 2 AND (3 OR 4)",
		},
		{
			"Should complex trees with more depth",
			&Op{
				Left: &Op{
					Left: &Leaf{1},
					Val:  "OR",
					Right: &Op{
						Left: &Leaf{5},
						Val:  "AND",
						Right: &Op{
							Left:  &Leaf{7},
							Val:   "OR",
							Right: &Leaf{8},
						},
					},
				},
				Val: "AND",
				Right: &Op{
					Left: &Op{
						Left: &Op{
							Left:  &Leaf{3},
							Val:   "AND",
							Right: &Leaf{2},
						},
						Val: "OR",
						Right: &Op{
							Left:  &Leaf{56},
							Val:   "AND",
							Right: &Leaf{1000},
						},
					},
					Val:   "OR",
					Right: &Leaf{4},
				},
			},
			"1 OR (5 AND (7 OR 8)) AND (3 AND 2 OR (56 AND 1000) OR 4)",
		},
	}
}

func TestParse(t *testing.T) {
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			assert := assert.New(t)
			actual, err := Parse(c.logic)

			assert.NoError(err, "Should not have an error")
			deepEql(assert, c.exprTree, actual)
		})
	}
}

func TestSerialize(t *testing.T) {
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			assert := assert.New(t)
			var b strings.Builder
			c.exprTree.Eval(&b)

			actual := b.String()
			assert.Equal(c.logic, actual)
		})
	}
}

func deepEql(assert *assert.Assertions, expected Node, actual Node) {
	switch e := expected.(type) {
	case *Leaf:
		a, isLeaf := actual.(*Leaf)
		assert.True(isLeaf, "Expected was a leaf, actual was not! expected %v, actual %v", expected, actual)
		if e != nil && a != nil {
			assert.Equal(e.Val, a.Val, "Expected leaf values to match")
		} else if e != nil || a != nil {
			assert.Fail(fmt.Sprintf("One was nil and the other had a value. Expected %v, Actual %v", e, a))
		}
	case *Op:
		a, isOp := actual.(*Op)
		assert.True(isOp, "Expected was an Operation, actual was not! expected %v, actual %v", expected, actual)
		if e != nil && a != nil {
			assert.Equal(e.Val, a.Val, "Expected operation to be the same")
			deepEql(assert, e.Left, a.Left)
			deepEql(assert, e.Right, a.Right)
		} else if e != nil || a != nil {
			assert.Fail(fmt.Sprintf("One was nil and the other had a value. Expected %v, Actual %v", e, a))
		}
	default:
		assert.Fail("Node was neither a leaf or an op.")
	}
}
