package main

import (
	"net/http"

	"github.com/joe-nguhi/gophercises/cyoa"
)

func main() {

	story, err := cyoa.GetStory("gopher.json")
	if err != nil {
		panic(err)
	}

	handler := cyoa.ArchHandler{Story: story}

	http.ListenAndServe(":8080", handler)
}
