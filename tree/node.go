package tree

func New[T any](data T, children ...Node[T]) Node[T] {
	return Node[T]{
		Data:     data,
		Children: children,
	}
}

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
