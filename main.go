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

var DB_IP = "127.0.0.1"

type bodyMessageUsers struct {
	Usr  string `json:"usr"`
	Pass string `json:"pass"`
}

type bodyMessageUserInfo struct {
	UserToken string `json:"userToken"`
}

type bodyMessageGetTasks struct {
	UserToken string `json:"usertoken"`
}

type bodyMessageSendTasks struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Lat     string `json:"lat"`
	Long    string `json:"long"`
}

type bodyMessageAddTask struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	UserToken string `json:"userToken"`
	Lat       string `json:"lat"`
	Long      string `json:"long"`
}

func handleRequest() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/getTasks", allUserTasks).Methods("POST")
	myRouter.HandleFunc("/Auth", authUser).Methods("POST")
	myRouter.HandleFunc("/addTask", addTasks).Methods("POST")
	myRouter.HandleFunc("/delTask", deleteTask).Methods("POST")
	myRouter.HandleFunc("/getUserInfo", getUserInfo).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}
func getUserInfo(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", "host="+DB_IP+" port=5432 user=postgres dbname=locari_db sslmode=disable password=postgres")
	jsons := simplejson.New()

	defer db.Close()

	if err != nil {
		panic(err)
	} else {
		fmt.Printf("DB call get User info\n")
	}

	rows, _ := db.Query(fmt.Sprintf("SELECT * FROM users.users"))
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Printf("ioutil")
		panic(err.Error())
	}
	var msg bodyMessageUserInfo

	json.Unmarshal([]byte(body), &msg)

	userToken := msg.UserToken

	for rows.Next() {
		var (
			username string
			password string
			token    string
			email    string
		)
		if err := rows.Scan(&token, &username, &password, &email); err != nil {
			log.Fatal(err)
		}
		if userToken == token {
			jsons.Set("Token", token)
		}
	}
	jsons.Set("User", jsons)
	payload, err := jsons.MarshalJSON()
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func authUser(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", "host="+DB_IP+" port=5432 user=postgres dbname=locari_db sslmode=disable password=postgres")
	jsons := simplejson.New()
	jsons2 := simplejson.New()

	if err != nil {
		panic(err)
	} else {
		fmt.Printf("DB call auth user\n")
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
			email    string
		)
		if err := rows.Scan(&token, &username, &password, &email); err != nil {
			log.Fatal(err)
		}
		if userapp == username && passapp == password {
			jsons2.Set("Token", token)
			// fmt.Printf(userapp + " - " + username + " - " + passapp + " - " + password + "\n")
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
	db, err := sql.Open("postgres", "host="+DB_IP+" port=5432 user=postgres dbname=locari_db sslmode=disable password=postgres")
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
			lat       string
			long      string
		)
		if err := rows.Scan(&id, &title, &content, &usertoken, &lat, &long); err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("dziala" + "\n")

		taskojb := bodyMessageSendTasks{ID: id, Title: title, Content: content, Lat: lat, Long: long}
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
	db, err := sql.Open("postgres", "host="+DB_IP+" port=5432 user=postgres dbname=locari_db sslmode=disable password=postgres")

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
	userToken := msg.UserToken

	_, err = db.Exec(fmt.Sprintf("DELETE FROM users.tasks WHERE id = '%s' AND usertoken = '%s'", id, userToken))
	if err != nil {
		fmt.Printf("\nexec\n")
		panic(err.Error())
	}

	fmt.Fprintf(w, fmt.Sprintf("Task (ID = %s, usrT = %s) was deleted", id, userToken))
}

func addTasks(w http.ResponseWriter, r *http.Request) {
	// db.Exec(fmt.Sprintf("INSERT INTO tasks (title, description, content) VALUES ('%s', '%s', '%s')", r.FormValue("title"), r.FormValue("description"), r.FormValue("content")))

	db, err := sql.Open("postgres", "host="+DB_IP+" port=5432 user=postgres dbname=locari_db sslmode=disable password=postgres")

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
	lat := msg.Lat
	long := msg.Long

	_, err = db.Exec(fmt.Sprintf("INSERT INTO users.tasks (id, title, content, usertoken, lat, long) VALUES ('%s', '%s', '%s', '%s', '%s', '%s')", id, title, content, userToken, lat, long))
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
