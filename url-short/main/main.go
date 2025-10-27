package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	urlshort "github.com/joe-nguhi/gophercises/url-short"
)

var yamlFlag string
var jsonFlag string

func init() {
	flag.StringVar(&yamlFlag, "yaml", "paths.yml", "a yaml file in the format 'path: url'")
	flag.StringVar(&jsonFlag, "json", "paths.json", "a json file in the format 'path: url'")
	flag.Parse()
}

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the fallback
	yaml, err := getFileData(yamlFlag)
	if err != nil {
		panic(err)
	}

	yamlHandler, err := urlshort.YAMLHandler(yaml, mapHandler)
	if err != nil {
		panic(err)
	}

	// Build the JSONHandler using the yamlHandler as the fallback
	json, err := getFileData(jsonFlag)
	if err != nil {
		panic(err)
	}

	jsonHandler, err := urlshort.JSONHandler(json, yamlHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func getFileData(file string) ([]byte, error) {
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

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
