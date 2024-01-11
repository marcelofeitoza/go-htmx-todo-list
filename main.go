package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

type Todo struct {
	Title  string
	Status bool
}

func main() {
	fmt.Println("Hello, world")

	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))

		todos := map[string][]Todo{
			"Todos": {
				{Title: "Task 1", Status: false},
				{Title: "Task 2", Status: true},
				{Title: "Task 3", Status: false},
				{Title: "Task 4", Status: true},
			},
		}

		tmpl.Execute(w, todos)
	}

	h2 := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		title := r.PostFormValue("title")
		status := false
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w, "todo-list-element", Todo{Title: title, Status: status})
	}

	http.HandleFunc("/", h1)
	http.HandleFunc("/add-todo/", h2)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
