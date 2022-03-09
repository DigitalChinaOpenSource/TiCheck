package executor

import "TiCheck/internal/model"

type Result struct {
	executionError error // script level error
	result         model.CheckData
}

type Executor interface {
	ExecuteCheck(chan Result) // executes one round of check
}
