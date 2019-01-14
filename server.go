package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	validator "gopkg.in/validator.v2"
	yaml "gopkg.in/yaml.v2"
)

// Appdata is a struct for app info
type Appdata struct {
	Title       string       `yaml:"title" validate:"nonzero"`
	Version     string       `yaml:"version" validate:"nonzero"`
	Maintainers []Maintainer `yaml: ", inline" validate:"nonzero"`
	Company     string       `yaml:"company" validate:"nonzero"`
	Website     string       `yaml:"website" validate:"nonzero"`
	Source      string       `yaml:"source" validate:"nonzero"`
	License     string       `yaml:"license" validate:"nonzero"`
	Description string       `yaml:"description" validate:"nonzero"`
}

// Maintainer is a struct for maintainer info
type Maintainer struct {
	Name  string `yaml:"name" validate:"nonzero"`
	Email string `yaml:"email" validate:"regexp=^[0-9a-z]+@[0-9a-z]+(\\.[0-9a-z]+)+$"`
}

var dataStore []Appdata

// GetHandler retrieves all records
func GetHandler(w http.ResponseWriter, req *http.Request) {
	yamlData, err := yaml.Marshal(dataStore)
	if err != nil {
		fmt.Println(err)
	}
	result := string(yamlData)
	w.Header().Set("Content-Type", "application/x-yaml")
	fmt.Fprintf(w, "%s", result)
}

// QueryHandler retrieves records by title
func QueryHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var result []string
	for _, item := range dataStore {
		if strings.ToLower(item.Title) == strings.ToLower(params["title"]) {
			yamlItem, err := yaml.Marshal(item)
			if err != nil {
				fmt.Println(err)
			}
			result = append(result, string(yamlItem))
		}
	}
	if result != nil {
		w.Header().Set("Content-Type", "application/x-yaml")
		fmt.Fprintf(w, "%s", result)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

// PostHandler creates a new record in the data store
func PostHandler(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newData := Appdata{}
	err = yaml.UnmarshalStrict(body, &newData)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if errs := validator.Validate(newData); errs != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "missing or invalid required field")
		return
	}
	dataStore = append(dataStore, newData)
	yamlData, err := yaml.Marshal(newData)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/x-yaml")
	w.WriteHeader(http.StatusCreated)
	result := string(yamlData)
	fmt.Fprintf(w, "%v", result)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/applications", GetHandler).Methods("GET")
	router.HandleFunc("/applications", PostHandler).Methods("POST")
	router.HandleFunc("/applications/{title}", QueryHandler).Methods("GET")

	server := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
