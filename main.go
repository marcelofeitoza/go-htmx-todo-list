package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Todo struct {
	Id     int
	Title  string
	Status bool
}

var (
	currentId = 2
	todos     = map[string][]Todo{
		"Todos": {
			{Id: 1, Title: "Go to colleg", Status: false},
			{Id: 2, Title: "Play with dog", Status: true},
		},
	}
)

func main() {
	fmt.Println("Hello, world")

	h1 := func(w http.ResponseWriter, _ *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, todos)
	}

	h2 := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		title := r.PostFormValue("title")
		currentId += 1
		todo := Todo{Id: currentId, Title: title, Status: false}
		todos["Todos"] = append(todos["Todos"], todo)
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w, "todo-list-element", todo)
	}

	h3 := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		r.ParseForm()
		taskIDStr := r.FormValue("taskID")
		taskID, err := strconv.Atoi(taskIDStr)
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		for i, task := range todos["Todos"] {
			if task.Id == taskID {
				todos["Todos"][i].Status = !task.Status
				tmpl, err := template.ParseFiles("index.html")
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				err = tmpl.ExecuteTemplate(w, "todo-list-element", todos["Todos"][i])
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		}
		http.Error(w, "Todo item not found", http.StatusNotFound)
	}

	h4 := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		r.ParseForm()
		taskIDStr := r.FormValue("taskID")
		taskID, err := strconv.Atoi(taskIDStr)
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		for i, task := range todos["Todos"] {
			if task.Id == taskID {
				todos["Todos"] = append(todos["Todos"][:i], todos["Todos"][i+1:]...)
				w.WriteHeader(http.StatusOK)
				return
			}
		}
		http.Error(w, "Todo item not found", http.StatusNotFound)
	}

	http.HandleFunc("/", h1)
	http.HandleFunc("/add-todo/", h2)
	http.HandleFunc("/toggle-todo/", h3)
	http.HandleFunc("/delete-todo/", h4)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
