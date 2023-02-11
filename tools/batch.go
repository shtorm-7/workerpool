package tools

import (
	C "github.com/shtorm-7/workerpool/constant"
	"github.com/shtorm-7/workerpool/generator"
)

func AwaitBatch(queue C.Queue, tasks generator.Scheme[func()]) <-chan struct{} {
	await := make(chan struct{})
	go func() {
		state := newOnceState(await)
		tasks.Process(
			func(task func()) {
				state.Add(1)
				queue <- func() {
					task()
					state.Done()
				}
			},
		)
		state.Done()
	}()
	return await
}

func Batch[R any](queue C.Queue, resultsSize int, tasks generator.Scheme[func() (R, error)]) <-chan TaskResult[R] {
	results := make(chan TaskResult[R], resultsSize)
	go func() {
		state := newOnceState(results)
		tasks.Process(
			func(task func() (R, error)) {
				state.Add(1)
				queue <- func() {
					result, err := task()
					results <- TaskResult[R]{result, err}
					state.Done()
				}
			},
		)
		state.Done()
	}()
	return results
}
