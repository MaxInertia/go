package rxgo

func NewBehaviorSubject[T any](v T) BehaviorSubject[T] {
	stream := make(chan T, 1)
	stream <- v
	return BehaviorSubject[T]{
		lastValue: v,
		Observable: Observable[T]{
			stream: stream,
		},
	}
}

type BehaviorSubject[T any] struct {
	lastValue T
	Observable[T]
}

func (s BehaviorSubject[T]) readFromStream() T {
	v := <-s.stream
	s.lastValue = v
	return v
}

func (s BehaviorSubject[T]) ToObservable() Observable[T] {
	return s.Observable
}

func (s BehaviorSubject[T]) Emit(v T) {
	s.lastValue = v
	s.stream <- v
}
