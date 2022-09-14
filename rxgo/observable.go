package rxgo

type observable[T any] interface {
	Subscribe(omEmit func(t T)) Subscription
}

type Observable[T any] struct {
	observable[T]
	stream chan T
}

func (o Observable[T]) Subscribe(omEmit func(t T)) Subscription {
	unsubscribe := make(chan bool)
	go func() {
		for {
			select {
			case <-unsubscribe:
				return
			case v := <-o.stream:
				omEmit(v)
			}
		}
	}()
	return Subscription{
		unsubscribeChannel: unsubscribe,
	}
}

func Map[T any, R any](fn func(t T) R) func(o Observable[T]) Observable[R] {
	return func(o Observable[T]) Observable[R] {
		modifiedStream := make(chan R)
		go func() {
			select {
			case v := <-o.stream:
				modifiedStream <- fn(v)
			}
		}()
		return Observable[R]{
			stream: modifiedStream,
		}
	}
}
