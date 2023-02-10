package tools

import (
	"time"

	C "github.com/shtorm-7/workerpool/constant"
)

func Future[R any](queue C.Queue, task func() (R, error)) <-chan TaskResult[R] {
	results := make(chan TaskResult[R])
	go func() {
		queue <- func() {
			result, err := task()
			results <- TaskResult[R]{result, err}
			close(results)
		}
	}()
	return results
}

func TryFuture[R any](queue C.Queue, task func() (R, error)) (<-chan TaskResult[R], bool) {
	results := make(chan TaskResult[R])
	select {
	case queue <- func() {
		result, err := task()
		results <- TaskResult[R]{result, err}
		close(results)
	}:
		return results, true
	default:
		return nil, false
	}
}

func TryFutureWithTimeout[R any](queue C.Queue, task func() (R, error), timeout time.Duration) (<-chan TaskResult[R], bool) {
	results := make(chan TaskResult[R])
	select {
	case queue <- func() {
		result, err := task()
		results <- TaskResult[R]{result, err}
		close(results)
	}:
		return results, true
	case <-time.After(timeout):
		return nil, false
	}
}
