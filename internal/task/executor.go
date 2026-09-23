package task

type Executor interface {
	Execute(payload string) error
}
