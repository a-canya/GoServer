package main

import (
	"fmt"
	"net/http"
	"strings"
)

// Server serves HTTP requests
func Server(w http.ResponseWriter, r *http.Request) {
	option := strings.Split(r.URL.Path, "/")[1]
	switch option {
	case "getUsers":
		fmt.Fprint(w, GetUsers())
	}
	//fmt.Fprint(w, option)
}

// GetUsers returns a list fo all username sin the social network
// STUB
func GetUsers() []string {
	users := []string{"arnau"}
	return users
}
