package tree

type tree[T any, Children any] interface {
	GetData() T
	GetChildren() []Children
}
