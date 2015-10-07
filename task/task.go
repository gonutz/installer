package task

type Task interface {
	Name() string
	Execute() error
}

func New(name string, execute func() error) Task {
	return &simpleTask{name, execute}
}

type simpleTask struct {
	name    string
	execute func() error
}

func (t *simpleTask) Name() string   { return t.name }
func (t *simpleTask) Execute() error { return t.execute() }
