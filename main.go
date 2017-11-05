package main

import (
	"fmt"
	"html"
	"net/http"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Todo struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Text      string
	Completed bool
	Timestamp time.Time
}

func errorHandler(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	url := "mongo"
	session, err := mgo.Dial(url)
	errorHandler(err)

	c := session.DB("todo_api").C("todos")
	err = c.Insert(&Todo{Text: "Test", Completed: false, Timestamp: time.Now()})
	errorHandler(err)

	port := "8070"
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "Its Alive!, %q", html.EscapeString(request.URL.Path))
	})
	fmt.Println("Listening on " + port)
	http.ListenAndServe(":"+port, nil)
}
