package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"./urlshortener"
)

var yamlFilePath string
var jsonFilePath string
var mysqlDsn string

func init() {
	// Parse command line arguments
	flag.StringVar(&yamlFilePath, "yaml", "paths.yml", "a file that contains paths in yaml format")
	flag.StringVar(&jsonFilePath, "json", "paths.json", "a file that contains paths in json format")
	flag.StringVar(&mysqlDsn, "dsn", "root:toor@tcp(127.0.0.1:3306)/urlshortener", "a dsn for the mysql server which tells us how to connect to it")

	flag.Parse()
}

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshortener.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the fallback
	yamlFileHandler, err := urlshortener.YAMLFileHandler(yamlFilePath, mapHandler)
	if err != nil {
		log.Fatal(err)
	}

	// Build the JSONFileHandler using the YAMLHandler as the fallback
	jsonFileHandler, err := urlshortener.JSONFileHandler(jsonFilePath, yamlFileHandler)
	if err != nil {
		log.Fatal(err)
	}

	mysqlHandler, err := urlshortener.MYSQLHandler(mysqlDsn, jsonFileHandler)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting the server on :8080")

	// Start listening...
	http.ListenAndServe(":8080", mysqlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Create a landing page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world!")
	})

	return mux
}
