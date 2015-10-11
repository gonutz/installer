package task

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// CreateShortcut takes the exePath which is the destination file to execute
// from the link, the linkPath which is the folder in which the link is to be
// created, and the linkName which is the name the user sees for the link. It
// must not contain an extension, the extension is automatically added by the
// task.
func CreateShortcut(exePath, linkPath, linkName string) Task {
	_, exe := filepath.Split(exePath)
	if strings.Contains(exe, " ") {
		panic("for now the file part of a link destination must not contain spaces")
	}
	return &shortcut{exePath, linkPath, linkName}
}

type shortcut struct {
	exePath, linkPath, linkName string
}

func (t *shortcut) Name() string {
	return "Creating shortcut to\n" +
		t.exePath + "\nin\n" +
		t.linkPath +
		"\nwith the name " + t.linkName
}

func (t *shortcut) Execute() error {
	fileName := filepath.Join(t.linkPath, t.linkName+".bat")
	exePath, exe := filepath.Split(t.exePath)
	content := fmt.Sprintf(`start /d "%v" %v`, exePath, exe)
	return makeError(
		"creating link file '"+fileName+"'",
		ioutil.WriteFile(fileName, []byte(content), 0666))
}
