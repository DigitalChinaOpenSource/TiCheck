package executor

type Result struct {
	executionError error // script level error
	result         CheckData
}

type Executor interface {
	ExecuteCheck(chan Result) // executes one round of check
}
