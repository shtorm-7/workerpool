package pool

import (
	"fmt"

	"github.com/shtorm-7/workerpool/callbackfield"
	C "github.com/shtorm-7/workerpool/constant"
)

type ResizablePool struct {
	*Pool

	factory C.WorkerFactory

	postAddWorkerHandlers    []WorkerHandler
	postRemoveWorkerHandlers []WorkerHandler
}

func NewResizablePool(factory C.WorkerFactory, n int, opts ...ResizablePoolOption) *ResizablePool {
	if n <= 0 {
		panic(fmt.Sprintf("the value '%d' is not valid. the value must be greater than 0", n))
	}
	pool := &ResizablePool{
		Pool: &Pool{
			workers:      make([]C.BaseWorker, n),
			status:       callbackfield.NewCallbackField[C.Status](),
			startHandler: ParallelStart,
			stopHandler:  ParallelStop,
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

func NewResizablePoolFactory(factory C.WorkerFactory, n int, opts ...ResizablePoolOption) C.WorkerFactory {
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
		p.workers = append(p.workers, worker)
		for _, handler := range p.postAddWorkerHandlers {
			handler(worker)
		}
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
	removingWorkers := p.workers[len(p.workers)-n:]
	if p.status.Get() == C.Started {
		p.stopHandler(removingWorkers)
	}
	p.workers = p.workers[:len(p.workers)-n]
	for _, worker := range removingWorkers {
		for _, handler := range p.postRemoveWorkerHandlers {
			handler(worker)
		}
	}
}
