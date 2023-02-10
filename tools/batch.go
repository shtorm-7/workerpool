package tools

import (
	C "github.com/shtorm-7/workerpool/constant"
	"github.com/shtorm-7/workerpool/generator"
)

func AwaitBatch(queue C.Queue, tasks generator.Scheme[func()]) <-chan struct{} {
	await := make(chan struct{})
	go func() {
		state := int32(1)
		tasks.Process(
			func(task func()) {
				atomicIncrement(&state)
				queue <- func() {
					task()
					atomicDecrementAndCloseIfZero(&state, await)
				}
			},
		)
		atomicDecrementAndCloseIfZero(&state, await)
	}()
	return await
}

func Batch[R any](queue C.Queue, resultsSize int, tasks generator.Scheme[func() (R, error)]) <-chan TaskResult[R] {
	results := make(chan TaskResult[R], resultsSize)
	go func() {
		state := int32(1)
		tasks.Process(
			func(task func() (R, error)) {
				atomicIncrement(&state)
				queue <- func() {
					result, err := task()
					results <- TaskResult[R]{result, err}
					atomicDecrementAndCloseIfZero(&state, results)
				}
			},
		)
		atomicDecrementAndCloseIfZero(&state, results)
	}()
	return results
}
