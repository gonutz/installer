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
	return makeError(
		"creating folder '"+t.path+"'",
		os.MkdirAll(t.path, 0666))
}
