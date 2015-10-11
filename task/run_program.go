package task

import (
	"os/exec"
	"strings"
)

func RunProgram(path string, params ...string) Task {
	return &runProgram{path, params}
}

type runProgram struct {
	path   string
	params []string
}

func (t *runProgram) Name() string {
	return "Running " + t.path
}

func (t *runProgram) Execute() error {
	params := []string{"/C", t.path}
	for _, param := range t.params {
		params = append(params, param)
	}
	err := exec.Command("cmd", params...).Run()
	return makeError("running program '"+strings.Join(params, " ")+"'", err)
}
