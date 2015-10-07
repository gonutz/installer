package task

import (
	"bytes"
	"io"
	"io/ioutil"
	"os/exec"
)

func AddToPathEnv(add string) Task {
	return &addToPath{add}
}

type addToPath struct {
	add string
}

func (t *addToPath) Name() string {
	return "Adding " + t.add + " to PATH environment variable"
}

func (t *addToPath) Execute() error {
	if len(pathedPath) == 0 {
		if err := initPathed(); err != nil {
			return err
		}
	}

	return exec.Command("cmd", "/C", pathedPath, "-f", "-a", t.add).Run()

	//path := os.Getenv("PATH")

	//// check if the addition is already contained in the PATH
	//if strings.Contains(path, ";"+t.add+";") || // in the middle?
	//	strings.HasPrefix(path, t.add+";") || // at the start?
	//	strings.HasSuffix(path, ";"+t.add) || // at the end?
	//	path == t.add { // PATH only contains addition
	//	return nil // if PATH contains it, we are done
	//}

	//newPath := path
	//if newPath != "" {
	//	newPath += ";"
	//}
	//newPath += t.add

	//return exec.Command("setx", "PATH", newPath).Run()
}

var pathedPath string

func initPathed() error {
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}

	file, err := ioutil.TempFile(tempDir, "pathed.exe")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, bytes.NewReader(pathedExe))
	if err != nil {
		return err
	}

	pathedPath = file.Name()
	return nil
}
