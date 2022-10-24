package server

import (
	"encoding/json"
	"fmt"
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

	http.Error(w, "No such user", http.StatusNotFound)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "Invalid Json", http.StatusBadRequest)
		fmt.Println(err.Error())
		return
	}

	users = append(users, u)

	fmt.Println("Users updated: ")
	fmt.Println(users)
}

func HandleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/users/{name}", getUserAge).Methods(http.MethodGet)
	router.HandleFunc("/users", updateUser).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":8080", router))
}
