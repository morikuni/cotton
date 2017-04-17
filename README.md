# yacm(八雲) - yet another composable middleware

[![Build Status](https://travis-ci.org/morikuni/yacm.svg?branch=master)](https://travis-ci.org/morikuni/yacm)
[![GoDoc](https://godoc.org/github.com/morikuni/yacm?status.svg)](https://godoc.org/github.com/morikuni/yacm)

yacm provides a way to compose handlers and middlewares for `net/http`, inspired by Twitter's Finagle.

## Install

```sh
go get github.com/morikuni/yacm
```

## Design

A typical middleware stack is like this.

```go
type Middleware func(http.Handler) http.Handler

// from net/http
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```

Composition rule is

```
Middleware + Handler = Handler
```

and some libraries provide

```
Middleware + Middleware = Middleware
```

```
 +----------------+
 |     Client     |
 +---+--------^---+
     |        |
  req|     res|
     |        |
 +---v--------+---+ ----------+
 |   Middleware   |           |
 +---+--------^---+           |
     |        |               |
  req|     res|               |
     |        |               |
 +---v--------+---+ -+        |
 |   Middleware   |  |        |Handler
 +---+--------^---+  |        |
     |        |      |        |
  req|     res|      |Handler |
     |        |      |        |
 +---v--------+---+  |        |
 |    Handler     |  |        |
 +----------------+ -+ -------+
```

yacm's stack is this.

```go
type Filter interface {
	WrapService(w http.ResponseWriter, r *http.Request, s Service) error
}

type Service interface {
	TryServeHTTP(w http.ResponseWriter, r *http.Request) error
}

type Catcher interface {
	CatchError(w http.ResponseWriter, r *http.Request, err error) error
}

type Shutter interface {
	ShutError(w http.ResponseWriter, r *http.Request, err error)
}

// from net/http
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```

`Service` corresponds to `Handler`, and `Filter` corresponds to `Middleware`.

Composition rule is

```
Filter  + Filter  = Filter
Filter  + Service = Service
Cather  + Cather  = Cather
Cather  + Shutter = Shutter
Service + Shutter = Handler
```

```
                   +----------------+
                   |     Client     |
                   +---+--------^---+
                       |        |
                    req|     res+----------+
                       |        |          |
       +------- +- +---v--------+---+      |  +---------------+ ----------+
       |        |  |     Filter     +--+   +--+    Shutter    |           |
       |        |  +---+--------^---+  |   |  +------------^--+           |
       |        |      |        |      |   |               |              |
       |  Filter|   req|     res|      |   |res         err|              |
       |        |      |        |      |   |               |              |
       |        |  +---v--------+---+  |   |  +------------+--+ -+        |
Service|        |  |     Filter     +--+   +--+    Catcher    |  |        |Shutter
       |        +- +---+--------^---+  |   |  +------------^--+  |        |
       |               |        |      |   |               |     |        |
       |            req|     res|      |   +------+     err|     |Catcher |
       |               |        |      |          |        |     |        |
       |           +---v--------+---+  |      +---+--------+--+  |        |
       |           |    Service     +--+------>    Catcher    |  |        |
       +---------- +----------------+   err   +---------------+ -+ -------+
```

You can also use `Handler` as a `Service`, `Middleware` as a `Filter`.

## Example

```go
package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/morikuni/yacm"
)

func Logging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL)
		h.ServeHTTP(w, r)
	})
}

var ErrNotGet = errors.New("method is not GET")

func GetOnly(w http.ResponseWriter, r *http.Request, next yacm.Service) error {
	if r.Method != "GET" {
		return ErrNotGet
	}
	return next.TryServeHTTP(w, r)
}

func Catcher(w http.ResponseWriter, r *http.Request, err error) error {
	switch err {
	case ErrNotGet:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return nil
	default:
		return err
	}
}

func main() {
	b := yacm.EmptyBuilder.
		AppendMiddlewares(Logging).
		AppendFilterFunc(GetOnly).
		AppendCatcherFunc(Catcher)

	http.Handle("/hello", b.ApplyFunc(func(w http.ResponseWriter, r *http.Request) error {
		w.Write([]byte("hello"))
		return nil
	}))
	http.Handle("/error", b.ApplyFunc(func(w http.ResponseWriter, r *http.Request) error {
		// Since this error will never be handled by Catcher,
		// it will be 500 internal server error by yacm.DefaultShutter
		return errors.New("unknown error")
	}))

	http.ListenAndServe(":8080", nil)
}
```


