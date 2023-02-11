package tools

type TaskResult[R any] struct {
	Result R
	Err    error
}
