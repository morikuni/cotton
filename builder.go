package yacm

import (
	"net/http"
)

// Builder builds a http.Handler from Filters, Catchers, Shutters, http.Handlers and Services.
type Builder struct {
	Filter  Filter
	Catcher Catcher
	Shutter Shutter
}

// NewBuilder creates a new Builder.
func NewBuilder() Builder {
	return Builder{}
}

func (b Builder) cloneWithFilter(f Filter) Builder {
	return Builder{f, b.Catcher, b.Shutter}
}

// AppendFilters appends Filters to the tail of the current Filter.
func (b Builder) AppendFilters(fs ...Filter) Builder {
	f := ComposeFilters(fs...)
	if b.Filter != nil {
		f = ComposeFilters(b.Filter, f)
	}
	return b.cloneWithFilter(f)
}

// AppendFilterFunc is same as AppendFilters but it takes FilterFunc.
func (b Builder) AppendFilterFunc(f func(http.ResponseWriter, *http.Request, Service) error) Builder {
	return b.AppendFilters(FilterFunc(f))
}

// AppendMiddlewares converts Middlewares to Filters, then appends to the tail of the current Filter.
func (b Builder) AppendMiddlewares(ms ...Middleware) Builder {
	filters := make([]Filter, len(ms))
	for i, m := range ms {
		filters[i] = m
	}
	return b.AppendFilters(filters...)
}

func (b Builder) cloneWithCatcher(c Catcher) Builder {
	return Builder{b.Filter, c, b.Shutter}
}

// AppendCatchers appends Catchers to the tail of the current Catcher.
func (b Builder) AppendCatchers(cs ...Catcher) Builder {
	c := ComposeCatchers(cs...)
	if b.Catcher != nil {
		c = ComposeCatchers(b.Catcher, c)
	}
	return b.cloneWithCatcher(c)
}

// AppendCatcherFunc is same as AppendCatchers but it takes CatcherFunc.
func (b Builder) AppendCatcherFunc(f func(http.ResponseWriter, *http.Request, error) error) Builder {
	return b.AppendCatchers(CatcherFunc(f))
}

func (b Builder) cloneWithShutter(s Shutter) Builder {
	return Builder{b.Filter, b.Catcher, s}
}

// WithShutter replace the Shutter by the argument.
func (b Builder) WithShutter(s Shutter) Builder {
	return b.cloneWithShutter(s)
}

// WithShutterFunc is save as WithShutterFunc but it takes WithShutterFunc.
func (b Builder) WithShutterFunc(f func(http.ResponseWriter, *http.Request, error)) Builder {
	return b.WithShutter(ShutterFunc(f))
}

// Apply applys the Builder to Service and creates a http.Handler.
func (b Builder) Apply(s Service) http.Handler {
	shutter := b.Shutter
	if shutter == nil {
		shutter = DefaultShutter
	}
	if b.Catcher != nil {
		shutter = ApplyCatcher(b.Catcher, shutter)
	}
	if b.Filter != nil {
		s = ApplyFilter(b.Filter, s)
	}
	return ServiceToHandler(s, shutter)
}

// ApplyFunc is same as Apply but it takes ServiceFunc.
func (b Builder) ApplyFunc(f func(http.ResponseWriter, *http.Request) error) http.Handler {
	return b.Apply(ServiceFunc(f))
}

// ApplyHandler converts http.Handler to Service, then apply the Builder.
func (b Builder) ApplyHandler(h http.Handler) http.Handler {
	return b.Apply(HandlerToService(h))
}

// ApplyHandlerFunc is same as ApplyHandler but it takes http.HandlerFunc.
func (b Builder) ApplyHandlerFunc(f func(http.ResponseWriter, *http.Request)) http.Handler {
	return b.ApplyHandler(http.HandlerFunc(f))
}
