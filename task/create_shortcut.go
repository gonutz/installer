package task

import "errors"

func CreateShortcut(exePath, linkPath string) Task {
	return &shortcut{exePath, linkPath}
}

type shortcut struct {
	exePath, linkPath string
}

func (t *shortcut) Name() string {
	return "Creating shortcut to\n" + t.exePath + "\nin\n" + t.linkPath
}

func (t *shortcut) Execute() error {
	return errors.New("TODO create shortcut")
}
