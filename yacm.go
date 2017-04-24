// yacm provides a way to build a http handlers with middlewares.
package yacm

import (
	"errors"
)

var (
	// ErrEmptyArgs is used when the arguments to compose function are empty.
	ErrEmptyArgs = errors.New("arguments are required")
)
