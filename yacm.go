// yacm provides a way to compose handlers and middlewares for net/http, inspired by Twitter's Finagle.
package yacm

import (
	"errors"
)

var (
	// ErrEmptyArgs is used when the arguments to compose function are empty.
	ErrEmptyArgs = errors.New("arguments are required")
)
