package task

import (
	"fmt"
	"strings"
)

func Inform(msg string) Task {
	return information(msg)
}

type information string

func (information) Name() string { return "" }

func (info information) Execute() error {
	fmt.Println(line)
	fmt.Println(info)
	fmt.Println(line)
	fmt.Println("Press ENTER to continue...")
	fmt.Scanln()
	return nil
}

var line = strings.Repeat("-", 79)
