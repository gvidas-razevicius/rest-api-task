package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var users map[string]User
var apps map[string]App

func getUser(w http.ResponseWriter, r *http.Request) {
	var names NamesArray
	err := names.DecodePayload(r.Body, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var res UserArray
	res.Array = make([]User, len(names.Array))
	for i, name := range names.Array {
		user, found := users[name]

		if !found {
			http.Error(w, "One or more user was not found", http.StatusNotFound)
			return
		}

		res.Array[i] = user
	}
	json.NewEncoder(w).Encode(res)
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
		fmt.Println("User created: ", users)
		w.WriteHeader(http.StatusCreated)
	} else {
		http.Error(w, "User already exists!", http.StatusForbidden)
	}
}

func delUser(w http.ResponseWriter, r *http.Request) {
	var names NamesArray
	err := names.DecodePayload(r.Body, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := names.Array[0]
	if _, found := users[name]; found {
		delete(users, name)
		fmt.Println("User deleted: ", users)
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Error(w, "User not found!", http.StatusNotFound)
	}
}

func getApp(w http.ResponseWriter, r *http.Request) {
	var names NamesArray
	err := names.DecodePayload(r.Body, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := names.Array[0]

	if res, found := apps[name]; found {
		json.NewEncoder(w).Encode(res)
		return
	}
	http.Error(w, "No such app", http.StatusNotFound)
}

func createApp(w http.ResponseWriter, r *http.Request) {
	var obj App
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		http.Error(w, "Invalid Json", http.StatusBadRequest)
		fmt.Printf("Error creating app: %v", err)
		return
	}

	if _, found := apps[obj.Name]; !found {
		obj.Created = StringInt(time.Now().Year())
		apps[obj.Name] = obj
		fmt.Println("App created: ", apps)
		w.WriteHeader(http.StatusCreated)
	} else {
		http.Error(w, "App already exists!", http.StatusForbidden)
	}
}

func delApp(w http.ResponseWriter, r *http.Request) {
	var names NamesArray
	err := names.DecodePayload(r.Body, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := names.Array[0]
	if _, found := apps[name]; found {
		delete(apps, name)
		fmt.Println("App deleted: ", apps)
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Error(w, "App not found!", http.StatusNotFound)
	}
}

// Decodes the payload into input. Returns an error message to be pased into the body of response
// if anything goes wrong.
// -payload: the body of the request
// -enfSingle: enforce single value. If true returns error if payload has more than one object
func (input *NamesArray) DecodePayload(payload io.Reader, enfSingle bool) error {
	if err := json.NewDecoder(payload).Decode(&input); err != nil || len(input.Array) == 0 {
		return fmt.Errorf("Bad payload")
	}

	if len(input.Array) > 1 && enfSingle {
		return fmt.Errorf("Cannot process multiple objects")
	}

	return nil
}

func HandleRequests() {
	users = make(map[string]User)
	apps = make(map[string]App)
	users["John"] = User{Name: "John", Age: 1}
	users["Peter"] = User{Name: "Peter", Age: 2}
	users["Mike"] = User{Name: "Mike", Age: 3}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/users", getUser).Methods(http.MethodGet)
	router.HandleFunc("/users", delUser).Methods(http.MethodDelete)
	router.HandleFunc("/users", createUser).Methods(http.MethodPost)
	router.HandleFunc("/app", getApp).Methods(http.MethodGet)
	router.HandleFunc("/app", delApp).Methods(http.MethodDelete)
	router.HandleFunc("/app", createApp).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":8080", router))
}