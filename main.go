package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type task struct {
	ID   int    `json:"id"`
	Info string `json:"info"`
}

var tasks []task

func main() {
	tasks = append(tasks, task{1, "Eat"})
	tasks = append(tasks, task{2, "Sleep"})
	tasks = append(tasks, task{3, "Rave"})
	fmt.Println(tasks)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", hello)
	router.HandleFunc("/task", createTask).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func createTask(w http.ResponseWriter, r *http.Request) {
	//example: curl -d '{"id":4, "info":"Learn"}' -H "Content-Type: application/json" -i -X POST http://localhost:8080/task
	var newTask task
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Error")
	}
	json.Unmarshal(reqBody, &newTask)
	tasks = append(tasks, newTask)
	w.WriteHeader(http.StatusCreated)
	fmt.Println(tasks)
}
