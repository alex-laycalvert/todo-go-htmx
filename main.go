package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/alex-laycalvert/todo/todos"
)

const (
	PAGES_DIR      = "templates/pages"
	LAYOUTS_DIR    = "templates/layouts"
	COMPONENTS_DIR = "templates/components"
)

type PageData struct {
	Todos []todos.Todo
}

func main() {
	log.Println("Initializing Database")
	err := todos.Init()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v\n", err.Error())
	}
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/api/todo", apiHandler)
	http.HandleFunc("/api/todo/", apiHandler)
	http.HandleFunc("/", indexHandler)
	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(res http.ResponseWriter, req *http.Request) {
	todos, err := todos.Todos()
	if err != nil {
		handleError(res, err)
		return
	}
	if err := renderPage(res, "index", PageData{Todos: todos}); err != nil {
		handleError(res, err)
	}
}

func apiHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		data, err := parseBody(req)
		if err != nil {
			handleError(res, err)
			return
		}
		todos.Add(data.Get("description"))
	case "DELETE":
		params := strings.Split(req.URL.Path, "/")
		if len(params) == 0 {
			break
		}
		todos.Remove(params[len(params)-1])
	default:
		http.Error(res, "method not supported", http.StatusMethodNotAllowed)
		log.Printf("Error: %v method not supported", req.Method)
		return
	}
	todos, err := todos.Todos()
	if err != nil {
		handleError(res, err)
		return
	}
	if err := renderComponent(res, "todos", PageData{Todos: todos}); err != nil {
		handleError(res, err)
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

func handleError(res http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	http.Error(res, err.Error(), http.StatusInternalServerError)
	log.Printf("Error: %v\n", err)
}
