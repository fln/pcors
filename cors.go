// Package pcors provides implementation of fully permissive CORS middleware.
//
// This package is intended for web services that does not require a strict
// cross-origin resource sharing rules and would like to expose API endpoints to
// all origins, allow all methods, all request headers including credentials
// with minimal configuration. CORS pre-fligh request also hints client to
// cache the response for up to 24 hours.
//
// Example usage:
//	package main
//
//	import (
//		"net/http"
//
//		"github.com/fln/pcors"
//	)
//
//	func main() {
//		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//			w.Write([]byte("Hello world!"))
//		})
//		http.ListenAndServe(":8123", pcors.Default(http.DefaultServeMux))
//	}
package pcors

import (
	"net/http"
	"strings"
)

type cors struct {
	exposeHeaders string
	next          http.Handler
}

// ExposeHeaders generates permissive CORS middleware similar to Default one,
// but also exposes given list of response headers.
func ExposeHeaders(headers ...string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &cors{
			exposeHeaders: strings.Join(headers, ", "),
			next:          h,
		}
	}
}

// Default is a standard permissive CORS middleware. It will allow all
// cross-origin requests. However this middleware will not expose any of the
// response headers.
func Default(h http.Handler) http.Handler {
	return &cors{
		exposeHeaders: "",
		next:          h,
	}
}

// ServeHTTP implements http.Handler interface.
func (c *cors) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "" {
		c.next.ServeHTTP(w, r)
		return
	}

	// Common headers
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Add("Vary", "Origin")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	requestMethod := r.Header.Get("Access-Control-Request-Method")
	// Preflight request case
	if r.Method == "OPTIONS" && requestMethod != "" {
		w.Header().Set("Access-Control-Allow-Methods", requestMethod)
		w.Header().Add("Vary", "Access-Control-Request-Method")
		if headers := r.Header.Get("Access-Control-Request-Headers"); headers != "" {
			w.Header().Set("Access-Control-Allow-Headers", headers)
		}
		w.Header().Add("Vary", "Access-Control-Request-Headers")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Normal request case
	if c.exposeHeaders != "" {
		w.Header().Set("Access-Control-Expose-Headers", c.exposeHeaders)
	}
	c.next.ServeHTTP(w, r)
}
