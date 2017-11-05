package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
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

var c *mgo.Collection

func AllTodos(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var allTodos []Todo
	err := c.Find(nil).Sort("-timestamp").All(&allTodos)
	json, err := json.Marshal(allTodos)
	errorHandler(err)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func GetOneById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	todo := Todo{}
	id := bson.ObjectIdHex(ps.ByName("id"))
	err := c.Find(bson.M{"_id": id}).One(&todo)
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprint(w, err.Error())
	} else {
		json, err := json.Marshal(todo)
		errorHandler(err)
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	}
}

func CreateTodo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var err error
	var text string

	decoder := json.NewDecoder(r.Body)
	id := bson.NewObjectId()

	err = decoder.Decode(&text)
	errorHandler(err)

	todo := &Todo{ID: id, Text: text, Completed: false, Timestamp: time.Now()}
	err = c.Insert(todo)
	errorHandler(err)

	json, err := json.Marshal(todo)
	errorHandler(err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Todo API !")
}

func RunServer() {
	port := "8070"
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/api/todos", AllTodos)
	router.GET("/api/todos/:id", GetOneById)
	router.POST("/api/todos/create", CreateTodo)
	fmt.Println(http.ListenAndServe(":"+port, router))
}

func ConnectToDb() {
	url := "mongo"
	session, err := mgo.Dial(url)
	errorHandler(err)
	c = session.DB("todo_api").C("todos")
}

func InsetMockDocs() {
	err := c.Insert(&Todo{Text: "Test1", Completed: false, Timestamp: time.Now()}, &Todo{Text: "Test2", Completed: false, Timestamp: time.Now()})
	errorHandler(err)
}

func main() {
	ConnectToDb()
	InsetMockDocs()
	RunServer()
}
