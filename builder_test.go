package yacm

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Tester struct {
	a      *assert.Assertions
	count  *int
	expect int
}

func (f Tester) WrapService(w http.ResponseWriter, r *http.Request, next Service) error {
	f.a.Equal(f.expect, *f.count)
	*f.count++
	return next.TryServeHTTP(w, r)
}

func (f Tester) CatchError(w http.ResponseWriter, r *http.Request, err error) error {
	f.a.Equal(f.expect, *f.count)
	*f.count++
	return err
}

func TestBuilder(t *testing.T) {
	assert := assert.New(t)

	b := EmptyBuilder

	count := 0
	handler := b.AppendMiddlewares(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(0, count)
			count++
			next.ServeHTTP(w, r)
		})
	}).AppendFilters(
		Tester{assert, &count, 1},
		Tester{assert, &count, 2},
	).AppendCatcherFunc(func(w http.ResponseWriter, r *http.Request, err error) error {
		assert.Equal(6, count)
		count++
		return err
	}).AppendCatchers(
		Tester{assert, &count, 4},
		Tester{assert, &count, 5},
	).WithShutterFunc(func(w http.ResponseWriter, r *http.Request, err error) {
		assert.Equal(7, count)
		w.WriteHeader(http.StatusTeapot)
		w.Write([]byte(err.Error()))
	}).ApplyFunc(func(w http.ResponseWriter, r *http.Request) error {
		assert.Equal(3, count)
		count++
		return errors.New("test")
	})

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, nil)
	assert.Equal("test", rec.Body.String())
	assert.Equal(http.StatusTeapot, rec.Code)
	assert.Equal(nil, EmptyBuilder.filter)
	assert.Equal(nil, EmptyBuilder.catcher)
}
