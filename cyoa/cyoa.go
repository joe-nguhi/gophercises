package cyoa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type ArchHandler struct {
	Story
	Page *template.Template
}

type StoryArc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

type Story map[string]StoryArc

func (a ArchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path
	arc := strings.Replace(path, "/", "", 1)

	if arc == "intro" || arc == "" {
		a.printStory(w, a.Story["intro"])
		return
	}

	for k, v := range a.Story {
		if k == arc {
			a.printStory(w, v)
			return
		}
	}

	fmt.Fprintf(w, "<h1>%s</h1>", "Not Found")
}

func (a ArchHandler) printStory(w http.ResponseWriter, data StoryArc) {

	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	err := a.Page.Execute(w, data)
	check(err)

}

func GetStory(file string) (Story, error) {
	data, err := getBytes(file)
	if err != nil {
		return nil, err
	}

	var a Story
	d := json.NewDecoder(bytes.NewReader(data))
	err = d.Decode(&a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func getBytes(file string) ([]byte, error) {

	f, err := os.Open(file)
	if err != nil {
		return []byte(""), err
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		return []byte(""), err
	}

	return bytes, nil

}
