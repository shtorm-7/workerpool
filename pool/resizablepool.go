package pool

import (
	"fmt"

	"github.com/shtorm-7/workerpool/callbackfield"
	C "github.com/shtorm-7/workerpool/constant"
)

type ResizablePool struct {
	*Pool

	factory C.WorkerFactory

	postAddHandlers    []WorkersHandler
	postRemoveHandlers []WorkersHandler
}

func NewResizablePool(factory C.WorkerFactory, n int, opts ...ResizablePoolOption) C.ResizablePool {
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
	addingWorkers := make([]C.BaseWorker, n)
	for i := range addingWorkers {
		addingWorkers[i] = p.factory()
	}
	p.workers = append(p.workers, addingWorkers...)
	for _, handler := range p.postAddHandlers {
		handler(addingWorkers)
	}
	if p.status.Get() == C.Started {
		p.startHandler(addingWorkers)
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
	p.workers = p.workers[:len(p.workers)-n]
	if p.status.Get() == C.Started {
		p.stopHandler(removingWorkers)
	}
	for _, handler := range p.postRemoveHandlers {
		handler(removingWorkers)
	}
}
