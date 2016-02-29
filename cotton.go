// Package cotton is simple, lightweight and composable HTTP Handler/Middleware
package cotton

import (
	"net/http"
)

// Error is a result of Middleware.
type Error interface{}

// RecoverFunc is a callback for a Service.
// Called only when the Service returned a Error.
type RecoverFunc func(w http.ResponseWriter, r *http.Request, err Error)
