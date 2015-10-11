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
	err := exec.Command("setx", t.name, t.value).Run()
	if err == nil {
		return nil
	}

	// SETX did not work so it might not exists, try another way
	return makeError(
		"setting environment variable "+t.name,
		exec.Command("cmd", "/C", "reg", "add", `HKCU\Environment`,
			"/v", t.name,
			"/t", "REG_SZ",
			"/d", t.value).Run(),
	)
}
