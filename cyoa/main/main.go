package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/joe-nguhi/gophercises/cyoa"
)

const port = "8080"

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseFiles("./temp.html"))
}

func main() {

	story, err := cyoa.GetStory("gopher.json")
	if err != nil {
		panic(err)
	}

	handler := cyoa.ArchHandler{Story: story, Page: tmpl}

	fmt.Println("Starting server on port", port)
	http.ListenAndServe(":"+port, handler)
}
