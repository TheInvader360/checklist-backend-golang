package main

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestBootstrapData(t *testing.T) {
	bootstrapData()
	assert.Equal(t, 3, len(tasks))
	assert.Equal(t, 1004, nextTaskID)
	assert.Equal(t, task{1001, "Eat", true}, tasks[0])
	assert.Equal(t, task{1002, "Sleep", false}, tasks[1])
	assert.Equal(t, task{1003, "Rave", false}, tasks[2])
}

func TestCreateTask2(t *testing.T) {
	//curl -d '{"info":"Learn", "done":false}' -H "Content-Type: application/json" -i -X POST http://localhost:8080/task
	type test struct {
		body, expectedInfo                                                      string
		nextTaskID, expectedCode, expectedTasks, expectedID, expectedNextTaskID int
		expectedDone                                                            bool
	}
	tests := []test{
		{
			body:               "{\"info\":\"Learn\", \"done\":false}",
			nextTaskID:         1004,
			expectedCode:       201,
			expectedTasks:      4,
			expectedID:         1004,
			expectedInfo:       "Learn",
			expectedDone:       false,
			expectedNextTaskID: 1005,
		},
		{
			body:               "{\"info\":\"Learn\", \"done\":false}",
			nextTaskID:         1001, // force conflict
			expectedCode:       409,
			expectedTasks:      3,
			expectedID:         1004,
			expectedInfo:       "Learn",
			expectedDone:       false,
			expectedNextTaskID: 1002,
		},
	}
	for _, tc := range tests {
		bootstrapData()
		nextTaskID = tc.nextTaskID
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "http://localhost:8080/task", strings.NewReader(tc.body))
		createTask(w, r)
		assert.Equal(t, tc.expectedCode, w.Code)
		assert.Equal(t, tc.expectedTasks, len(tasks))
		if len(tasks) > 3 {
			assert.Equal(t, tc.expectedID, tasks[3].ID)
			assert.Equal(t, tc.expectedInfo, tasks[3].Info)
			assert.Equal(t, tc.expectedDone, tasks[3].Done)
		}
		assert.Equal(t, tc.expectedNextTaskID, nextTaskID)
	}
}

func TestReadTasks(t *testing.T) {
	//curl -i -X GET http://localhost:8080/tasks
	bootstrapData()
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://localhost:8080/tasks", nil)
	readTasks(w, r)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "[{\"id\":1001,\"info\":\"Eat\",\"done\":true},{\"id\":1002,\"info\":\"Sleep\",\"done\":false},{\"id\":1003,\"info\":\"Rave\",\"done\":false}]\n", w.Body.String())
}

func TestReadTask(t *testing.T) {
	//curl -i -X GET http://localhost:8080/tasks/1001
	type test struct {
		id, expectedBody string
		expectedCode     int
	}
	tests := []test{
		{
			id:           "1001",
			expectedCode: 200,
			expectedBody: "{\"id\":1001,\"info\":\"Eat\",\"done\":true}\n",
		},
		{
			id:           "ERR",
			expectedCode: 404,
			expectedBody: "",
		},
		{
			id:           "999",
			expectedCode: 404,
			expectedBody: "",
		},
	}
	for _, tc := range tests {
		bootstrapData()
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:8080/tasks/%s", tc.id), nil)
		r = mux.SetURLVars(r, map[string]string{"id": tc.id})
		readTask(w, r)
		assert.Equal(t, tc.expectedCode, w.Code)
		assert.Equal(t, tc.expectedBody, w.Body.String())
	}
}

func TestUpdateTask(t *testing.T) {
	//curl -d '{"info":"Drink", "done":true}' -H "Content-Type: application/json" -i -X PUT http://localhost:8080/tasks/1002
	type test struct {
		id, body, expectedInfo string
		expectedCode           int
		expectedDone           bool
	}
	tests := []test{
		{
			id:           "1002",
			body:         "{\"info\":\"Drink\", \"done\":true}",
			expectedCode: 200,
			expectedInfo: "Drink",
			expectedDone: true,
		},
		{
			id:           "ERR",
			body:         "{\"info\":\"Drink\", \"done\":true}",
			expectedCode: 404,
			expectedInfo: "Sleep",
			expectedDone: false,
		},
		{
			id:           "999",
			body:         "{\"info\":\"Drink\", \"done\":true}",
			expectedCode: 404,
			expectedInfo: "Sleep",
			expectedDone: false,
		},
	}
	for _, tc := range tests {
		bootstrapData()
		w := httptest.NewRecorder()
		r := httptest.NewRequest("PUT", fmt.Sprintf("http://localhost:8080/tasks/%s", tc.id), strings.NewReader(tc.body))
		r = mux.SetURLVars(r, map[string]string{"id": tc.id})
		updateTask(w, r)
		assert.Equal(t, tc.expectedCode, w.Code)
		assert.Equal(t, 1002, tasks[1].ID)
		assert.Equal(t, tc.expectedInfo, tasks[1].Info)
		assert.Equal(t, tc.expectedDone, tasks[1].Done)
	}
}

func TestDeleteTask(t *testing.T) {
	//curl -i -X DELETE http://localhost:8080/tasks/1003
	type test struct {
		id                          string
		expectedCode, expectedTasks int
	}
	tests := []test{
		{
			id:            "1003",
			expectedCode:  200,
			expectedTasks: 2,
		},
		{
			id:            "ERR",
			expectedCode:  404,
			expectedTasks: 3,
		},
		{
			id:            "999",
			expectedCode:  404,
			expectedTasks: 3,
		},
	}
	for _, tc := range tests {
		bootstrapData()
		w := httptest.NewRecorder()
		r := httptest.NewRequest("DELETE", fmt.Sprintf("http://localhost:8080/tasks/%s", tc.id), nil)
		r = mux.SetURLVars(r, map[string]string{"id": tc.id})
		deleteTask(w, r)
		assert.Equal(t, tc.expectedCode, w.Code)
		assert.Equal(t, tc.expectedTasks, len(tasks))
	}
}
