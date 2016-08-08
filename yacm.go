// Package yacm is simple, lightweight and composable HTTP Handler/Middleware
package yacm

import (
	"net/http"
)

// ErrorHandler is a callback for a Service.
// Called only when the Service returned a error.
type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)
