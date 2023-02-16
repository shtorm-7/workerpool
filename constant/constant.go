package constant

import "github.com/shtorm-7/workerpool/callbackfield"

type (
	BaseWorker interface {
		Start()
		Stop()
		Status() callbackfield.ReadOnlyCallbackField[Status]
		Metrics() Metrics
	}

	WorkerFactory func() BaseWorker
)
