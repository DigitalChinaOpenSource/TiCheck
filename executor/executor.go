package executor

type Result struct {
	executionError error // script level error
	result         CheckData
}

type Executor interface {
	ExecuteCheck(string, chan Result) // executes one round of check
}
