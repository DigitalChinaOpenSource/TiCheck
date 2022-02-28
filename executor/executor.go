package executor

type Result struct {
	executionError error
	resultID       int
}

type Executor interface {
	Execute(string, chan Result)
}
