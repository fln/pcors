Permissive CORS
---------------

[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/fln/pcors)

Package pcors provides implementation of fully permissive
[CORS]("https://www.w3.org/TR/cors/") middleware.

This package is intended for web services that does not require a strict
cross-origin resource sharing rules and would like to expose API endpoints to
all origins, allow all methods, all request headers including credentials with
minimal configuration. CORS pre-fligh request also hints client to cache the
response for up to 24 hours.

Example usage:

```go
package main

import (
	"net/http"

	"github.com/fln/pcors"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})
	http.ListenAndServe(":8123", pcors.Default(http.DefaultServeMux))
}
```
