package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func getUser(w http.ResponseWriter, r *http.Request) {
	names := r.URL.Query()["names"]

	var res UserArray
	res.Array = make([]User, len(names))
	for i, name := range names {
		user, found := cache.Users[name]

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

	if _, found := cache.Users[obj.Name]; !found {
		cache.Users[obj.Name] = obj
		fmt.Println("User created: ", cache.Users)
		w.WriteHeader(http.StatusCreated)
		go writeToDisk()
	} else {
		http.Error(w, "User already exists!", http.StatusForbidden)
	}
}

func delUser(w http.ResponseWriter, r *http.Request) {
	names := r.URL.Query()["names"]
	if len(names) > 1 {
		http.Error(w, "Cannot process multiple objects", http.StatusBadRequest)
		return
	}

	name := names[0]
	if _, found := cache.Users[name]; found {
		delete(cache.Users, name)
		fmt.Println("User deleted: ", cache.Users)
		w.WriteHeader(http.StatusNoContent)
		go writeToDisk()
	} else {
		http.Error(w, "User not found!", http.StatusNotFound)
	}
}

func getApp(w http.ResponseWriter, r *http.Request) {
	names := r.URL.Query()["names"]
	if len(names) > 1 {
		http.Error(w, "Cannot process multiple objects", http.StatusBadRequest)
		return
	}

	name := names[0]

	if res, found := cache.Apps[name]; found {
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

	if _, found := cache.Apps[obj.Name]; !found {
		obj.Created = StringInt(time.Now().Year())
		cache.Apps[obj.Name] = obj
		fmt.Println("App created: ", cache.Apps)
		w.WriteHeader(http.StatusCreated)
		go writeToDisk()
	} else {
		http.Error(w, "App already exists!", http.StatusForbidden)
	}
}

func delApp(w http.ResponseWriter, r *http.Request) {
	names := r.URL.Query()["names"]
	if len(names) > 1 {
		http.Error(w, "Cannot process multiple objects", http.StatusBadRequest)
		return
	}

	name := names[0]
	if _, found := cache.Apps[name]; found {
		delete(cache.Apps, name)
		fmt.Println("App deleted: ", cache.Apps)
		w.WriteHeader(http.StatusNoContent)
		go writeToDisk()
	} else {
		http.Error(w, "App not found!", http.StatusNotFound)
	}
}

func HandleRequests() {
	if err := initDb(); err != nil {
		log.Fatalf("Could not initiliase the database: %v", err)
	}
	fmt.Println(cache)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/users", getUser).Methods(http.MethodGet)
	router.HandleFunc("/users", delUser).Methods(http.MethodDelete)
	router.HandleFunc("/users", createUser).Methods(http.MethodPost)
	router.HandleFunc("/apps", getApp).Methods(http.MethodGet)
	router.HandleFunc("/apps", delApp).Methods(http.MethodDelete)
	router.HandleFunc("/apps", createApp).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":8080", router))
}
