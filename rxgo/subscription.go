package rxgo

type Subscription struct {
	unsubscribeChannel chan<- bool
}

func (s *Subscription) Unsubscribe() {
	s.unsubscribeChannel <- true
	//close(s.unsubscribeChannel)
}
