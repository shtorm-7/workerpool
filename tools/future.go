package tools

import (
	"time"

	C "github.com/shtorm-7/workerpool/constant"
)

func Future[R any](queue C.Queue, task func() R) <-chan R {
	result := make(chan R)
	queue <- func() {
		result <- task()
		close(result)
	}
	return result
}

func TryFuture[R any](queue C.Queue, task func() R) (<-chan R, bool) {
	result := make(chan R)
	select {
	case queue <- func() {
		result <- task()
		close(result)
	}:
		return result, true
	default:
		return nil, false
	}
}

func TryFutureWithTimeout[R any](queue C.Queue, task func() R, timeout time.Duration) (<-chan R, bool) {
	result := make(chan R)
	select {
	case queue <- func() {
		result <- task()
		close(result)
	}:
		return result, true
	case <-time.After(timeout):
		return nil, false
	}
}
