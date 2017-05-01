package yacm

import (
	"net/http"
)

// Env is a environment for building http.Handler and Service.
// Env contains Filter, Catcher and Shutter and it is immutable,
// so any methods never changes its state and return a new Env.
type Env struct {
	Filter  Filter
	Catcher Catcher
	Shutter Shutter
}

// NewEnv creates a new Env.
func NewEnv() Env {
	return Env{}
}

func (e Env) cloneWithFilter(f Filter) Env {
	return Env{f, e.Catcher, e.Shutter}
}

// AppendFilters appends Filters to the tail of the current Filter.
func (e Env) AppendFilters(fs ...Filter) Env {
	f := ComposeFilters(fs...)
	if e.Filter != nil {
		f = ComposeFilters(e.Filter, f)
	}
	return e.cloneWithFilter(f)
}

// AppendFilterFunc is same as AppendFilters but it takes FilterFunc.
func (e Env) AppendFilterFunc(f func(http.ResponseWriter, *http.Request, Service) error) Env {
	return e.AppendFilters(FilterFunc(f))
}

// AppendMiddlewares converts Middlewares to Filters, then appends to the tail of the current Filter.
func (e Env) AppendMiddlewares(ms ...Middleware) Env {
	filters := make([]Filter, len(ms))
	for i, m := range ms {
		filters[i] = m
	}
	return e.AppendFilters(filters...)
}

func (e Env) cloneWithCatcher(c Catcher) Env {
	return Env{e.Filter, c, e.Shutter}
}

// AppendCatchers appends Catchers to the tail of the current Catcher.
func (e Env) AppendCatchers(cs ...Catcher) Env {
	c := ComposeCatchers(cs...)
	if e.Catcher != nil {
		c = ComposeCatchers(e.Catcher, c)
	}
	return e.cloneWithCatcher(c)
}

// AppendCatcherFunc is same as AppendCatchers but it takes CatcherFunc.
func (e Env) AppendCatcherFunc(f func(http.ResponseWriter, *http.Request, error) error) Env {
	return e.AppendCatchers(CatcherFunc(f))
}

func (e Env) cloneWithShutter(s Shutter) Env {
	return Env{e.Filter, e.Catcher, s}
}

// WithShutter replace the Shutter by the argument.
func (e Env) WithShutter(s Shutter) Env {
	return e.cloneWithShutter(s)
}

// WithShutterFunc is same as WithShutterFunc but it takes WithShutterFunc.
func (e Env) WithShutterFunc(f func(http.ResponseWriter, *http.Request, error)) Env {
	return e.WithShutter(ShutterFunc(f))
}

// Serve applys the Env to Service and creates a http.Handler.
func (e Env) Serve(s Service) http.Handler {
	shutter := e.Shutter
	if shutter == nil {
		shutter = DefaultShutter
	}
	if e.Catcher != nil {
		shutter = ApplyCatcher(e.Catcher, shutter)
	}
	if e.Filter != nil {
		s = ApplyFilter(e.Filter, s)
	}
	return ServiceToHandler(s, shutter)
}

// ServeFunc is same as Serve but it takes ServiceFunc.
func (e Env) ServeFunc(f func(http.ResponseWriter, *http.Request) error) http.Handler {
	return e.Serve(ServiceFunc(f))
}

// Handle converts http.Handler to Service, then call the Serve with it.
func (e Env) Handle(h http.Handler) http.Handler {
	return e.Serve(HandlerToService(h))
}

// HandleFunc is same as Handle but it takes http.HandlerFunc.
func (e Env) HandleFunc(f func(http.ResponseWriter, *http.Request)) http.Handler {
	return e.Handle(http.HandlerFunc(f))
}
