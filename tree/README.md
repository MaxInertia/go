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
```

Where `op` is a function that is applied to each node.

### Traversal Operators

#### Accumulate

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

#### Emit

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

