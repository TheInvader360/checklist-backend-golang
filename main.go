package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type task struct {
	ID   int    `json:"id"`
	Info string `json:"info"`
	Done bool   `json:"done"`
}

var tasks []task
var nextTaskID int

func main() {
	tasks = append(tasks, task{1001, "Eat", false})
	tasks = append(tasks, task{1002, "Sleep", false})
	tasks = append(tasks, task{1003, "Rave", false})
	nextTaskID = 1004

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/task", createTask).Methods("POST")
	router.HandleFunc("/tasks", readTasks).Methods("GET")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	headers := handlers.AllowedHeaders([]string{"Content-Type"})
	methods := handlers.AllowedMethods([]string{"POST", "GET", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(headers, methods, origins)(router)))
}

func createTask(w http.ResponseWriter, r *http.Request) {
	//curl -d '{"info":"Learn", "done":false}' -H "Content-Type: application/json" -i -X POST http://localhost:8080/task
	var reqTask task
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(reqBody, &reqTask)

	reqTask.ID = nextTaskID
	nextTaskID++

	conflict := false
	for _, task := range tasks {
		if task.ID == reqTask.ID {
			conflict = true
		}
	}

	if conflict {
		w.WriteHeader(http.StatusConflict)
	} else {
		tasks = append(tasks, reqTask)
		w.WriteHeader(http.StatusCreated)
		fmt.Println("Created")
		json.NewEncoder(w).Encode(reqTask)
	}
}

func readTasks(w http.ResponseWriter, r *http.Request) {
	//curl -i -X GET http://localhost:8080/tasks
	json.NewEncoder(w).Encode(tasks)
	fmt.Println("Read")
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	//curl -d '{"info":"Drink", "done":true}' -H "Content-Type: application/json" -i -X PUT http://localhost:8080/tasks/1002
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

	found := false

	for index, task := range tasks {
		if task.ID == updateTaskID {
			found = true
			tasks[index].Info = reqTask.Info
			tasks[index].Done = reqTask.Done
			fmt.Println("Updated")
		}
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	//curl -i -X DELETE http://localhost:8080/tasks/1003
	deleteTaskID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Println(err)
	}

	found := false

	for index, task := range tasks {
		if task.ID == deleteTaskID {
			found = true
			tasks = append(tasks[:index], tasks[index+1:]...)
			fmt.Println("Deleted")
		}
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
	}
}
