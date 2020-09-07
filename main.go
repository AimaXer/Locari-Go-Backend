package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	db *sql.DB
)

func allTasks(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", "host=34.77.175.38 user=aimaxer dbname=locari_db sslmode=disable password=Maciek0808")

	if err != nil {

		panic(err)

	} else {
		fmt.Printf("DB connected")
	}

	defer db.Close()
	rows, _ := db.Query(fmt.Sprintf("SELECT * FROM tasks"))
	for rows.Next() {
		var (
			title   string
			desc    string
			content string
		)
		if err := rows.Scan(&title, &desc, &content); err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "%s - %s - %s", title, desc, content)
	}
	json.NewEncoder(w).Encode(&rows)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "delete task")
}

func addTasks(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", "host=34.77.175.38 user=aimaxer dbname=locari_db sslmode=disable password=Maciek0808")

	if err != nil {
		panic(err)

	} else {
		fmt.Printf("DB connected\n")
	}

	defer db.Close()
	db.Exec(fmt.Sprintf("INSERT INTO tasks (title, description, content) VALUES ('%s', '%s', '%s')", r.FormValue("title"), r.FormValue("description"), r.FormValue("content")))
	fmt.Printf("Inserted \n")
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

	handleRequest()
}
