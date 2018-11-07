package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseErrors(t *testing.T) {
	cases := []struct {
		desc    string
		fixture string
		err     *ParseError
	}{
		{
			"Should fail with only an operation",
			"AND",
			&ParseError{
				Position: 0,
				Logic:    "AND",
				Reason:   "unexpected operation",
			},
		},
		{
			"Should fail with an unacceptable operation",
			"1 FOO 3",
			&ParseError{
				Position: 2,
				Logic:    "1 FOO 3",
				Reason:   "FOO is an unacceptable operation",
			},
		},
		{
			"Should fail with two operations in a row",
			"1 AND (OR 2)",
			&ParseError{
				Position: 7,
				Logic:    "1 AND (OR 2)",
				Reason:   "unexpected operation",
			},
		},
		{
			"Should fail with a bad character",
			"1 ! AND 2",
			&ParseError{
				Position: 2,
				Logic:    "1 ! AND 2",
				Reason:   "general error",
			},
		},
		{
			"Should fail with an unexpected number",
			"1 AN132D 2",
			&ParseError{
				Position: 4,
				Logic:    "1 AN132D 2",
				Reason:   "unexpected number",
			},
		},
		{
			"Should fail with an unexpected character at the end",
			"1 AND 2AA",
			&ParseError{
				Position: 7,
				Logic:    "1 AND 2AA",
				Reason:   "unexpected character",
			},
		},
		{
			"Should fail with an unexpected character",
			"1A AND 2AA",
			&ParseError{
				Position: 1,
				Logic:    "1A AND 2AA",
				Reason:   "unexpected character",
			},
		},
		{
			"Should fail with unbalanced parens",
			"1 AND 2 (3 AND 4",
			&ParseError{
				Position: 16,
				Logic:    "1 AND 2 (3 AND 4",
				Reason:   "unbalanced parenthesis",
			},
		},
		{
			"Should fail with a closing parens without an opening parens",
			"1 AND 2 )3 AND 4",
			&ParseError{
				Position: 8,
				Logic:    "1 AND 2 )3 AND 4",
				Reason:   "unexpected closing parenthesis",
			},
		},
		{
			"Should fail with too many closing parens",
			"1 AND 2 OR (3 AND 4) OR 5) AND 6",
			&ParseError{
				Position: 25,
				Logic:    "1 AND 2 OR (3 AND 4) OR 5) AND 6",
				Reason:   "unexpected closing parenthesis",
			},
		},
		{
			"Should fail with an unexpected opening parens",
			"1 (3 AND 4)",
			&ParseError{
				Position: 2,
				Logic:    "1 (3 AND 4)",
				Reason:   "unexpected opening parenthesis",
			},
		},
		{
			"Should fail with an unexpected opening parens",
			"1 AND 3 (2 AND 3)",
			&ParseError{
				Position: 16,
				Logic:    "1 AND 3 (2 AND 3)",
				Reason:   "invalid syntax",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			assert := assert.New(t)
			_, err := Parse(c.fixture)
			assert.Equal(c.err, err)
		})
	}
}

func TestValidButOddParse(t *testing.T) {
	cases := []struct {
		desc     string
		fixture  string
		expected Node
	}{
		{
			"Should handle multiple spaces",
			"    1 AND    2      ",
			&Op{
				Left:  &Leaf{1},
				Val:   "AND",
				Right: &Leaf{2},
			},
		},
		{
			"Should handle extra parens on the left",
			"((1 AND 2) AND 5) OR (3 AND 4)",
			&Op{
				Left: &Op{
					Left: &Op{
						Left:  &Leaf{1},
						Val:   "AND",
						Right: &Leaf{2},
					},
					Val:   "AND",
					Right: &Leaf{5},
				},
				Val: "OR",
				Right: &Op{
					Left:  &Leaf{3},
					Val:   "AND",
					Right: &Leaf{4},
				},
			},
		},
		{
			"Should handle extra parens with or without spaces",
			"1 AND 2 AND 5 OR(  3 AND 4)",
			&Op{
				Left: &Op{
					Left: &Op{
						Left:  &Leaf{1},
						Val:   "AND",
						Right: &Leaf{2},
					},
					Val:   "AND",
					Right: &Leaf{5},
				},
				Val: "OR",
				Right: &Op{
					Left:  &Leaf{3},
					Val:   "AND",
					Right: &Leaf{4},
				},
			},
		},
		{
			"Should handle extraneous parens",
			"1 AND 2 AND 5 OR (   ((3 AND 4))   )",
			&Op{
				Left: &Op{
					Left: &Op{
						Left:  &Leaf{1},
						Val:   "AND",
						Right: &Leaf{2},
					},
					Val:   "AND",
					Right: &Leaf{5},
				},
				Val: "OR",
				Right: &Op{
					Left:  &Leaf{3},
					Val:   "AND",
					Right: &Leaf{4},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			assert := assert.New(t)
			actual, err := Parse(c.fixture)

			assert.NoError(err, "Should not have an error")
			deepEql(assert, c.expected, actual)
		})
	}
}
