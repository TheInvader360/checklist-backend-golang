# Minimal REST API

Run:

    git clone https://github.com/TheInvader360/checklist-backend-golang
    cd checklist-backend-golang
    go run main.go

Usage (CRUD):

    curl -d '{"info":"Learn", "done":false}' -H "Content-Type: application/json" -i -X POST http://localhost:8080/task
    curl -i -X GET http://localhost:8080/tasks
    curl -i -X GET http://localhost:8080/tasks/1001
    curl -d '{"info":"Drink", "done":true}' -H "Content-Type: application/json" -i -X PUT http://localhost:8080/tasks/1002
    curl -i -X DELETE http://localhost:8080/tasks/1003
