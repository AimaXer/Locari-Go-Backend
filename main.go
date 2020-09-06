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
  var tasks Task
  db.Find(&task)
	json.NewEncoder(w).Encode(&task)
}

func deleteTask(w http.ResponseWriter, r *http.Request){
  fmt.Fprintf(w, "delete task")
}

func addTasks(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  taskAdd = Task{
    "Title" : vars["title"],
    "desc" : vars["description"],
    "content" : vars["content"],
  }
  
  db.Create(&taskAdd)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage-karolina")
}

func handleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)
	
	db, err = gorm.Open( "postgres", "host=34.77.175.38 port=5432 user=postgres dbname=postgres sslmode=disable password=Maciek0808")

  if err != nil {

    panic("failed to connect database")

  }
  
  defer db.close()
  
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/Tasks", allTasks).Methods("GET")
	myRouter.HandleFunc("/Tasks", addTasks).Methods("POST")
	myRouter.HandleFunc("/Tasks", deleteTask).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	handleRequest()
}
