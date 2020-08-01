package main

import (
	"fmt"
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
		user := "arnau"
		password := "1234"
		if ok := s.store.AddUser(user, password); ok {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "User already exists")
		}
	}
}
