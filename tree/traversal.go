package tree

// DepthFirstTraversal applies the provided function on the nodes in depth first order
func DepthFirstTraversal[T any, N tree[T, N]](root N, op TraversalOperator[T, N]) {
	op(root)
	for _, child := range root.GetChildren() {
		DepthFirstTraversal[T](child, op)
	}
}

type TraversalOperator[T any, N tree[T, N]] func(n N)

// Accumulate generates a traversal operator that appends nodes to the given slice
func Accumulate[T any, N tree[T, N]](acc *[]N) TraversalOperator[T, N] {
	return func(n N) {
		*acc = append(*acc, n)
	}
}

// Emit generates a traversal operator sends nodes on a channel
func Emit[T any, N tree[T, N]](acc chan<- N) TraversalOperator[T, N] {
	return func(n N) {
		acc <- n
	}
}
