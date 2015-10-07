package task

import "errors"

// newCheckTask creates a task that fails with the given error message if the
// condition is false.
func Check(condition func() bool, errorMsg string) Task {
	return &check{condition, errors.New(errorMsg)}
}

type check struct {
	condition func() bool
	err       error
}

func (t *check) Name() string { return "" }

func (t *check) Execute() error {
	if !t.condition() {
		return t.err
	}
	return nil
}
