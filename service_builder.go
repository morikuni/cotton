package yacm

import (
	"net/http"
)

var (
	EmptyServiceBuilder = newServiceBuilder()
)

type ServiceBuilder struct {
	filter  Filter
	catcher Catcher
	shutter ErrorShutter
}

func newServiceBuilder() ServiceBuilder {
	return ServiceBuilder{nil, nil, DefaultErrorShutter}
}

func (b ServiceBuilder) cloneWithFilter(f Filter) ServiceBuilder {
	return ServiceBuilder{f, b.catcher, b.shutter}
}

func (b ServiceBuilder) AppendFilters(fs ...Filter) ServiceBuilder {
	f := ComposeFilters(fs...)
	if b.filter != nil {
		f = ComposeFilters(b.filter, f)
	}
	return b.cloneWithFilter(f)
}

func (b ServiceBuilder) AppendFilterFunc(f func(http.ResponseWriter, *http.Request, Service) error) ServiceBuilder {
	return b.AppendFilters(FilterFunc(f))
}

func (b ServiceBuilder) cloneWithCatcher(c Catcher) ServiceBuilder {
	return ServiceBuilder{b.filter, c, b.shutter}
}

func (b ServiceBuilder) AppendCatchers(cs ...Catcher) ServiceBuilder {
	c := ComposeCatchers(cs...)
	if b.catcher != nil {
		c = ComposeCatchers(b.catcher, c)
	}
	return b.cloneWithCatcher(c)
}

func (b ServiceBuilder) AppendCatcherFunc(f func(http.ResponseWriter, *http.Request, error) error) ServiceBuilder {
	return b.AppendCatchers(CatcherFunc(f))
}

func (b ServiceBuilder) cloneWithErrorShutter(es ErrorShutter) ServiceBuilder {
	return ServiceBuilder{b.filter, b.catcher, es}
}

func (b ServiceBuilder) WithErrorShutter(es ErrorShutter) ServiceBuilder {
	return b.cloneWithErrorShutter(es)
}

func (b ServiceBuilder) WithErrorShutterFunc(f func(http.ResponseWriter, *http.Request, error)) ServiceBuilder {
	return b.WithErrorShutter(ErrorShutterFunc(f))
}

func (b ServiceBuilder) Apply(s Service) http.Handler {
	shutter := b.shutter
	if b.catcher != nil {
		shutter = ApplyCatcher(b.catcher, shutter)
	}
	if b.filter != nil {
		s = ApplyFilter(b.filter, s)
	}
	return ServiceToHandler(s, shutter)
}

func (b ServiceBuilder) ApplyFunc(f func(http.ResponseWriter, *http.Request) error) http.Handler {
	return b.Apply(ServiceFunc(f))
}
