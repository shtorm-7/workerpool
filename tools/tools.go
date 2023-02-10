package tools

import "sync/atomic"

type TaskResult[R any] struct {
	Result R
	Err    error
}

func atomicIncrement(value *int32) {
	atomic.AddInt32(value, 1)
}

func atomicDecrementAndCloseIfZero[T any](value *int32, channel chan T) {
	if atomic.AddInt32(value, -1) == 0 {
		close(channel)
	}
}
