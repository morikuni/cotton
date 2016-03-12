package yacm

import (
	"net/http"
	"testing"

	"github.com/morikuni/yacm/testutil"
)

func TestService(t *testing.T) {
	tt := &testutil.T{t}

	callMe := tt.CallMe()
	service := Service(func(w http.ResponseWriter, r *http.Request) Error {
		callMe.Call()
		return nil
	})

	handler := service.Recover(func(w http.ResponseWriter, r *http.Request, err Error) {
		tt.Error("unreachable")
	})

	handler(nil, nil)
	callMe.MustCalled()
}
