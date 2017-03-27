// Package yacm is simple, lightweight and composable HTTP Handler/Middleware
package yacm

import (
	"errors"
)

var (
	ErrEmptyArgs = errors.New("arguments are required")
)
