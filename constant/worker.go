package constant

import "github.com/shtorm-7/workerpool/callbackfield"

type (
	Queue chan func()

	Worker interface {
		BaseWorker

		State() callbackfield.ReadOnlyCallbackField[State]
	}
)
