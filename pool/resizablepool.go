package pool

import (
	"fmt"

	"github.com/shtorm-7/workerpool/callbackfield"
	C "github.com/shtorm-7/workerpool/constant"
)

type ResizablePool struct {
	*Pool

	factory C.WorkerFactory

	addWorkerHandlers    []WorkerHandler
	removeWorkerHandlers []WorkerHandler
}

func NewResizablePool(
	factory C.WorkerFactory,
	n int,
	opts ...ResizablePoolOption,
) *ResizablePool {
	if n <= 0 {
		panic(fmt.Sprintf("the value '%d' is not valid. the value must be greater than 0", n))
	}
	pool := &ResizablePool{
		Pool: &Pool{
			workers:      make([]C.BaseWorker, n),
			status:       callbackfield.NewCallbackField[C.Status](),
			startHandler: ConcurrentStart,
			stopHandler:  ConcurrentStop,
		},
		factory: factory,
	}
	for i := range pool.workers {
		pool.workers[i] = factory()
	}
	for _, opt := range opts {
		opt(pool)
	}
	return pool
}

func NewResizablePoolFactory(
	factory C.WorkerFactory,
	n int,
	opts ...ResizablePoolOption,
) C.WorkerFactory {
	return func() C.BaseWorker {
		return NewResizablePool(factory, n, opts...)
	}
}

func (p *ResizablePool) AddWorkers(n int) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	if n <= 0 {
		panic(fmt.Sprintf("the value '%d' is not valid. the value must be greater than 0", n))
	}
	for i := 0; i < n; i++ {
		worker := p.factory()
		for _, handler := range p.addWorkerHandlers {
			handler(worker)
		}
		p.workers = append(p.workers, worker)
	}
	if p.status.Get() == C.Started {
		p.startHandler(p.workers[len(p.workers)-n:])
	}
}

func (p *ResizablePool) RemoveWorkers(n int) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	if n <= 0 {
		panic(fmt.Sprintf("the value '%d' is not valid. the value must be greater than 0", n))
	} else if n > len(p.workers) {
		panic(fmt.Sprintf("the value '%d' is not valid. the value must be less than length of current workers", n))
	}
	if p.status.Get() == C.Started {
		p.stopHandler(p.workers[len(p.workers)-n:])
	}
	for _, worker := range p.workers[len(p.workers)-n:] {
		for _, handler := range p.removeWorkerHandlers {
			handler(worker)
		}
	}
	p.workers = p.workers[:len(p.workers)-n]
}
