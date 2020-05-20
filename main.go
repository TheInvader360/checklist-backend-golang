package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type task struct {
	ID   int    `json:"id"`
	Info string `json:"info"`
	Done bool   `json:"done"`
}

var tasks []task

func main() {
	tasks = append(tasks, task{1001, "Eat", false})
	tasks = append(tasks, task{1002, "Sleep", false})
	tasks = append(tasks, task{1003, "Rave", false})
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/task", createTask).Methods("POST")
	router.HandleFunc("/tasks", readTasks).Methods("GET")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func createTask(w http.ResponseWriter, r *http.Request) {
	//curl -d '{"id":1004, "info":"Learn", "done":false}' -H "Content-Type: application/json" -i -X POST http://localhost:8080/task
	var reqTask task
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(reqBody, &reqTask)
	tasks = append(tasks, reqTask)
	w.WriteHeader(http.StatusCreated)
}

func readTasks(w http.ResponseWriter, r *http.Request) {
	//curl -i -X GET http://localhost:8080/tasks
	json.NewEncoder(w).Encode(tasks)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	//curl -d '{"id":1002, "info":"Drink", "done":true}' -H "Content-Type: application/json" -i -X PUT http://localhost:8080/tasks/1002
	updateTaskID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Println(err)
	}
	var reqTask task
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(reqBody, &reqTask)
	for index, task := range tasks {
		if task.ID == updateTaskID {
			tasks[index].Info = reqTask.Info
			tasks[index].Done = reqTask.Done
		}
	}
}
