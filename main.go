package main

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"

  "github.com/gorilla/mux"
  "github.com/jinzhu/gorm"
  //"github.com/lib/pq"
  //"github.com/rs/cors"
  "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
  _ "github.com/jinzhu/gorm/dialects/postgres"
)

type Task struct {
  gorm.Model
  Title   string
  Description    string
  Content string
}

var db *gorm.DB
var err error
type Tasks []Task

func allTasks(w http.ResponseWriter, r *http.Request) {
  var task Task
  db.Find(&task)
  json.NewEncoder(w).Encode(&task)
}

func deleteTask(w http.ResponseWriter, r *http.Request){
  fmt.Fprintf(w, "delete task")
}

func addTasks(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  taskAdd := Task{
    Title : vars["title"],
    Description : vars["description"],
    Content : vars["content"],
  }
  
  db.Create(&taskAdd)
}

func homePage(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Homepage-karolina")
}

func handleRequest() {
  myRouter := mux.NewRouter().StrictSlash(true)
  
  myRouter.HandleFunc("/", homePage)
  myRouter.HandleFunc("/Tasks", allTasks).Methods("GET")
  myRouter.HandleFunc("/Tasks", addTasks).Methods("POST")
  myRouter.HandleFunc("/Tasks", deleteTask).Methods("DELETE")
  log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
  
  db, err := sql.Open( "cloudsqlpostgres", "host=34.77.175.38 user=postgres dbname=locari_db sslmode=disable password=Maciek0808")

  if err != nil {

    panic(err)

  } else {
    fmt.Printf("DB connected")
  }
  
  defer db.Close()
  
  handleRequest()
}
