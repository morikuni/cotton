package yacm

import (
	"net/http"
	"testing"

	"github.com/morikuni/yacm/testutil"
)

func TestMethodFilter(t *testing.T) {
	tt := testutil.T{t}
	filter := MethodFilter(GET, POST)
	req, _ := http.NewRequest("", "", nil)

	service := filter.For(testutil.NOPHandler)
	h1 := service.Recover(func(w http.ResponseWriter, r *http.Request, err Error) {
		tt.Error("unreachable")
	})
	req.Method = GET
	h1(nil, req)
	req.Method = POST
	h1(nil, req)

	callMe := tt.CallMe()
	h2 := service.Recover(func(w http.ResponseWriter, r *http.Request, err Error) {
		e, _ := err.(MethodNotAllowed)
		tt.MustEqual(e.Method, DELETE)
		tt.MustEqual(e.Expect[0], GET)
		tt.MustEqual(e.Expect[1], POST)
		callMe.Call()
	})
	req.Method = DELETE
	h2(nil, req)
	callMe.MustCalled()
}
