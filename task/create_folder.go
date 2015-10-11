package task

import "os"

func CreateFolder(path string) Task {
	return &createFolder{path}
}

type createFolder struct{ path string }

func (t *createFolder) Name() string {
	return "Creating folder " + t.path
}

func (t *createFolder) Execute() error {
	if err := os.MkdirAll(t.path, 0666); err != nil {
		return makeError("creating folder '"+t.path+"'", err)
	}
	return nil
}
