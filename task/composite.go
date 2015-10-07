package task

import "fmt"

func ContinueAfterError(name string, tasks []Task) Task {
	return newComposite(name, true, tasks)
}

func FailOnFirstError(name string, tasks []Task) Task {
	return newComposite(name, false, tasks)
}

func newComposite(name string, continueAfterError bool, tasks []Task) *composite {
	return &composite{name, continueAfterError, tasks}
}

type composite struct {
	name               string
	continueAfterError bool
	tasks              []Task
}

func (t *composite) Name() string { return t.name }

func (t *composite) Execute() error {
	for _, task := range t.tasks {
		name := task.Name()
		if len(name) > 0 {
			fmt.Println(task.Name())
		}
		if err := task.Execute(); err != nil {
			if t.continueAfterError {
				fmt.Println("Error:", err)
			} else {
				return err
			}
		}
	}
	return nil
}
