package parse

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerializeErrors(t *testing.T) {
	cases := []struct {
		desc    string
		fixture Node
		err     *SerializeError
	}{
		{
			"Should return an exception if the operation has a nil left",
			&Op{
				Val:   "AND",
				Right: &Leaf{3},
			},
			&SerializeError{
				Op:     "AND",
				Reason: "nil left node",
			},
		},
		{
			"Should return an exception if the operation has a nil right",
			&Op{
				Val:  "AND",
				Left: &Leaf{3},
			},
			&SerializeError{
				Op:     "AND",
				Reason: "nil right node",
			},
		},
		{
			"Should return an exception if the operation doesn't have a value",
			&Op{
				Left:  &Leaf{3},
				Val:   "",
				Right: &Leaf{3},
			},
			&SerializeError{
				Op:     "",
				Reason: "bad operation",
			},
		},
		{
			"Should return an exception if the operation doesn't have an acceptable value",
			&Op{
				Left:  &Leaf{3},
				Val:   "FOO",
				Right: &Leaf{3},
			},
			&SerializeError{
				Op:     "FOO",
				Reason: "bad operation",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			assert := assert.New(t)
			var b strings.Builder
			err := c.fixture.Eval(&b)
			assert.Equal(c.err, err)
		})
	}
}
