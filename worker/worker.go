package worker

import (
	"sync"

	"github.com/shtorm-7/workerpool/callbackfield"
	C "github.com/shtorm-7/workerpool/constant"
)

type Worker struct {
	queue C.Queue

	flow Flow

	status *callbackfield.CallbackField[C.Status]
	state  *callbackfield.CallbackField[C.State]

	meta C.Meta

	mtx sync.Mutex
}

func NewWorker(queue C.Queue, opts ...WorkerOption) *Worker {
	worker := &Worker{
		queue:  queue,
		flow:   DefaultFlow,
		status: callbackfield.NewCallbackField[C.Status](),
		state:  callbackfield.NewCallbackField[C.State](),
	}
	for _, opt := range opts {
		opt(worker)
	}
	return worker
}

func NewWorkerFactory(queue C.Queue, opts ...WorkerOption) C.WorkerFactory {
	return func() C.BaseWorker {
		return NewWorker(queue, opts...)
	}
}

func (w *Worker) Start() {
	w.mtx.Lock()
	defer w.mtx.Unlock()
	if w.status.Get() == C.Stopped {
		w.status.Set(C.Starting)
		go func() {
			w.status.Set(C.Started)
			w.flow(w)
			w.status.Set(C.Stopped)
		}()
		<-w.status.Await(C.Started)
	}
}

func (w *Worker) Stop() {
	w.mtx.Lock()
	defer w.mtx.Unlock()
	if w.status.Get() == C.Started {
		w.status.Set(C.Stopping)
		<-w.status.Await(C.Stopped)
	}
}

func (w *Worker) Status() *callbackfield.CallbackField[C.Status] {
	return w.status
}

func (w *Worker) State() *callbackfield.CallbackField[C.State] {
	return w.state
}

func (w *Worker) Meta() C.Meta {
	return w.meta
}

func (w *Worker) processTask(task func()) {
	w.state.Set(C.Received)
	task()
	w.state.Set(C.Succeeded)
	w.state.Set(C.Pending)
}
