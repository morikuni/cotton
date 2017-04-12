package yacm

import (
	"net/http"
)

var (
	EmptyBuilder = newBuilder()
)

type Builder struct {
	filter  Filter
	catcher Catcher
	shutter Shutter
}

func newBuilder() Builder {
	return Builder{nil, nil, DefaultShutter}
}

func (b Builder) cloneWithFilter(f Filter) Builder {
	return Builder{f, b.catcher, b.shutter}
}

func (b Builder) AppendFilters(fs ...Filter) Builder {
	f := ComposeFilters(fs...)
	if b.filter != nil {
		f = ComposeFilters(b.filter, f)
	}
	return b.cloneWithFilter(f)
}

func (b Builder) AppendFilterFunc(f func(http.ResponseWriter, *http.Request, Service) error) Builder {
	return b.AppendFilters(FilterFunc(f))
}

func (b Builder) AppendMiddlewares(ms ...Middleware) Builder {
	filters := make([]Filter, len(ms))
	for i, m := range ms {
		filters[i] = m
	}
	return b.AppendFilters(filters...)
}

func (b Builder) cloneWithCatcher(c Catcher) Builder {
	return Builder{b.filter, c, b.shutter}
}

func (b Builder) AppendCatchers(cs ...Catcher) Builder {
	c := ComposeCatchers(cs...)
	if b.catcher != nil {
		c = ComposeCatchers(b.catcher, c)
	}
	return b.cloneWithCatcher(c)
}

func (b Builder) AppendCatcherFunc(f func(http.ResponseWriter, *http.Request, error) error) Builder {
	return b.AppendCatchers(CatcherFunc(f))
}

func (b Builder) cloneWithShutter(s Shutter) Builder {
	return Builder{b.filter, b.catcher, s}
}

func (b Builder) WithShutter(s Shutter) Builder {
	return b.cloneWithShutter(s)
}

func (b Builder) WithShutterFunc(f func(http.ResponseWriter, *http.Request, error)) Builder {
	return b.WithShutter(ShutterFunc(f))
}

func (b Builder) Apply(s Service) http.Handler {
	shutter := b.shutter
	if b.catcher != nil {
		shutter = ApplyCatcher(b.catcher, shutter)
	}
	if b.filter != nil {
		s = ApplyFilter(b.filter, s)
	}
	return ServiceToHandler(s, shutter)
}

func (b Builder) ApplyFunc(f func(http.ResponseWriter, *http.Request) error) http.Handler {
	return b.Apply(ServiceFunc(f))
}
