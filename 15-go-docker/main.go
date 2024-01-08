package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from Docker!"))
		fmt.Fprintf(w, "\nHello, you've requested: %s\n", r.URL.Path)
	})

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hi from Docker!"))
		fmt.Fprintf(w, "\nHi, you've requested: %s\n", r.URL.Path)
	})

	http.ListenAndServe(":8080", nil)

}
