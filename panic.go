package yacm

import (
	"net/http"
)

// PanicOccured is a Error of PanicFilter.
type PanicOccured struct {
	Reason interface{}
}

// PanicFilter recovers a panic.
func PanicFilter(w http.ResponseWriter, r *http.Request, s Service) (err Error) {
	defer func() {
		e := recover()
		if e != nil {
			err = PanicOccured{
				Reason: e,
			}
		}
	}()
	err = s(w, r)
	return
}
