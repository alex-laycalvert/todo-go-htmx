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
	Todos []*todos.Todo
}

func main() {
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/api/todo", apiHandler)
	http.HandleFunc("/api/todo/", apiHandler)
	http.HandleFunc("/", indexHandler)
	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(res http.ResponseWriter, req *http.Request) {
	if err := renderPage(res, "index", PageData{Todos: todos.Todos()}); err != nil {
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
	if err := renderComponent(res, "todos", PageData{Todos: todos.Todos()}); err != nil {
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
