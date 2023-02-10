package pool

import C "github.com/shtorm-7/workerpool/constant"

type (
	WorkerHandler func(worker C.BaseWorker)

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

func WithMeta(meta C.Meta) PoolOption {
	return func(pool *Pool) {
		pool.meta = meta
	}
}

func WithInitWorkerHandler(handler WorkerHandler) PoolOption {
	return func(pool *Pool) {
		for _, worker := range pool.workers {
			handler(worker)
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

func WithAddWorkerHandler(handler WorkerHandler) ResizablePoolOption {
	return func(pool *ResizablePool) {
		pool.addWorkerHandlers = append(pool.addWorkerHandlers, handler)
	}
}

func WithRemoveWorkerHandler(handler WorkerHandler) ResizablePoolOption {
	return func(pool *ResizablePool) {
		pool.removeWorkerHandlers = append(pool.removeWorkerHandlers, handler)
	}
}
