package urlshortener

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	yaml "gopkg.in/yaml.v2"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		redirectUrl, exists := pathsToUrls[r.URL.Path]

		if exists {
			log.Printf("Redirecting '%s' to '%s'...", r.URL.Path, redirectUrl)

			http.Redirect(w, r, redirectUrl, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

func YAMLFileHandler(yamlFilePath string, fallback http.Handler) (http.HandlerFunc, error) {
	// Open the yaml file containing paths and urls
	yamlFile, err := os.Open(yamlFilePath)
	if err != nil {
		return nil, err
	}
	// Defer closing the file after it has been read completely
	defer yamlFile.Close()

	// Read the entire file
	yml, err := ioutil.ReadAll(yamlFile)
	if err != nil {
		return nil, err
	}

	// Unmarshal the file contents into a golang map
	var parsedYaml []map[string]string
	if err := yaml.Unmarshal(yml, &parsedYaml); err != nil {
		return nil, err
	}

	// Process the map to something that the MapHandler can process
	pathsToUrls := make(map[string]string)
	for _, value := range parsedYaml {
		key := value["path"]
		pathsToUrls[key] = value["url"]
	}

	return MapHandler(pathsToUrls, fallback), nil
}

func JSONFileHandler(jsonFilePath string, fallback http.Handler) (http.HandlerFunc, error) {
	// Open the json file containing paths and urls
	jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		return nil, err
	}
	// Defer closing the file after it has been read completely
	defer jsonFile.Close()

	// Read the entire file
	jsn, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	// Unmarshal the file contents into a golang map
	var parsedYaml []map[string]string
	if err := json.Unmarshal(jsn, &parsedYaml); err != nil {
		return nil, err
	}

	// Process the map to something that the MapHandler can process
	pathsToUrls := make(map[string]string)
	for _, value := range parsedYaml {
		key := value["path"]
		pathsToUrls[key] = value["url"]
	}

	return MapHandler(pathsToUrls, fallback), nil
}

func MYSQLHandler(dsn string, fallback http.Handler) (http.HandlerFunc, error) {
	// Create a variable to hold the mysql connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// Make sure the correctly close the database connection after we're done
	defer db.Close()

	// Connect to the database and query all data from the paths table
	rows, err := db.Query("select `path`, `url` from `paths`")
	if err != nil {
		return nil, err
	}
	// Close the result reset after we extracted our data
	defer rows.Close()

	var path, url string
	pathsToUrls := make(map[string]string)

	// Loop through each result, and assign it to our result map
	for rows.Next() {
		err := rows.Scan(&path, &url)
		if err != nil {
			return nil, err
		}

		pathsToUrls[path] = url
	}

	return MapHandler(pathsToUrls, fallback), nil
}
