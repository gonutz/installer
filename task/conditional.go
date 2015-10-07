package task

func Conditional(condition func() bool, ifTask, elseTask Task) Task {
	return &conditional{condition, ifTask, elseTask}
}

type conditional struct {
	condition        func() bool
	ifTask, elseTask Task
}

func (t *conditional) Name() string {
	if t.condition() {
		return t.ifTask.Name()
	}
	return t.elseTask.Name()
}

func (t *conditional) Execute() error {
	if t.condition() {
		return t.ifTask.Execute()
	}
	return t.elseTask.Execute()
}
