package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	validator "gopkg.in/validator.v2"
	yaml "gopkg.in/yaml.v2"
)

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

type Maintainer struct {
	Name  string `yaml:"name" validate:"nonzero"`
	Email string `yaml:"email" validate:"regexp=^[0-9a-z]+@[0-9a-z]+(\\.[0-9a-z]+)+$"`
}

var ds []Appdata

func GetHandler(w http.ResponseWriter, req *http.Request) {

	yamlData, err := yaml.Marshal(ds)
	if err != nil {
		fmt.Println(err)
	}
	result := string(yamlData)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "results:", result)

}

func QueryHandler(w http.ResponseWriter, req *http.Request) {
	queryParams := req.URL.Query()
	for _, item := range ds {
		if item.Title == queryParams["title"][0] {
			yamlItem, err := yaml.Marshal(item)
			if err != nil {
				fmt.Println(err)
			}
			result := string(yamlItem)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "result:\n%s", result)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

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

	ds = append(ds, newData)
	yamlData, err := yaml.Marshal(newData)
	if err != nil {
		fmt.Println(err)
	}
	w.WriteHeader(http.StatusCreated)
	result := string(yamlData)
	fmt.Fprintf(w, "--- entry created:\n%v\n\n", result)
}

func seedData() {
	d1 := Appdata{"Valid App 2", "1.0.1", []Maintainer{{"Sam Brown", "Sam@gmail.com"}, {"Fam Brown", "Fam@gmail.com"}}, "1", "2", "3", "4", "5"}
	ds = append(ds, d1)
}

func main() {
	seedData()

	router := mux.NewRouter()
	router.HandleFunc("/", GetHandler).Methods("GET")
	router.HandleFunc("/applications", PostHandler).Methods("POST")
	router.HandleFunc("/applications", QueryHandler).Methods("GET")
	router.Queries("title", "version", "company")

	server := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
