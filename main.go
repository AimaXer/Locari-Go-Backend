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

type bodyMessageUsers struct {
	Usr  string `json:"usr"`
	Pass string `json:"pass"`
}

type bodyMessageGetTasks struct {
	UserToken string `json:"usertoken"`
}

type bodyMessageSendTasks struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type bodyMessageAddTask struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	UserToken string `json:"userToken"`
}

func handleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/getTasks", allUserTasks).Methods("POST")
	myRouter.HandleFunc("/Auth", authUser).Methods("POST")
	myRouter.HandleFunc("/addTask", addTasks).Methods("POST")
	// myRouter.HandleFunc("/Tasks", deleteTask).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}
func authUser(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres dbname=locari-db sslmode=disable password=Maciek0808")
	jsons := simplejson.New()
	jsons2 := simplejson.New()

	if err != nil {
		panic(err)
	} else {
		fmt.Printf("DB call get all Users ALL\n")
	}

	defer db.Close()
	rows, _ := db.Query(fmt.Sprintf("SELECT * FROM users.users"))

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("ioutil")
		panic(err.Error())
	}
	var msg bodyMessageUsers

	json.Unmarshal([]byte(body), &msg)

	userapp := msg.Usr
	passapp := msg.Pass
	found := false

	for rows.Next() {
		var (
			username string
			password string
			token    string
		)
		if err := rows.Scan(&token, &username, &password); err != nil {
			log.Fatal(err)
		}

		if userapp == username && passapp == password {
			jsons2.Set("Token", token)
			// fmt.Printf(userapp + " - " + username + " - " + passapp + " - " + password + "\n")
			fmt.Printf(token)
			found = true
		}
	}
	if !found {
		jsons2.Set("Token", "")
	}
	jsons.Set("Users", jsons2)
	payload, err := jsons.MarshalJSON()
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func allUserTasks(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres dbname=locari-db sslmode=disable password=Maciek0808")
	jsons := simplejson.New()

	if err != nil {
		panic(err)
	} else {
		fmt.Printf("DB connected get tasks" + "\n")
	}
	defer db.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("ioutil" + "\n")
		panic(err.Error())
	}
	var msg bodyMessageGetTasks
	json.Unmarshal([]byte(body), &msg)
	userT := msg.UserToken
	var tasks []bodyMessageSendTasks

	rows, _ := db.Query(fmt.Sprintf("SELECT * FROM users.tasks"))
	for rows.Next() {
		var (
			id        string
			title     string
			content   string
			usertoken string
		)
		if err := rows.Scan(&id, &title, &content, &usertoken); err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("dziala" + "\n")

		taskojb := bodyMessageSendTasks{ID: id, Title: title, Content: content}
		// task, _ := json.Marshal(taskj)

		if userT == usertoken {
			tasks = append(tasks, taskojb)
		}
	}

	jsons.Set("Tasks", tasks)
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

	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres dbname=locari-db sslmode=disable password=Maciek0808")

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

	var msg bodyMessageAddTask

	json.Unmarshal([]byte(body), &msg)

	id := msg.ID
	title := msg.Title
	content := msg.Content
	userToken := msg.UserToken

	_, err = db.Exec(fmt.Sprintf("INSERT INTO users.tasks (id, title, content, usertoken) VALUES ('%s', '%s', '%s', '%s')", id, title, content, userToken))
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
