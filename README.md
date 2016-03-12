# YACM(八雲) - Yet Another Composable Middleware.

[![Build Status](https://travis-ci.org/morikuni/yacm.svg?branch=master)](https://travis-ci.org/morikuni/yacm)
[![GoDoc](https://godoc.org/github.com/morikuni/yacm?status.svg)](https://godoc.org/github.com/morikuni/yacm)

Simple, Lightweight and Composable HTTP Handler/Middleware.

## Install

```sh
go get github.com/morikuni/aec
```

## Design

YACM is designed as Middleware for `http.HandlerFunc` and work with `net/http`.  
There are 4 important types.

- http.HandlerFunc
- Middleware
- Service
- ErrorHandler

Flexible `http.HandlerFunc` can be made by composing these types.

```
Middleware + Middleware       => Middleware
Middleware + http.HandlerFunc => Service
Middleware + Service          => Service
Service    + ErrorHandler     => http.HandlerFunc
```

## Example

```go
package main

import (
	"log"
	"net/http"

	"github.com/morikuni/yacm"
)

func main() {
	// Middleware + Middleware => Middleware
	myMiddleware := yacm.Middleware(yacm.PanicFilter).And(yacm.MethodFilter(yacm.GET))

	// Middleware + http.HandlerFunc => Service
	myService := myMiddleware.For(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello\n"))
	})

	// Service + ErrorHandler => http.HandlerFunc
	myHandler := myService.Recover(func(w http.ResponseWriter, r *http.Request, err yacm.Error) {
		switch e := err.(type) {
		case yacm.PanicOccured:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error\n"))
			log.Println(e.Reason)
		case yacm.MethodNotAllowed:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method Not Allowed\n"))
			log.Printf("expect %v but %s\n", e.Expect, e.Method)
		}
	})

	http.HandleFunc("/hello", myHandler)
	http.ListenAndServe("127.0.0.1:12345", nil)
}
```

```sh
$ go run main/main.go &

$ curl -X "GET" "http://127.0.0.1:12345/hello"
Hello

$ curl -X "PUT" "http://127.0.0.1:12345/hello"
2016/02/28 00:19:30 expect [GET] but PUT # This is from main.go
Method Not Allowed
```

