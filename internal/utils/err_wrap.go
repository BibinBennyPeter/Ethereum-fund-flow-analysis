package utils

import (
	"fmt"
)

// wrapErr gives error some context msg
// returns nil if err is nil
func wrapErr(err error, msg string) (errWithContext error) {
	if err == nil {
		return
	}

	errWithContext = fmt.Errorf("%s: %v", msg, err)
	return
}

// wrapfErr gives error some context msg
// with desired format and content
// returns nil if err is nil
func wrapfErr(err error, format string, a ...interface{}) (errWithContext error) {
	if err == nil {
		return
	}

	errWithContext = wrapErr(err, fmt.Sprintf(format, a...))
	return
}
