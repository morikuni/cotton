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

func (f Tester) HandleError(w http.ResponseWriter, r *http.Request, err error) error {
	f.a.Equal(f.expect, *f.count)
	*f.count++
	return err
}

func TestServiceBuilder(t *testing.T) {
	assert := assert.New(t)

	b := EmptyServiceBuilder

	count := 0
	handler := b.AppendFilterFunc(func(w http.ResponseWriter, r *http.Request, next Service) error {
		assert.Equal(0, count)
		count++
		return next.TryServeHTTP(w, r)
	}).AppendFilters(
		Tester{assert, &count, 1},
		Tester{assert, &count, 2},
	).AppendCatcherFunc(func(w http.ResponseWriter, r *http.Request, err error) error {
		assert.Equal(4, count)
		count++
		return err
	}).AppendCatchers(
		Tester{assert, &count, 5},
		Tester{assert, &count, 6},
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
	assert.Equal(nil, EmptyServiceBuilder.filter)
	assert.Equal(nil, EmptyServiceBuilder.catcher)
}
