package tree_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	. "tree"
)

// SliceTree is an example of how one could implement the tree
// interface to utilize methods in this package on other types.
//
// SliceTree is a tree where each node has at most 1 child.
// Example: [1,2] is a node with child [2]
type SliceTree[T any] []T

func (s SliceTree[T]) GetChildren() []SliceTree[T] {
	if len(s) == 0 {
		return nil
	}
	return []SliceTree[T]{s[1:]}
}

func (s SliceTree[T]) GetData() T {
	return s[0]
}

func Test_SliceAsTree(t *testing.T) {
	t.Run("Accumulate", func(t *testing.T) {
		var root = SliceTree[int]{1, 2, 3}
		nodes := make([]SliceTree[int], 0)
		DepthFirstTraversal[int](root, Accumulate[int, SliceTree[int]](&nodes))

		assert.Equal(t, 1, nodes[0][0])
		assert.Equal(t, 2, nodes[1][0])
		assert.Equal(t, 3, nodes[2][0])
	})
	t.Run("Emit", func(t *testing.T) {
		var root = SliceTree[int]{1, 2, 3}
		nodes := make(chan SliceTree[int])
		go DepthFirstTraversal[int](root, Emit[int, SliceTree[int]](nodes))

		assert.Equal(t, 1, (<-nodes)[0])
		assert.Equal(t, 2, (<-nodes)[0])
		assert.Equal(t, 3, (<-nodes)[0])
	})
}
