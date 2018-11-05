package condparse

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemove(t *testing.T) {
	cases := []struct {
		desc     string
		fixture  Node
		remove   []uint
		expected string
	}{
		{
			"Should not remove anything if not found",
			&Leaf{1},
			[]uint{3},
			"1",
		},
		{
			"Should remove a single leaf",
			&Leaf{1},
			[]uint{1},
			"",
		},
		{
			"Should parse a very simple tree with just one operation and two leafs",
			&Op{
				Left:  &Leaf{1},
				Val:   "AND",
				Right: &Leaf{2},
			},
			[]uint{2},
			"1",
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
			[]uint{2},
			"1 AND 3",
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
			[]uint{3},
			"1 OR 2 AND 4",
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
			[]uint{7, 56},
			"1 OR (5 AND 8) AND (3 AND 2 OR 1000 OR 4)",
		},
		{
			"Should remove the value from multiple places",
			&Op{
				Left: &Op{
					Left: &Leaf{1},
					Val:  "OR",
					Right: &Op{
						Left: &Leaf{5},
						Val:  "AND",
						Right: &Op{
							Left:  &Leaf{1},
							Val:   "OR",
							Right: &Leaf{1},
						},
					},
				},
				Val: "AND",
				Right: &Op{
					Left: &Op{
						Left: &Op{
							Left:  &Leaf{1},
							Val:   "AND",
							Right: &Leaf{2},
						},
						Val: "OR",
						Right: &Op{
							Left:  &Leaf{56},
							Val:   "AND",
							Right: &Leaf{1},
						},
					},
					Val:   "OR",
					Right: &Leaf{4},
				},
			},
			[]uint{1},
			"5 AND (2 OR 56 OR 4)",
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			assert := assert.New(t)

			n := c.fixture

			for _, v := range c.remove {
				n = n.Remove(v)
			}

			var b strings.Builder
			if n != nil {
				n.Eval(&b)
				actual := b.String()
				assert.Equal(c.expected, actual)
			} else {
				assert.Equal(c.expected, "")
			}
		})
	}
}
