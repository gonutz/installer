package task

import "os/exec"

func SetEnv(varName, varValue string) Task {
	return &setEnv{varName, varValue}
}

type setEnv struct {
	name, value string
}

func (t *setEnv) Name() string {
	return "Setting " + t.name + " environment variable"
}

func (t *setEnv) Execute() error {
	return makeError(
		"running 'setx "+t.name+" "+t.value+"'",
		exec.Command("setx", t.name, t.value).Run())
}
