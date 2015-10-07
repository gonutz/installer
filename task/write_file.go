package task

import (
	"io"
	"os"
)

func WriteFile(path string, r io.Reader) Task {
	return &writeFile{path, r}
}

type writeFile struct {
	path string
	r    io.Reader
}

func (t *writeFile) Name() string {
	return "Creating " + t.path
}

func (t *writeFile) Execute() error {
	file, err := os.Create(t.path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, t.r)
	return err
}
