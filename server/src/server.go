package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func getUser(w http.ResponseWriter, r *http.Request) {
	get(w, r, cache.Users)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	create(w, r, cache.Users)
}

func delUser(w http.ResponseWriter, r *http.Request) {
	del(w, r, cache.Users)
}

func getApp(w http.ResponseWriter, r *http.Request) {
	get(w, r, cache.Apps)
}

func createApp(w http.ResponseWriter, r *http.Request) {
	create(w, r, cache.Apps)
}

func delApp(w http.ResponseWriter, r *http.Request) {
	del(w, r, cache.Apps)
}

func get[T Object](w http.ResponseWriter, r *http.Request, mem map[string]T) {
	names := r.URL.Query()["names"]

	res := make([]T, len(names))
	for i, name := range names {
		obj, found := mem[name]

		if !found {
			http.Error(w, "One or more object was not found", http.StatusNotFound)
			return
		}

		res[i] = obj
	}
	json.NewEncoder(w).Encode(res)
}

func create[T Object](w http.ResponseWriter, r *http.Request, mem map[string]T) {
	var obj T
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		http.Error(w, "Invalid Json", http.StatusBadRequest)
		fmt.Printf("Error creating app: %v", err)
		return
	}

	name := obj.GetName()
	if _, found := cache.Apps[name]; !found {
		mem[obj.GetName()] = obj
		fmt.Printf("%T created: %v", mem[name], mem)
		w.WriteHeader(http.StatusCreated)
		go writeToDisk()
	} else {
		http.Error(w, fmt.Sprintf("%s not found!", mem[name].GetType()), http.StatusForbidden)
	}
}

func del[T Object](w http.ResponseWriter, r *http.Request, mem map[string]T) {
	names := r.URL.Query()["names"]
	if len(names) > 1 {
		http.Error(w, "Cannot process multiple objects", http.StatusBadRequest)
		return
	}
	if len(names) == 0 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	name := names[0]
	if _, found := mem[name]; found {
		delete(mem, name)
		fmt.Printf("%T deleted: %v", mem[name], mem)
		w.WriteHeader(http.StatusNoContent)
		go writeToDisk()
	} else {
		http.Error(w, fmt.Sprintf("%s not found!", mem[name].GetType()), http.StatusNotFound)
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