package parse

// Index Op will index the left and right nodes
func (o *Op) Index(start uint) Node {

	var left Node
	var right Node

	if o.Left != nil {
		left = o.Left.Index(start)
	}

	if o.Right != nil {
		right = o.Right.Index(start)
	}

	return &Op{
		Left:  left,
		Val:   o.Val,
		Right: right,
	}
}

// Index Leaf will add start to the current value
func (l *Leaf) Index(start uint) Node {
	return &Leaf{
		Val: l.Val + start,
	}
}

func Sequence(n Node) (Node, error) {
	// collect numbers from all leafs into an array

	// order the numbers

	// give each number a value a new ordinal starting with 1

	// return the result

	return n, nil
}
