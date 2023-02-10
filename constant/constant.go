package constant

import "github.com/shtorm-7/workerpool/callbackfield"

type (
	BaseWorker interface {
		Start()
		Stop()
		Status() *callbackfield.CallbackField[Status]
		Meta() Meta
	}

	WorkerFactory func() BaseWorker
)
