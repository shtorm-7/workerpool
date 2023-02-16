package tools

import (
	C "github.com/shtorm-7/workerpool/constant"
	"github.com/shtorm-7/workerpool/generator"
)

type (
	ChainResult[FR any] struct {
		Result FR
		Err    error
	}

	LinkingHandler[FR, V any] func(value V, finalResultHandler func(FR, error))

	Chain[FR, V any] struct {
		rootHandler LinkingHandler[FR, V]
	}

	Link[CFR, CV, V, R any] struct {
		queue C.Queue

		handler     func(V) (R, error)
		nextHandler LinkingHandler[CFR, R]

		chain *Chain[CFR, CV]
	}
)

func NewChain[CFR, CV, V any](link *Link[CFR, CV, V, CFR]) *Chain[CFR, CV] {
	link.nextHandler = func(value CFR, finalResultHandler func(CFR, error)) {
		finalResultHandler(value, nil)
	}
	return link.chain
}

func NewLink[CFR, CV, LR any](queue C.Queue, handler func(CV) (LR, error)) *Link[CFR, CV, CV, LR] {
	link := &Link[CFR, CV, CV, LR]{
		queue:   queue,
		handler: handler,
		chain:   new(Chain[CFR, CV]),
	}
	link.chain.rootHandler = link.linkingHandler()
	return link
}

func AddLink[CFR, CV, V, PR, R any](previousLink *Link[CFR, CV, V, PR], queue C.Queue, handler func(PR) (R, error)) *Link[CFR, CV, PR, R] {
	link := &Link[CFR, CV, PR, R]{
		queue:   queue,
		handler: handler,
		chain:   previousLink.chain,
	}
	previousLink.nextHandler = link.linkingHandler()
	return link
}

func (ch *Chain[FR, V]) Await(value V) <-chan struct{} {
	await := make(chan struct{})
	go func() {
		ch.rootHandler(
			value,
			func(FR, error) {
				close(await)
			},
		)
	}()
	return await
}

func (ch *Chain[FR, V]) Future(value V) <-chan ChainResult[FR] {
	results := make(chan ChainResult[FR])
	go func() {
		ch.rootHandler(
			value,
			func(finalResult FR, err error) {
				results <- ChainResult[FR]{finalResult, err}
			},
		)
	}()
	return results
}

func (ch *Chain[FR, V]) AwaitBatch(values generator.Generator[V]) <-chan struct{} {
	await := make(chan struct{})
	go func() {
		state := newOnceState(await)
		defer state.Done()
		finalResultHandler := func(FR, error) {
			state.Done()
		}
		values.Process(
			func(value V) {
				state.Add(1)
				ch.rootHandler(value, finalResultHandler)
			},
		)
	}()
	return await
}

func (ch *Chain[FR, V]) Batch(resultsSize int, values generator.Generator[V]) <-chan ChainResult[FR] {
	results := make(chan ChainResult[FR], resultsSize)
	go func() {
		state := newOnceState(results)
		defer state.Done()
		finalResultHandler := func(finalResult FR, err error) {
			results <- ChainResult[FR]{finalResult, err}
			state.Done()
		}
		values.Process(
			func(value V) {
				state.Add(1)
				ch.rootHandler(value, finalResultHandler)
			},
		)
	}()
	return results
}

func (link *Link[CFR, CV, V, R]) linkingHandler() LinkingHandler[CFR, V] {
	return func(value V, finalResultHandler func(CFR, error)) {
		link.queue <- func() {
			if result, err := link.handler(value); err == nil {
				link.nextHandler(result, finalResultHandler)
			} else {
				finalResultHandler(*new(CFR), err)
			}
		}
	}
}
