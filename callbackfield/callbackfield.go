package callbackfield

type (
	Callback func()

	OrAwait[T comparable] struct {
		values []T
		await  chan struct{}
	}

	CallbackField[T comparable] struct {
		value T

		anyAwait chan struct{}
		awaits   map[T]chan struct{}
		orAwaits []OrAwait[T]

		callbacks     map[T][]Callback
		onceCallbacks map[T][]Callback
		anyCallbacks  []Callback

		mtx RWMutex
	}

	ReadOnlyCallbackField[T comparable] interface {
		Get() T
		Await(values ...T) <-chan struct{}
		AddCallback(value T, callback Callback)
		AddOnceCallback(value T, callback Callback)
		AddAnyCallback(callback Callback)
	}
)

func NewCallbackField[T comparable](opts ...CallbackFieldOption[T]) *CallbackField[T] {
	field := &CallbackField[T]{
		awaits: make(map[T]chan struct{}),

		callbacks:     make(map[T][]Callback),
		onceCallbacks: make(map[T][]Callback),
	}
	for _, opt := range opts {
		opt(field)
	}
	return field
}

func (cf *CallbackField[T]) Get() T {
	cf.mtx.RLock()
	defer cf.mtx.RUnlock()
	return cf.value
}

func (cf *CallbackField[T]) Set(value T) {
	cf.mtx.Lock()
	if value == cf.value {
		cf.mtx.Unlock()
		return
	}
	cf.value = value
	cf.mtx.UnlockToRLock()
	defer cf.mtx.RUnlock()
	cf.callback()
}

func (cf *CallbackField[T]) Await(values ...T) <-chan struct{} {
	cf.mtx.RLock()
	if len(values) == 0 {
		cf.mtx.RUnlockToLock()
		defer cf.mtx.Unlock()
		if cf.anyAwait == nil {
			cf.anyAwait = make(chan struct{})
		}
		return cf.anyAwait
	} else if len(values) == 1 {
		value := values[0]
		if value == cf.value {
			cf.mtx.RUnlock()
			return ClosedChannel
		}
		cf.mtx.RUnlockToLock()
		defer cf.mtx.Unlock()
		if await, ok := cf.awaits[value]; ok {
			return await
		}
		await := make(chan struct{})
		cf.awaits[value] = await
		return await
	} else {
		for _, value := range values {
			if value == cf.value {
				cf.mtx.RUnlock()
				return ClosedChannel
			}
		}
		cf.mtx.RUnlockToLock()
		defer cf.mtx.Unlock()
		orAwait := OrAwait[T]{values, make(chan struct{})}
		cf.orAwaits = append(cf.orAwaits, orAwait)
		return orAwait.await
	}
}

func (cf *CallbackField[T]) AddCallback(value T, callback Callback) {
	cf.mtx.Lock()
	defer cf.mtx.Unlock()
	cf.callbacks[value] = append(cf.callbacks[value], callback)
}

func (cf *CallbackField[T]) AddOnceCallback(value T, callback Callback) {
	cf.mtx.Lock()
	defer cf.mtx.Unlock()
	cf.onceCallbacks[value] = append(cf.onceCallbacks[value], callback)
}

func (cf *CallbackField[T]) AddAnyCallback(callback Callback) {
	cf.mtx.Lock()
	defer cf.mtx.Unlock()
	cf.anyCallbacks = append(cf.anyCallbacks, callback)
}

func (cf *CallbackField[T]) callback() {
	if cf.anyAwait != nil {
		close(cf.anyAwait)
		cf.anyAwait = nil
	}
	if await, ok := cf.awaits[cf.value]; ok {
		close(await)
		delete(cf.awaits, cf.value)
	}
	for i, orAwait := range cf.orAwaits {
		for _, orValue := range orAwait.values {
			if orValue == cf.value {
				close(orAwait.await)
				cf.orAwaits = append(cf.orAwaits[:i], cf.orAwaits[i+1:]...)
			}
		}
	}
	if callbacks, ok := cf.callbacks[cf.value]; ok {
		for _, callback := range callbacks {
			callback()
		}
	}
	if callbacks, ok := cf.onceCallbacks[cf.value]; ok {
		for _, callback := range callbacks {
			callback()
		}
		delete(cf.onceCallbacks, cf.value)
	}
	for _, callback := range cf.anyCallbacks {
		callback()
	}
}
