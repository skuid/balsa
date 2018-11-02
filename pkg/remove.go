package condparse

// Remove a node by value
//	n.Remove(1)
// Removes any leaf that has a value of 1, shifting up the tree where needed
func (l *Leaf) Remove(v uint) Node {
	if l.Val == v {
		return nil
	}
	return l
}

// Remove a node by value from an operation
//	n.Remove(1)
// Removes any leaf that has a value of 1, shifting up the tree where needed
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
