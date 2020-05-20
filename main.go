package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe(":8080", router))
}
