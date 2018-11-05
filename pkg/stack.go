package condparse

func pop(slice []Node) (Node, []Node) {
	return slice[len(slice)-1], slice[:len(slice)-1]
}
