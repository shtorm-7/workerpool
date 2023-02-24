package tools

import (
	"time"

	C "github.com/shtorm-7/workerpool/constant"
)

func Await(queue C.Queue, task func()) <-chan struct{} {
	await := make(chan struct{})
	queue <- func() {
		task()
		close(await)
	}
	return await
}

func TryAwait(queue C.Queue, task func()) (<-chan struct{}, bool) {
	await := make(chan struct{})
	select {
	case queue <- func() {
		task()
		close(await)
	}:
		return await, true
	default:
		return nil, false
	}
}

func TryAwaitWithTimeout(queue C.Queue, task func(), timeout time.Duration) (<-chan struct{}, bool) {
	await := make(chan struct{})
	select {
	case queue <- func() {
		task()
		close(await)
	}:
		return await, true
	case <-time.After(timeout):
		return nil, false
	}
}
