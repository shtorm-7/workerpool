package tools

import (
	C "github.com/shtorm-7/workerpool/constant"
	"github.com/shtorm-7/workerpool/generator"
)

type (
	LinkingHandler[V, FR any] func(value V, finalResultHandler func(FR, error))

	Chain[V, FR any] struct {
		linkingHandler LinkingHandler[V, FR]
	}

	Link[CV, CFR, V, R any] struct {
		queue C.Queue

		handler            func(V) (R, error)
		nextLinkingHandler LinkingHandler[R, CFR]

		chain *Chain[CV, CFR]
	}
)

func NewChain[CV, CFR, V any](link *Link[CV, CFR, V, CFR]) *Chain[CV, CFR] {
	link.nextLinkingHandler = func(value CFR, finalResultHandler func(CFR, error)) {
		finalResultHandler(value, nil)
	}
	return link.chain
}

func NewLink[CV, CFR, LR any](queue C.Queue, handler func(CV) (LR, error)) *Link[CV, CFR, CV, LR] {
	link := &Link[CV, CFR, CV, LR]{
		queue:   queue,
		handler: handler,
		chain:   new(Chain[CV, CFR]),
	}
	link.chain.linkingHandler = link.linkingHandler()
	return link
}

func AddLink[CV, CFR, V, PR, R any](previousLink *Link[CV, CFR, V, PR], queue C.Queue, handler func(PR) (R, error)) *Link[CV, CFR, PR, R] {
	link := &Link[CV, CFR, PR, R]{
		queue:   queue,
		handler: handler,
		chain:   previousLink.chain,
	}
	previousLink.nextLinkingHandler = link.linkingHandler()
	return link
}

func (ch Chain[V, FR]) AwaitBatch(values generator.Scheme[V]) <-chan struct{} {
	await := make(chan struct{})
	go func() {
		state := int32(1)
		finalResultHandler := func(FR, error) {
			atomicDecrementAndCloseIfZero(&state, await)
		}
		values.Process(
			func(value V) {
				atomicIncrement(&state)
				ch.linkingHandler(value, finalResultHandler)
			},
		)
		atomicDecrementAndCloseIfZero(&state, await)
	}()
	return await
}

func (ch Chain[V, FR]) Batch(resultsSize int, values generator.Scheme[V]) <-chan TaskResult[FR] {
	results := make(chan TaskResult[FR], resultsSize)
	go func() {
		state := int32(1)
		finalResultHandler := func(finalResult FR, err error) {
			results <- TaskResult[FR]{finalResult, err}
			atomicDecrementAndCloseIfZero(&state, results)
		}
		values.Process(
			func(value V) {
				atomicIncrement(&state)
				ch.linkingHandler(value, finalResultHandler)
			},
		)
		atomicDecrementAndCloseIfZero(&state, results)
	}()
	return results
}

func (link *Link[CV, CFR, V, R]) linkingHandler() LinkingHandler[V, CFR] {
	return func(value V, finalResultHandler func(CFR, error)) {
		link.queue <- func() {
			if result, err := link.handler(value); err == nil {
				link.nextLinkingHandler(result, finalResultHandler)
			} else {
				finalResultHandler(*new(CFR), err)
			}
		}
	}
}
