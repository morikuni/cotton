package yacm

import (
	"net/http"
	"testing"

	"github.com/morikuni/yacm/testutil"
)

func panicHandler(w http.ResponseWriter, r *http.Request) {
	panic("panic")
}

func TestPanicFilter(t *testing.T) {
	tt := testutil.T{t}
	filter := Middleware(PanicFilter)

	ps := filter.For(panicHandler)
	callMe := tt.CallMe()
	ph := ps.Recover(func(w http.ResponseWriter, r *http.Request, err Error) {
		tt.MustEqual(err, PanicOccured{"panic"})
		callMe.Call()
	})
	ph(nil, nil)
	callMe.MustCalled()

	es := filter.For(testutil.NOPHandler)
	eh := es.Recover(func(w http.ResponseWriter, r *http.Request, err Error) {
		tt.Error("unreachable")
	})
	eh(nil, nil)
}
