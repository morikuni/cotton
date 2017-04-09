package yacm

import (
	"net/http"
)

type Catcher interface {
	HandleError(w http.ResponseWriter, r *http.Request, err error) error
}

type CatcherFunc func(w http.ResponseWriter, r *http.Request, err error) error

func (f CatcherFunc) HandleError(w http.ResponseWriter, r *http.Request, err error) error {
	return f(w, r, err)
}

type chainedCatcher []Catcher

func (c chainedCatcher) HandleError(w http.ResponseWriter, r *http.Request, err error) error {
	for _, h := range c {
		err = h.HandleError(w, r, err)
		if err == nil {
			return nil
		}
	}
	return err
}

func ComposeCatchers(cs ...Catcher) Catcher {
	l := len(cs)
	switch l {
	case 0:
		panic(ErrEmptyArgs)
	default:
		return chainedCatcher(cs)
	}
}

func ApplyCatcher(c Catcher, s Shutter) Shutter {
	return ShutterFunc(func(w http.ResponseWriter, r *http.Request, err error) {
		err = c.HandleError(w, r, err)
		if err != nil {
			s.ShutError(w, r, err)
		}
	})
}

func ApplyCatcherToFunc(c Catcher, s func(http.ResponseWriter, *http.Request, error)) Shutter {
	return ApplyCatcher(c, ShutterFunc(s))
}
