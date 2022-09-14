package tree

// Flatten returns the tree as a slice with the nodes in depth first order
func Flatten[T any, N tree[T, N]](root N) []N {
	var nodes = make([]N, 0)
	DepthFirstTraversal[T](root, Accumulate[T](&nodes))
	return nodes
}
