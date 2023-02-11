package callbackfield

type CallbackFieldOption[T comparable] func(cf *CallbackField[T])

func WithValue[T comparable](value T) CallbackFieldOption[T] {
	return func(cf *CallbackField[T]) {
		cf.value = value
	}
}
