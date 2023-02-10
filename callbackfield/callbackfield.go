package callbackfield

type (
	Callback func()

	OrChannel[T comparable] struct {
		values []T
		ch     chan struct{}
	}

	CallbackField[T comparable] struct {
		value T

		channel    chan struct{}
		channels   map[T]chan struct{}
		orChannels []OrChannel[T]

		callbacks     map[T][]Callback
		onceCallbacks map[T][]Callback
		anyCallbacks  []Callback

		mtx RWMutex
	}
)

func NewCallbackField[T comparable](opts ...CallbackFieldOption[T]) *CallbackField[T] {
	field := &CallbackField[T]{
		channels: make(map[T]chan struct{}),

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
		if cf.channel == nil {
			cf.channel = make(chan struct{})
		}
		return cf.channel
	} else if len(values) == 1 {
		value := values[0]
		if value == cf.value {
			cf.mtx.RUnlock()
			return ClosedChannel
		}
		cf.mtx.RUnlockToLock()
		defer cf.mtx.Unlock()
		if _, ok := cf.channels[value]; !ok {
			cf.channels[value] = make(chan struct{})
		}
		return cf.channels[value]
	} else {
		for _, value := range values {
			if value == cf.value {
				cf.mtx.RUnlock()
				return ClosedChannel
			}
		}
		cf.mtx.RUnlockToLock()
		defer cf.mtx.Unlock()
		orChannel := OrChannel[T]{values, make(chan struct{})}
		cf.orChannels = append(cf.orChannels, orChannel)
		return orChannel.ch
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
	if cf.channel != nil {
		close(cf.channel)
		cf.channel = nil
	}
	if ch, ok := cf.channels[cf.value]; ok {
		close(ch)
		delete(cf.channels, cf.value)
	}
	for i, orChannel := range cf.orChannels {
		for _, orValue := range orChannel.values {
			if orValue == cf.value {
				close(orChannel.ch)
				cf.orChannels = append(cf.orChannels[:i], cf.orChannels[i+1:]...)
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
