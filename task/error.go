package task

import "errors"

func makeError(msg string, err error) error {
	if err == nil {
		return nil
	}
	return errors.New(msg + ": " + err.Error())
}
