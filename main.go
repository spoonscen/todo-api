package main

import (
	"fmt"
	"html"
	"net/http"

	"gopkg.in/mgo.v2"
)

func main() {
	url := "localhost:27017"
	session, err := mgo.Dial(url)
	fmt.Println(err)

	c := session.DB("todo_api").C("todos")
	fmt.Println(c)

	port := "8070"
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "Its Alive!, %q", html.EscapeString(request.URL.Path))
	})
	fmt.Println("Listening on " + port)
	http.ListenAndServe(":"+port, nil)
}
