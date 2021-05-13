package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	cases := []struct {
		desc     string
		fixture  Node
		start    int
		expected Node
	}{
		{
			"Should reindex a single leaf",
			&Leaf{1},
			5,
			&Leaf{6},
		},
		{
			"Should reindex a simple op",
			&Op{
				Left:  &Leaf{2},
				Val:   "AND",
				Right: &Leaf{5},
			},
			21,
			&Op{
				Left:  &Leaf{23},
				Val:   "AND",
				Right: &Leaf{26},
			},
		},
		{
			"Should reindex a complex op",
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
			3,
			&Op{
				Left: &Op{
					Left: &Leaf{4},
					Val:  "OR",
					Right: &Op{
						Left: &Leaf{8},
						Val:  "AND",
						Right: &Op{
							Left:  &Leaf{10},
							Val:   "OR",
							Right: &Leaf{11},
						},
					},
				},
				Val: "AND",
				Right: &Op{
					Left: &Op{
						Left: &Op{
							Left:  &Leaf{6},
							Val:   "AND",
							Right: &Leaf{5},
						},
						Val: "OR",
						Right: &Op{
							Left:  &Leaf{59},
							Val:   "AND",
							Right: &Leaf{1003},
						},
					},
					Val:   "OR",
					Right: &Leaf{7},
				},
			},
		},
		{
			"Should handle nil leafs",
			&Op{
				Left: &Leaf{3},
				Val:  "OR",
			},
			5,
			&Op{
				Left: &Leaf{8},
				Val:  "OR",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			assert := assert.New(t)
			actual := c.fixture.Index(c.start)

			deepEql(assert, c.expected, actual)
		})
	}
}

func TestSequence(t *testing.T) {
	cases := []struct {
		desc     string
		fixture  Node
		expected Node
	}{
		{
			"Should sequence a single leaf",
			&Leaf{5},
			&Leaf{0},
		},
		{
			"Should sequence a simple operation",
			&Op{
				Left:  &Leaf{5},
				Val:   "AND",
				Right: &Leaf{3},
			},
			&Op{
				Left:  &Leaf{1},
				Val:   "AND",
				Right: &Leaf{0},
			},
		},
		{
			"Should reindex a complex op",
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
			&Op{
				Left: &Op{
					Left: &Leaf{0},
					Val:  "OR",
					Right: &Op{
						Left: &Leaf{4},
						Val:  "AND",
						Right: &Op{
							Left:  &Leaf{5},
							Val:   "OR",
							Right: &Leaf{6},
						},
					},
				},
				Val: "AND",
				Right: &Op{
					Left: &Op{
						Left: &Op{
							Left:  &Leaf{2},
							Val:   "AND",
							Right: &Leaf{1},
						},
						Val: "OR",
						Right: &Op{
							Left:  &Leaf{7},
							Val:   "AND",
							Right: &Leaf{8},
						},
					},
					Val:   "OR",
					Right: &Leaf{3},
				},
			},
		},
		{
			"Should handle nil leafs",
			&Op{
				Left: &Leaf{3},
				Val:  "OR",
			},
			&Op{
				Left: &Leaf{0},
				Val:  "OR",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			assert := assert.New(t)

			actual := Sequence(c.fixture)

			deepEql(assert, c.expected, actual)
		})
	}
}
