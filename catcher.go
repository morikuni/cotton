package yacm

import (
	"net/http"
)

// Catcher handles an error from Service or Filter.
type Catcher interface {
	// CatchError handles an error, but it may returns an error if it was unexpected.
	CatchError(w http.ResponseWriter, r *http.Request, err error) error
}

// CatcherFunc is the adapter to use a function as Catcher.
type CatcherFunc func(w http.ResponseWriter, r *http.Request, err error) error

// CatchError implements Catcher.
func (f CatcherFunc) CatchError(w http.ResponseWriter, r *http.Request, err error) error {
	return f(w, r, err)
}

type chainedCatcher []Catcher

func (c chainedCatcher) CatchError(w http.ResponseWriter, r *http.Request, err error) error {
	for _, h := range c {
		err = h.CatchError(w, r, err)
		if err == nil {
			return nil
		}
	}
	return err
}

// ComposeCatchers composes multiple Catchers to a single Catcher.
// NOPCatcher will be returned for empty arguments.
func ComposeCatchers(cs ...Catcher) Catcher {
	l := len(cs)
	switch l {
	case 0:
		return NOPCatcher
	default:
		return chainedCatcher(cs)
	}
}

// ApplyCatcher apply Catcher to Shutter and creates a new Shutter.
func ApplyCatcher(c Catcher, s Shutter) Shutter {
	return ShutterFunc(func(w http.ResponseWriter, r *http.Request, err error) {
		err = c.CatchError(w, r, err)
		if err != nil {
			s.ShutError(w, r, err)
		}
	})
}

var (
	// NOPCatcher is Catcher that handles nothing.
	NOPCatcher = CatcherFunc(nopCatcher)
)

func nopCatcher(w http.ResponseWriter, r *http.Request, err error) error {
	return err
}
