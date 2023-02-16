package pool

import (
	"sync"

	"github.com/shtorm-7/workerpool/callbackfield"
	C "github.com/shtorm-7/workerpool/constant"
)

type Pool struct {
	workers []C.BaseWorker

	startHandler StartHandler
	stopHandler  StopHandler

	status *callbackfield.CallbackField[C.Status]

	metricHandlers []C.MetricHandler

	mtx sync.Mutex
}

func NewPool(factories []C.WorkerFactory, opts ...PoolOption) C.Pool {
	if len(factories) == 0 {
		panic("factories cant be blank")
	}
	pool := &Pool{
		workers:      make([]C.BaseWorker, len(factories)),
		status:       callbackfield.NewCallbackField[C.Status](),
		startHandler: SequentialStart,
		stopHandler:  SequentialStop,
	}
	for i, factory := range factories {
		pool.workers[i] = factory()
	}
	for _, opt := range opts {
		opt(pool)
	}
	return pool
}

func NewPoolFactory(factories []C.WorkerFactory, opts ...PoolOption) C.WorkerFactory {
	return func() C.BaseWorker {
		return NewPool(factories, opts...)
	}
}

func (p *Pool) Start() {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	if p.status.Get() == C.Stopped {
		p.status.Set(C.Starting)
		p.startHandler(p.workers)
		p.status.Set(C.Started)
	}
}

func (p *Pool) Stop() {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	if p.status.Get() == C.Started {
		p.status.Set(C.Stopping)
		p.stopHandler(p.workers)
		p.status.Set(C.Stopped)
	}
}

func (p *Pool) Status() callbackfield.ReadOnlyCallbackField[C.Status] {
	return p.status
}

func (p *Pool) Workers() []C.BaseWorker {
	return p.workers
}

func (p *Pool) Metrics() C.Metrics {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	metrics := make(C.Metrics, len(p.metricHandlers)+1)
	for _, handler := range p.metricHandlers {
		metricName, metric := handler()
		metrics[metricName] = metric
	}
	workerMetrics := make([]C.Metrics, len(p.workers))
	for i, worker := range p.workers {
		workerMetrics[i] = worker.Metrics()
	}
	metrics["workers"] = workerMetrics
	return metrics
}
