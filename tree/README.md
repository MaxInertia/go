# Tree

Another tree implementation the world didn't need.

## Creating trees

```go
// Root -> Leaf
New("Root", New("Leaf"))

// Left <- Root -> Right
New("Root", New("Left"), New("Right"))
```

## Traversals

To iterate over all nodes you can use a traversal.
```go
DepthFirstTraversal(root, op)
// or
BreadthFirstTraversal(root, op)
```

Where `op` is a function that is applied to each node.

## Traversal Operators

### Accumulate

To accumulate the traversed nodes in a slice you can use `Accumulate`.
```go
left := New("Left")
right := New("Right")
root := New("Root", left, right)

nodes := make([]Node[string], 0)
DepthFirstTraversal(root, Accumulate[string](nodes))
nodes[0] // == root
nodes[1] // == left
nodes[2] // == right
```

### Emit

To emit traversed nodes on a channel you can use `Emit`.
```go
left := New("Left")
right := New("Right")
root := New("Root", left, right)

// Run traversal in another goroutine to avoid blocking
channel := make(chan Node[string])
go DepthFirstTraversal(root, Emit[string](channel))
<- channel // == root
<- channel // == left
<- channel // == right

// Or have a buffer >= the number of nodes in the tree
channel := make(chan Node[string], 4)
DepthFirstTraversal(root, Emit[string](channel))
<- channel // == root
<- channel // == left
<- channel // == right
```

## Interface `tree`

You can implement `tree` on your own types to utilize the functions in this package on them.

```go
type tree[T any, Children any] interface {
	GetData() T
	GetChildren() []Children
}
```

### Example 1: `Node`
Exposed by this package [here](https://github.com/MaxInertia/go/blob/main/tree/node.go)
```go
type Node[T any] struct {
	Data     T
	Children []Node[T]
}

func (n Node[T]) GetData() T {
	return n.Data
}

func (n Node[T]) GetChildren() []Node[T] {
	return n.Children
}
```

### Example 2: `SliceTree`
See [slice_test.go](https://github.com/MaxInertia/go/blob/main/tree/slice_test.go)

```go
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
```
