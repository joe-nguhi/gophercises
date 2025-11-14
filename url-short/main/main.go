package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	urlshort "github.com/joe-nguhi/gophercises/url-short"
	_ "github.com/mattn/go-sqlite3"
)

const dbName = "./paths.db"

var yamlFlag string
var jsonFlag string

func init() {
	flag.StringVar(&yamlFlag, "yaml", "paths.yml", "a yaml file in the format 'path: url'")
	flag.StringVar(&jsonFlag, "json", "paths.json", "a json file in the format 'path: url'")
	flag.Parse()
}

func main() {
	db, err := initDB()
	if err != nil {
		panic(err)
	}

	defer db.Close()

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

	sqlHandler, err := urlshort.SQLHandler(db, jsonHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", sqlHandler)
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

func initDB() (*sql.DB, error) {
	os.Remove(dbName)

	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}

	// Create table
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS paths (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL ,
			value TEXT NOT NULL
		);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return nil, err
	}

	// Insert Data
	sqlStmt = `
		INSERT INTO paths(name, value) 
		VALUES("/my-portfolio","https://joe-nguhi.netlify.app/"), ("/my-github","https://github.com/joe-nguhi")
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q \n", err)
		return nil, err
	}

	return db, nil
}
