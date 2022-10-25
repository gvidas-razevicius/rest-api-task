package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var users map[string]User

func getUserAge(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	res, found := users[name]

	if found {
		json.NewEncoder(w).Encode(res)
		return
	}
	http.Error(w, "No such user", http.StatusNotFound)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "Invalid Json", http.StatusBadRequest)
		return
	}

	if _, found := users[u.Name]; !found {
		users[u.Name] = u
		fmt.Println("User created: ")
		fmt.Println(users)
	} else {
		http.Error(w, "User already exists!", http.StatusForbidden)
	}
}

func HandleRequests() {
	users = make(map[string]User)
	users["John"] = User{Name: "John", Age: 1}
	users["Peter"] = User{Name: "Peter", Age: 2}
	users["Mike"] = User{Name: "Mike", Age: 3}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/users/{name}", getUserAge).Methods(http.MethodGet)
	router.HandleFunc("/users", createUser).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":8080", router))
}
