package main

import (
	"fmt"
	"html"
	"net/http"
)

func main() {
	port := "8070"
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "It's Alive!, %q", html.EscapeString(request.URL.Path))
	})
	fmt.Println("Listening on " + port)
	http.ListenAndServe(":"+port, nil)
}
