package rxgo

import (
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestBehaviorSubject(t *testing.T) {
	t.Run("Subscribe emits initial value", func(t *testing.T) {
		subject := NewBehaviorSubject("1")
		capture := make(chan string)
		_ = subject.ToObservable().Subscribe(func(v string) {
			capture <- v
		})
		require.Equal(t, "1", <-capture)
	})
	t.Run("Subscribe before emit", func(t *testing.T) {
		subject := NewBehaviorSubject("1")

		capture := make(chan string)
		_ = subject.Subscribe(func(v string) {
			capture <- v
		})

		require.Equal(t, "1", <-capture)
		subject.Emit("2")
		require.Equal(t, "2", <-capture)
	})
	t.Run("Subscribe before emit - with unsubscribe", func(t *testing.T) {
		subject := NewBehaviorSubject("1")

		capture := make(chan string)
		sub := subject.Subscribe(func(v string) {
			capture <- v
		})

		require.Equal(t, "1", <-capture)
		subject.Emit("2")
		require.Equal(t, "2", <-capture)

		sub.Unsubscribe()
	})
	t.Run("Pipe.map emitted values", func(t *testing.T) {
		subject := NewBehaviorSubject("1")
		obs := Map[string, int](func(t string) int {
			i, _ := strconv.ParseInt(t, 10, 64)
			return int(i)
		})(subject.ToObservable())

		capture := make(chan int)
		sub := obs.Subscribe(func(v int) {
			capture <- v
		})

		require.Equal(t, 1, <-capture)
		subject.Emit("2")
		require.Equal(t, 2, <-capture)

		sub.Unsubscribe()
	})
	t.Run("Subscribe after emit", func(t *testing.T) {
		t.Skip()
		subject := NewBehaviorSubject("1")
		subject.Emit("2")
		capture := make(chan string)
		subject.ToObservable().Subscribe(func(v string) {
			capture <- v
		})
		require.Equal(t, "2", <-capture)
	})
	t.Run("multiple emits before Subscribe", func(t *testing.T) {
		t.Skip()
		subject := NewBehaviorSubject("1")
		capture := make(chan string)
		subject.Subscribe(func(v string) {
			capture <- v
		})
		// Capture must emit twice
		subject.Emit("1")
		subject.Emit("2")
		require.Equal(t, "1", <-capture)
		require.Equal(t, "2", <-capture)
	})
}
