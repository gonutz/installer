package task

import (
	"bytes"
	"io"
	"os"
)

func WriteFile(path string, data []byte) Task {
	return &writeFile{path, data}
}

type writeFile struct {
	path string
	data []byte
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

	_, err = io.Copy(file, bytes.NewReader(t.data))
	return err
}
