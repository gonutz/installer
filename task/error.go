package task

import "errors"

func makeError(msg string, err error) error {
	return errors.New(msg + ": " + err.Error())
}
