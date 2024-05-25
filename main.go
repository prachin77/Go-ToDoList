package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prachin77/ToDoList/handlers"
)

func main() {
	// handlers.Init()

	r := mux.NewRouter()

	fmt.Println("listening to port 8080 .....")
	r.HandleFunc("/", handlers.DefaultRoute).Methods("GET")
	r.HandleFunc("/app", handlers.GetApp).Methods("GET")

	r.HandleFunc("/addtask", handlers.AddTask).Methods("POST")
	r.HandleFunc("/delete/{id}", handlers.DeleteTask).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}
