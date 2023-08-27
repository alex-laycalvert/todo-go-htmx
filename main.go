package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/alex-laycalvert/todo/todo"
)

const (
	PAGES_DIR      = "templates/pages"
	LAYOUTS_DIR    = "templates/layouts"
	COMPONENTS_DIR = "templates/components"
)

type PageData struct {
	Todos []*todo.Todo
}

var todos []*todo.Todo

func main() {
	todos = []*todo.Todo{}
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/api/todo", apiHandler)
	http.HandleFunc("/api/todo/", apiHandler)
	http.HandleFunc("/", indexHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(res http.ResponseWriter, req *http.Request) {
	if err := renderPage(res, "index", PageData{Todos: todos}); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Printf("Error: %v\n", err)
	}
}

func apiHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		data, err := parseBody(req)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			log.Printf("Error: %v\n", err)
			return
		}
		description := data.Get("description")
		if len(description) > 0 {
			todos = append(todos, todo.New(description))
		}
	case "DELETE":
		params := strings.Split(req.URL.Path, "/")
		if len(params) == 0 {
			break
		}
		id := params[len(params)-1]
		index := -1
		for i, t := range todos {
			if t.Id.String() == id {
				index = i
				break
			}
		}
		if index < 0 {
			break
		}
		copy(todos[index:], todos[index+1:])
		todos[len(todos)-1] = nil
		todos = todos[:len(todos)-1]
	default:
		http.Error(res, "method not supported", http.StatusMethodNotAllowed)
		log.Printf("Error: %v method not supported", req.Method)
		return
	}
	if err := renderComponent(res, "todos", PageData{Todos: todos}); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Printf("Error: %v\n", err)
		return
	}
}

func parseBody(req *http.Request) (url.Values, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	return url.ParseQuery(string(body))
}

func renderComponent(res http.ResponseWriter, component string, data interface{}) error {
	component = component + ".html"
	tmpl, err := template.New(component).ParseFiles(COMPONENTS_DIR + "/" + component)
	if err != nil {
		return err
	}
	return tmpl.ExecuteTemplate(res, component, data)
}

func renderPage(res http.ResponseWriter, page string, data interface{}) error {
	if page[0] != '/' {
		page = "/" + page
	}
	tmpl := template.New("template")
	patterns := []string{LAYOUTS_DIR + "/base.html", PAGES_DIR + page + ".html", COMPONENTS_DIR + "/*.html"}
	var err error
	for _, pattern := range patterns {
		tmpl, err = tmpl.ParseGlob(pattern)
		if err != nil {
			return err
		}
	}
	return tmpl.ExecuteTemplate(res, "base", data)
}
