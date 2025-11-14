package cyoa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type ArchHandler struct {
	Story
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
		printStory(w, a.Story["intro"])
	}

	for k, v := range a.Story {
		fmt.Printf("key: %s, url: %s\n", k, arc)
		if k == arc {
			printStory(w, v)
			return
		}
	}

	fmt.Fprintf(w, "<h1>%s</h1>", "Not Found")
}

func printStory(w http.ResponseWriter, arc StoryArc) {
	fmt.Fprintf(w, "<h1>%s</h1>", arc.Title)
	for _, line := range arc.Story {
		fmt.Fprintf(w, "<p>%s</p>", line)
	}
	for _, option := range arc.Options {
		fmt.Fprintf(w, "<div><a href=\"%s\">%s</a></div>", option.Arc, option.Text)
	}
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
