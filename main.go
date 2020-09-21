package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	db *sql.DB
)

func handleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/Tasks", allTasks).Methods("GET")
	myRouter.HandleFunc("/Users", allUsers).Methods("GET")
	// myRouter.HandleFunc("/Tasks", addTasks).Methods("POST")
	// myRouter.HandleFunc("/Tasks", deleteTask).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}
func allUsers(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", "host=34.77.175.38 user=aimaxer dbname=locari_db sslmode=disable password=Maciek0808")
	jsons := simplejson.New()
	jsons2 := simplejson.New()

	if err != nil {
		panic(err)
	} else {
		fmt.Printf("DB call get all Users\n")
	}

	defer db.Close()

	rows, _ := db.Query(fmt.Sprintf("SELECT * FROM users"))
	for rows.Next() {
		var (
			Username string
			Password string
			Token    string
		)
		if err := rows.Scan(&Username, &Password, &Token); err != nil {
			log.Fatal(err)
		}

		jsons2.Set("Username", Username)
		jsons2.Set("Password", Password)
		jsons2.Set("Token", Token)
	}
	jsons.Set("Users", jsons2)

	payload, err := jsons.MarshalJSON()
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func allTasks(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", "host=34.77.175.38 user=aimaxer dbname=locari_db sslmode=disable password=Maciek0808")
	jsons := simplejson.New()

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
			content string
		)
		if err := rows.Scan(&title, &content); err != nil {
			log.Fatal(err)
		}

		jsons.Set(title, content)
	}

	payload, err := jsons.MarshalJSON()
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "delete task")
}

func addTasks(w http.ResponseWriter, r *http.Request) {
	// db.Exec(fmt.Sprintf("INSERT INTO tasks (title, description, content) VALUES ('%s', '%s', '%s')", r.FormValue("title"), r.FormValue("description"), r.FormValue("content")))

	db, err := sql.Open("postgres", "host=34.77.175.38 user=aimaxer dbname=locari_db sslmode=disable password=Maciek0808")

	if err != nil {
		panic(err)
	} else {
		fmt.Printf("DB connected")
	}

	defer db.Close()

	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("ioutil")
		panic(err.Error())
	}

	keyVal := make(map[string]string)

	json.Unmarshal(body, &keyVal)

	title := keyVal["title"]
	fmt.Printf(title, "\n")
	content := keyVal["content"]

	_, err = db.Exec(fmt.Sprintf("INSERT INTO tasks (title, content) VALUES ('%s', '%s')", title, content))
	if err != nil {
		fmt.Printf("\nexec\n")
		panic(err.Error())
	}

	fmt.Fprintf(w, "New task was inserted")
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage")
}

func main() {

	handleRequest()
}
