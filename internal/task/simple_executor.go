package task

import "errors"

type SimpleExecutor struct{}

func NewSimpleExecutor() *SimpleExecutor {
	return &SimpleExecutor{}
}

func (e *SimpleExecutor) Execute(payload string) error {
	if payload == "" {
		return errors.New("task payload is empty")
	}

	return nil
}
