package tree_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	. "tree"
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
	t.Run("Accumulate with one node", func(t *testing.T) {
		root := New("Root")
		nodes := make([]Node[string], 0)
		DepthFirstTraversal[string](root, Accumulate[string](&nodes))
		require.Len(t, nodes, 1)
		require.Equal(t, nodes[0], root)
	})

	t.Run("Accumulate 1<-0->2", func(t *testing.T) {
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

	t.Run("Accumulate 2<-1<-0->3", func(t *testing.T) {
		leftLeft := New(2)
		right := New(3)
		left := New(1, leftLeft)
		root := New(0, left, right)

		nodes := make([]Node[int], 0)

		DepthFirstTraversal[int](root, Accumulate[int](&nodes))
		require.Len(t, nodes, 4)
		assert.Equal(t, nodes[0].Data, 0)
		assert.Equal(t, nodes[1].Data, 1)
		assert.Equal(t, nodes[2].Data, 2)
		assert.Equal(t, nodes[3].Data, 3)
	})

	t.Run("Emit without buffer 2<-1<-0->3", func(t *testing.T) {
		root := New(0, New(1, New(2)), New(3))
		nodes := make(chan Node[int])
		go DepthFirstTraversal[int](root, Emit[int](nodes))
		n := <-nodes
		require.Equal(t, n.GetData(), 0)
		n = <-nodes
		require.Equal(t, n.GetData(), 1)
		n = <-nodes
		require.Equal(t, n.GetData(), 2)
		n = <-nodes
		require.Equal(t, n.GetData(), 3)
	})

	t.Run("Emit with buffer 2<-1<-0->3", func(t *testing.T) {
		root := New(0, New(1, New(2)), New(3))
		nodes := make(chan Node[int], 4)
		DepthFirstTraversal[int](root, Emit[int](nodes))
		n := <-nodes
		require.Equal(t, n.GetData(), 0)
		n = <-nodes
		require.Equal(t, n.GetData(), 1)
		n = <-nodes
		require.Equal(t, n.GetData(), 2)
		n = <-nodes
		require.Equal(t, n.GetData(), 3)
	})
}

func TestBreadthFirstTraversal(t *testing.T) {
	t.Run("Accumulate with one node", func(t *testing.T) {
		root := New("Root")
		nodes := make([]Node[string], 0)
		BreadthFirstTraversal[string](root, Accumulate[string](&nodes))
		require.Len(t, nodes, 1)
		require.Equal(t, nodes[0], root)
	})

	t.Run("Accumulate 1<-0->2", func(t *testing.T) {
		right := New(2)
		left := New(1)
		root := New(0, left, right)

		nodes := make([]Node[int], 0)

		BreadthFirstTraversal[int](root, Accumulate[int](&nodes))
		require.Len(t, nodes, 3)
		assert.Equal(t, nodes[0].Data, 0)
		assert.Equal(t, nodes[1].Data, 1)
		assert.Equal(t, nodes[2].Data, 2)
	})

	t.Run("Accumulate 3<-1<-0->2", func(t *testing.T) {
		root := New(0, New(1, New(3)), New(2))
		nodes := make([]Node[int], 0)
		BreadthFirstTraversal[int](root, Accumulate[int](&nodes))
		require.Len(t, nodes, 4)
		assert.Equal(t, nodes[0].Data, 0)
		assert.Equal(t, nodes[1].Data, 1)
		assert.Equal(t, nodes[2].Data, 2)
		assert.Equal(t, nodes[3].Data, 3)
	})

	t.Run("Emit without buffer 3<-1<-0->2", func(t *testing.T) {
		root := New(0, New(1, New(3)), New(2))
		nodes := make(chan Node[int])
		go BreadthFirstTraversal[int](root, Emit[int](nodes))
		n := <-nodes
		require.Equal(t, n.GetData(), 0)
		n = <-nodes
		require.Equal(t, n.GetData(), 1)
		n = <-nodes
		require.Equal(t, n.GetData(), 2)
		n = <-nodes
		require.Equal(t, n.GetData(), 3)
	})

	t.Run("Emit with buffer 3<-1<-0->2", func(t *testing.T) {
		leftLeft := New(3)
		left := New(1, leftLeft)
		right := New(2)
		root := New(0, left, right)
		nodes := make(chan Node[int], 4)
		BreadthFirstTraversal[int](root, Emit[int](nodes))
		n := <-nodes
		require.Equal(t, n.GetData(), 0)
		n = <-nodes
		require.Equal(t, n.GetData(), 1)
		n = <-nodes
		require.Equal(t, n.GetData(), 2)
		n = <-nodes
		require.Equal(t, n.GetData(), 3)
	})
}
