package pool

import C "github.com/shtorm-7/workerpool/constant"

type (
	WorkersHandler func(workers []C.BaseWorker)

	PoolOption func(pool *Pool)

	ResizablePoolOption func(pool *ResizablePool)
)

func WithStartHandler(startHandler StartHandler) PoolOption {
	return func(pool *Pool) {
		pool.startHandler = startHandler
	}
}

func WithStopHandler(stopHandler StopHandler) PoolOption {
	return func(pool *Pool) {
		pool.stopHandler = stopHandler
	}
}

func WithMetrics(factories ...C.MetricHandlerFactory[C.Pool]) PoolOption {
	return func(pool *Pool) {
		for _, factory := range factories {
			pool.metricHandlers = append(pool.metricHandlers, factory(pool))
		}
	}
}

func WithPostInitHandlers(handlers ...WorkersHandler) PoolOption {
	return func(pool *Pool) {
		for _, handler := range handlers {
			handler(pool.workers)
		}
	}
}

func WithPoolOptions(opts ...PoolOption) ResizablePoolOption {
	return func(pool *ResizablePool) {
		for _, opt := range opts {
			opt(pool.Pool)
		}
	}
}

func WithPostAddHandlers(handlers ...WorkersHandler) ResizablePoolOption {
	return func(pool *ResizablePool) {
		pool.postAddHandlers = append(pool.postAddHandlers, handlers...)
	}
}

func WithPostRemoveHandlers(handlers ...WorkersHandler) ResizablePoolOption {
	return func(pool *ResizablePool) {
		pool.postRemoveHandlers = append(pool.postRemoveHandlers, handlers...)
	}
}
