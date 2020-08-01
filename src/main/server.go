package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
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

		if ok := CheckUsernameAndPassword(info["name"], info["pass"], &w); !ok {
			return
		}

		if ok := s.store.AddUser(info["name"], info["pass"]); ok {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "User already exists")
		}
	}
}

// CheckUsernameAndPassword returns true iff username has 5-10 alphanum characters and password has 8-12 alphanum chars.
// If conditions are not fulfilled, prints an error message in w
func CheckUsernameAndPassword(username, password string, w *http.ResponseWriter) bool {
	ok := true

	var isStringAlphabetic = regexp.MustCompile(`^[a-zA-Z0-9_]*$`).MatchString

	if !isStringAlphabetic(username) {
		(*w).WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(*w, "Username has invalid characters! Username must be unique, from 5 to 10 alphanumeric characters.")
		ok = false
	}

	if !isStringAlphabetic(password) {
		(*w).WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(*w, "Password has invalid characters! Password must have from 8 to 12 alphanumeric characters.")
		ok = false
	}

	// Note: checking len(var) returns the length in bytes but this might not correspond to the number of characters because
	// standard allows characters of multiple bytes. However, since we only accept alphanumeric characters this is ok
	if len(username) < 5 {
		(*w).WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(*w, "Username too short! Username must be unique, from 5 to 10 alphanumeric characters.")
		ok = false
	} else if len(username) > 10 {
		(*w).WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(*w, "Username too long! Username must be unique, from 5 to 10 alphanumeric characters.")
		ok = false
	}

	if len(password) < 8 {
		(*w).WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(*w, "Password too short! Password must have from 8 to 12 alphanumeric characters.")
		ok = false
	} else if len(password) > 12 {
		(*w).WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(*w, "Password too long! Password must have from 8 to 12 alphanumeric characters.")
		ok = false
	}

	return ok
}
