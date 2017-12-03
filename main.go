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
	Created   int64
}

func errorHandler(err error) {
	if err != nil {
		panic(err)
	}
}

var c *mgo.Collection

func AllTodos(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var allTodos []Todo
	err := c.Find(nil).Sort("Created").All(&allTodos)
	json, err := json.Marshal(allTodos)
	errorHandler(err)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func CompleteTodos(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var allTodos []Todo
	err := c.Find(bson.M{"completed": true}).Sort("Created").All(&allTodos)
	json, err := json.Marshal(allTodos)
	errorHandler(err)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func IncompleteTodos(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var allTodos []Todo
	err := c.Find(bson.M{"completed": false}).Sort("Created").All(&allTodos)
	json, err := json.Marshal(allTodos)
	errorHandler(err)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func GetOneByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

type Update struct {
	Text      string
	Completed bool
}

func UpdateOneByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	update := Update{}
	todo := Todo{}
	decoder := json.NewDecoder(r.Body)
	id := bson.ObjectIdHex(ps.ByName("id"))

	err := decoder.Decode(&update)
	errorHandler(err)

	err = c.FindId(id).One(&todo)
	errorHandler(err)

	if todo.Completed != update.Completed {
		todo.Completed = update.Completed
	}

	if len(update.Text) > 0 {
		todo.Text = update.Text
	}

	err = c.UpdateId(id, todo)
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

func ReplaceOneByID(w http.ResponseWriter, r *http.Request, _ps httprouter.Params) {
	update := Todo{}
	todo := Todo{}
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&update)
	errorHandler(err)

	err = c.FindId(update.ID).One(&todo)
	errorHandler(err)

	err = c.UpdateId(update.ID, update)
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
	currentTime := time.Now()
	todo := &Todo{ID: id, Text: text, Completed: false, Created: currentTime.Unix()}
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
	router.GET("/api/todos-complete", CompleteTodos)
	router.GET("/api/todos-incomplete", IncompleteTodos)
	router.POST("/api/todos", CreateTodo)
	router.PUT("/api/todos", ReplaceOneByID)
	router.GET("/api/todos/:id", GetOneByID)
	router.PATCH("/api/todos/:id", UpdateOneByID)
	fmt.Println(http.ListenAndServe(":"+port, router))
}

func ConnectToDb() {
	url := "mongo"
	session, err := mgo.Dial(url)
	errorHandler(err)
	c = session.DB("todo_api").C("todos")
}

func InsetMockDocs() {
	currentTime := time.Now()
	err := c.Insert(
		&Todo{
			Text:      "Test1",
			Completed: false,
			Created:   currentTime.Unix()},
		&Todo{
			Text:      "Test2",
			Completed: false,
			Created:   currentTime.Unix()})
	errorHandler(err)
}

func main() {
	ConnectToDb()
	InsetMockDocs()
	RunServer()
}
