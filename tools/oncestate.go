package tools

import "sync/atomic"

type onceState[T any] struct {
	channel chan T
	state   atomic.Uint32
}

func newOnceState[T any](channel chan T) *onceState[T] {
	onceState := onceState[T]{channel: channel}
	onceState.state.Store(1)
	return &onceState
}

func (s *onceState[T]) Add(delta int) uint32 {
	return s.state.Add(uint32(delta))
}

func (s *onceState[T]) Done() {
	if s.Add(-1) == 0 {
		close(s.channel)
	}
}
