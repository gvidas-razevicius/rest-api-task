package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var users map[string]User
var apps map[string]App

func getUser(w http.ResponseWriter, r *http.Request) {
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
	var obj User
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		http.Error(w, "Invalid Json", http.StatusBadRequest)
		return
	}

	if _, found := users[obj.Name]; !found {
		users[obj.Name] = obj
		fmt.Println("User created: ")
		fmt.Println(users)
		w.WriteHeader(http.StatusCreated)
	} else {
		http.Error(w, "User already exists!", http.StatusForbidden)
	}
}

func delUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	if _, found := users[name]; found {
		delete(users, name)
		fmt.Println("User deleted: ")
		fmt.Println(users)
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Error(w, "User not found!", http.StatusNotFound)
	}
}

func getApp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	res, found := apps[name]

	if found {
		json.NewEncoder(w).Encode(res)
		return
	}
	http.Error(w, "No such app", http.StatusNotFound)
}

func createApp(w http.ResponseWriter, r *http.Request) {
	var obj App
	// TODO: add year created setting
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		http.Error(w, "Invalid Json", http.StatusBadRequest)
		fmt.Printf("Error creating app: %v", err)
		return
	}

	if _, found := apps[obj.Name]; !found {
		apps[obj.Name] = obj
		fmt.Println("App created: ")
		fmt.Println(apps)
		w.WriteHeader(http.StatusCreated)
	} else {
		http.Error(w, "App already exists!", http.StatusForbidden)
	}
}

func delApp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	if _, found := apps[name]; found {
		delete(apps, name)
		fmt.Println("App deleted: ")
		fmt.Println(apps)
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Error(w, "App not found!", http.StatusNotFound)
	}
}

func HandleRequests() {
	users = make(map[string]User)
	apps = make(map[string]App)
	users["John"] = User{Name: "John", Age: 1}
	users["Peter"] = User{Name: "Peter", Age: 2}
	users["Mike"] = User{Name: "Mike", Age: 3}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/users/{name}", getUser).Methods(http.MethodGet)
	router.HandleFunc("/users/{name}", delUser).Methods(http.MethodDelete)
	router.HandleFunc("/users", createUser).Methods(http.MethodPost)
	router.HandleFunc("/app/{name}", getApp).Methods(http.MethodGet)
	router.HandleFunc("/app/{name}", delApp).Methods(http.MethodDelete)
	router.HandleFunc("/app", createApp).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":8080", router))
}
