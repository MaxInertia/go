package tree_test

import (
	. "generics/tree"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("Can create tree: A->B->C", func(t *testing.T) {
		actual := New("A", New("B", New("C")))
		expected := Node[string]{
			Data: "A",
			Children: []Node[string]{
				{
					Data: "B",
					Children: []Node[string]{
						{
							Data: "C",
						},
					},
				},
			},
		}
		require.Equal(t, expected, actual)
	})
}

func TestNode_GetChildren(t *testing.T) {
	a := New("A")
	ba := New("B", a)
	require.Equal(t, ba.GetChildren()[0], a)
	require.Len(t, a.GetChildren(), 0)
}

func Test_Flatten(t *testing.T) {
	a := New("A")
	ba := New("B", a)
	cba := New("C", ba)
	t.Run("(1->2)=[1,2]", func(t *testing.T) {
		nodes := Flatten[string, Node[string]](ba)
		require.Equal(t, nodes[0], ba)
		require.Equal(t, nodes[1], a)
	})

	t.Run("(1->2->3)=[1,2,3]", func(t *testing.T) {
		nodes := Flatten[string, Node[string]](cba)
		require.Equal(t, nodes[0], cba)
		require.Equal(t, nodes[1], ba)
		require.Equal(t, nodes[2], a)
	})

	t.Run("(Left<-Root->Right)=[Root,Left,Right]", func(t *testing.T) {
		left := New("Left")
		right := New("Right")
		root := New("Root", left, right)
		nodes := Flatten[string, Node[string]](root)
		require.Equal(t, nodes[0], root)
		require.Equal(t, nodes[1], left)
		require.Equal(t, nodes[2], right)
	})
}

func TestDepthFirstTraversal(t *testing.T) {
	t.Run("Accumulate", func(t *testing.T) {
		left := New("Left")
		right := New("Right")
		root := New("Root", left, right)
		nodes := make([]Node[string], 0)
		DepthFirstTraversal[string](root, Accumulate[string](&nodes))
		require.Len(t, nodes, 3)
		require.Equal(t, nodes[0], root)
		require.Equal(t, nodes[1], left)
		require.Equal(t, nodes[2], right)
	})

	t.Run("Emit without buffer", func(t *testing.T) {
		left := New("Left")
		right := New("Right")
		root := New("Root", left, right)
		nodes := make(chan Node[string])
		go DepthFirstTraversal[string](root, Emit[string](nodes))
		n0 := <-nodes
		require.Equal(t, n0, root)
		n1 := <-nodes
		require.Equal(t, n1, left)
		n2 := <-nodes
		require.Equal(t, n2, right)
	})

	t.Run("Emit with buffer", func(t *testing.T) {
		left := New("Left")
		right := New("Right")
		root := New("Root", left, right)
		nodes := make(chan Node[string], 3)
		DepthFirstTraversal[string](root, Emit[string](nodes))
		n0 := <-nodes
		require.Equal(t, n0, root)
		n1 := <-nodes
		require.Equal(t, n1, left)
		n2 := <-nodes
		require.Equal(t, n2, right)
	})
}
