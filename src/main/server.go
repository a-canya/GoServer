package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// UsersStore is an interface for a DB in which we can add and retrieve users
type UsersStore interface {
	GetUsers() []string
	AddUser(name string, password string) bool
}

// UsersServer is a strcuture which contains an interface to interact with the users DB
type UsersServer struct {
	store UsersStore
}

// ServeHTTP serves HTTP requests
func (s *UsersServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	option := strings.Split(r.URL.Path, "/")[1]
	switch option {
	case "getUsers":
		fmt.Fprint(w, s.store.GetUsers())
	case "signUp":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Couldn't read the data")
		}

		var info map[string]string
		json.Unmarshal(body, &info)

		if ok := s.store.AddUser(info["name"], info["pass"]); ok {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "User already exists")
		}
	}
}
