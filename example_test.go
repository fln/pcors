package pcors

import (
	"net/http"
)

func ExampleDefault() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})
	http.ListenAndServe(":8123", Default(http.DefaultServeMux))
}

func ExampleExposeHeaders() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-First", "first exposed")
		w.Header().Set("X-Second", "second exposed")
		w.Header().Set("X-Third", "not exposed")
		w.Write([]byte("Hello world!"))
	})
	http.ListenAndServe(
		":8123",
		ExposeHeaders("X-First", "X-Second")(http.DefaultServeMux),
	)
}
