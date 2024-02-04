package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Todo struct {
	Title     string
	Id        int
	Completed bool
}

type TodoUpdate struct {
	TaskID int `json:"taskID"`
}

var (
	currentId = 0
	todos     = map[string][]Todo{
		"Todos": {},
	}
)

var db *sqlx.DB

func init() {
	var err error

	for i := 0; i < 5; i++ {
		db, err = sqlx.Connect("postgres", "user=postgres_db dbname=postgres_db password=postgres_db sslmode=disable host=postgres")
		if err == nil {
			break
		}
		log.Println("Retrying to connect to the database...")
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalln(err)
	}

	db.Select(&todos, "SELECT id, title, completed FROM todos ORDER BY id")
}

func main() {
	fmt.Println("Listening at http://127.0.0.1:8080/")

	defer db.Close()

	render := func(w http.ResponseWriter, _ *http.Request) {
		var todos []Todo
		err := db.Select(&todos, "SELECT id, title, completed FROM todos ORDER BY id")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, map[string]interface{}{"Todos": todos})
	}

	addTodo := func(w http.ResponseWriter, r *http.Request) {
		title := r.PostFormValue("title")

		var newTodo Todo
		err := db.QueryRow("INSERT INTO todos (title, completed) VALUES ($1, $2) RETURNING id, title, completed", title, false).Scan(&newTodo.Id, &newTodo.Title, &newTodo.Completed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w, "todo-list-element", newTodo)
	}

	updateTodo := func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		taskIDStr := r.FormValue("taskID")
		taskID, err := strconv.Atoi(taskIDStr)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		res, err := db.Exec("UPDATE todos SET completed = NOT completed WHERE id = $1", taskID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if rowsAffected == 0 {
			http.Error(w, "Todo item not found", http.StatusNotFound)
			return
		}

		var updatedTodo Todo
		err = db.Get(&updatedTodo, "SELECT * FROM todos WHERE id = $1", taskID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl := template.Must(template.ParseFiles("index.html"))
		err = tmpl.ExecuteTemplate(w, "todo-list-element", updatedTodo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	deleteTodo := func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		taskIDStr := r.URL.Path[len("/delete-todo/"):]

		taskID, err := strconv.Atoi(taskIDStr)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		res, err := db.Exec("DELETE FROM todos WHERE id = $1", taskID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if rowsAffected > 0 {
			w.WriteHeader(http.StatusOK)
			return
		}
		if rowsAffected == 0 {
			http.Error(w, "Todo item not found", http.StatusNotFound)
			return
		}
	}

	http.HandleFunc("/", render)
	http.HandleFunc("/add-todo/", addTodo)
	http.HandleFunc("/toggle-todo/", updateTodo)
	http.HandleFunc("/delete-todo/", deleteTodo)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
