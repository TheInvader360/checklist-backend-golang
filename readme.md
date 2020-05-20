### Minimal REST API

Run:

    cd ~/go/src/github.com/TheInvader360/checklist-backend-golang
    go build
    ./checklist-backend-golang

Usage (CRUD):

    curl -d '{"id":1004, "info":"Learn", "done":false}' -H "Content-Type: application/json" -i -X POST http://localhost:8080/task
    curl -i -X GET http://localhost:8080/tasks
    curl -d '{"info":"Drink", "done":true}' -H "Content-Type: application/json" -i -X PUT http://localhost:8080/tasks/1002
    curl -i -X DELETE http://localhost:8080/tasks/1003
