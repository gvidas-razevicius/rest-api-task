package server

import (
	"encoding/json"
	"log"
	"net/http"

	mux "github.com/gorilla/mux"
)

var users = []User{
	{Name: "John", Age: 1},
	{Name: "Peter", Age: 2},
	{Name: "Mike", Age: 3},
}

func getUserAge(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	for _, user := range users {
		if user.Name == name {
			json.NewEncoder(w).Encode(user)
			return
		}
	}

	w.WriteHeader(http.StatusBadRequest)
}

func HandleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/users/{name}", getUserAge).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8080", router))
}
