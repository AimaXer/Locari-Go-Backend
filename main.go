package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
  "github.com/jinzhu/gorm"
  "github.com/lib/pq"
  "github.com/rs/cors"
)

type Task struct {
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

type Tasks []Task

func allTasks(w http.ResponseWriter, r *http.Request) {
	Tasks := Tasks{
		Task{Title: "Test title", Desc: "Test Desc", Content: "Hello World"},
	}
	fmt.Println("Tasks Endpoint")
	json.NewEncoder(w).Encode(Tasks)
}

func deleteTask(w http.ResponseWriter, r *http.Request){
  fmt.Fprintf(w, "delete task")
}

func addTasks(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Post Tasks")
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage-karolina")
}

func handleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/Tasks", allTasks).Methods("GET")
	myRouter.HandleFunc("/Tasks", addTasks).Methods("POST")
	myRouter.HandleFunc("/Tasks", deleteTask).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	handleRequest()
}
