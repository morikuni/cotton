package yacm

import (
	"net/http"
)

var (
	EmptyServiceBuilder = newServiceBuilder()
)

type ServiceBuilder struct {
	filter  Filter
	handler ErrorHandler
	shutter ErrorShutter
}

func newServiceBuilder() ServiceBuilder {
	return ServiceBuilder{nil, nil, DefaultErrorShutter}
}

func (b ServiceBuilder) cloneWithFilter(f Filter) ServiceBuilder {
	return ServiceBuilder{f, b.handler, b.shutter}
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

func (b ServiceBuilder) cloneWithErrorHandler(eh ErrorHandler) ServiceBuilder {
	return ServiceBuilder{b.filter, eh, b.shutter}
}

func (b ServiceBuilder) AppendErrorHandlers(ehs ...ErrorHandler) ServiceBuilder {
	eh := ComposeErrorHandlers(ehs...)
	if b.handler != nil {
		eh = ComposeErrorHandlers(b.handler, eh)
	}
	return b.cloneWithErrorHandler(eh)
}

func (b ServiceBuilder) AppendErrorHandlerFunc(f func(http.ResponseWriter, *http.Request, error) error) ServiceBuilder {
	return b.AppendErrorHandlers(ErrorHandlerFunc(f))
}

func (b ServiceBuilder) cloneWithErrorShutter(es ErrorShutter) ServiceBuilder {
	return ServiceBuilder{b.filter, b.handler, es}
}

func (b ServiceBuilder) WithErrorShutter(es ErrorShutter) ServiceBuilder {
	return b.cloneWithErrorShutter(es)
}

func (b ServiceBuilder) WithErrorShutterFunc(f func(http.ResponseWriter, *http.Request, error)) ServiceBuilder {
	return b.WithErrorShutter(ErrorShutterFunc(f))
}

func (b ServiceBuilder) Apply(s Service) http.Handler {
	shutter := b.shutter
	if b.handler != nil {
		shutter = ApplyErrorHandler(b.handler, shutter)
	}
	if b.filter != nil {
		s = ApplyFilter(b.filter, s)
	}
	return ServiceToHandler(s, shutter)
}

func (b ServiceBuilder) ApplyFunc(f func(http.ResponseWriter, *http.Request) error) http.Handler {
	return b.Apply(ServiceFunc(f))
}
