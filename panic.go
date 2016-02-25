package cotton

import (
	"net/http"
)

type PanicOccured struct {
	Reason interface{}
}

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
