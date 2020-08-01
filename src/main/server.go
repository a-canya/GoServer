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
	UserExists(name string) bool
	RequestFriendship(from, to string) bool
	CheckUsersPassword(user, password string) bool
	RespondToFriendshipRequest(user, otherUser string, acceptRequest bool) bool
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
		s.SignUp(&w, r)

	case "requestFriendship":
		s.RequestFriendship(&w, r)

	case "respondToFriendshipRequest":
		s.RespondToFriendshipRequest(&w, r)

	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

// SignUp takes a signUp HTTP request (r) to the UsersServer (s), processes it and populates the ResponseWriter (w)
func (s *UsersServer) SignUp(w *http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(*w, "Couldn't read the data")
	}

	var info map[string]string
	json.Unmarshal(body, &info)

	if ok, msg := CheckUsernameAndPassword(info["user"], info["pass"]); !ok {
		(*w).WriteHeader(http.StatusBadRequest)
		fmt.Fprint(*w, msg)
		return
	}

	if ok := s.store.AddUser(info["user"], info["pass"]); ok {
		(*w).WriteHeader(http.StatusOK)
	} else {
		(*w).WriteHeader(http.StatusBadRequest)
		fmt.Fprint(*w, "User already exists")
	}
}

// CheckUsernameAndPassword returns true iff username has 5-10 alphanum characters and password has 8-12 alphanum chars.
// If conditions are not fulfilled, msg holds an error message
func CheckUsernameAndPassword(username, password string) (bool, string) {
	ok := true
	msg := ""

	var isStringAlphabetic = regexp.MustCompile(`^[a-zA-Z0-9_]*$`).MatchString

	if !isStringAlphabetic(username) {
		msg += "Username has invalid characters! Username must be unique, from 5 to 10 alphanumeric characters."
		ok = false
	}

	if !isStringAlphabetic(password) {
		msg += "Password has invalid characters! Password must have from 8 to 12 alphanumeric characters."
		ok = false
	}

	// Note: checking len(var) returns the length in bytes but this might not correspond to the number of characters because
	// standard allows characters of multiple bytes. However, since we only accept alphanumeric characters this is ok
	if len(username) < 5 {
		msg += "Username too short! Username must be unique, from 5 to 10 alphanumeric characters."
		ok = false
	} else if len(username) > 10 {
		msg += "Username too long! Username must be unique, from 5 to 10 alphanumeric characters."
		ok = false
	}

	if len(password) < 8 {
		msg += "Password too short! Password must have from 8 to 12 alphanumeric characters."
		ok = false
	} else if len(password) > 12 {
		msg += "Password too long! Password must have from 8 to 12 alphanumeric characters."
		ok = false
	}

	return ok, msg
}

// RequestFriendship takes a requestFriendship HTTP request (r) to the UsersServer (s), processes it and populates the ResponseWriter (w)
func (s *UsersServer) RequestFriendship(w *http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(*w, "Couldn't read the data")
	}

	var info map[string]string
	json.Unmarshal(body, &info)

	// Check credentials
	if !s.store.CheckUsersPassword(info["user"], info["pass"]) {
		(*w).WriteHeader(http.StatusUnauthorized)
		return
	}

	// Check if users exist
	if !s.store.UserExists(info["user"]) || !s.store.UserExists(info["userTo"]) {
		(*w).WriteHeader(http.StatusBadRequest)
		fmt.Fprint(*w, "User does not exist") // error msg could be more explicit: which user does not exist?
		return
	}

	// Add request to the DB
	if ok := s.store.RequestFriendship(info["user"], info["userTo"]); ok {
		(*w).WriteHeader(http.StatusOK)
	} else {
		(*w).WriteHeader(http.StatusBadRequest)
		fmt.Fprint(*w, "Friendship request already exists")
	}
}

// RespondToFriendshipRequest takes a respondToFriendshipRequest HTTP request (r) to the UsersServer (s),
// processes it and populates the ResponseWriter (w)
func (s *UsersServer) RespondToFriendshipRequest(w *http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(*w, "Couldn't read the data")
	}

	var info map[string]string
	json.Unmarshal(body, &info)

	user := info["user"]
	pass := info["pass"]
	otherUser := info["userFriendshipRequest"]
	accept := false
	if info["acceptRequest"] == "1" {
		accept = true
	} else if info["acceptRequest"] != "0" {
		(*w).WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(*w, "acceptRequest field must be either 1 or 0")
	}

	// Check credentials
	if !s.store.CheckUsersPassword(user, pass) {
		(*w).WriteHeader(http.StatusUnauthorized)
		return
	}

	// Check if users exist
	if !s.store.UserExists(user) || !s.store.UserExists(otherUser) {
		(*w).WriteHeader(http.StatusBadRequest)
		fmt.Fprint(*w, "User does not exist") // error msg could be more explicit: which user does not exist?
		return
	}

	// Respond to friendship request
	if ok := s.store.RespondToFriendshipRequest(user, otherUser, accept); ok {
		(*w).WriteHeader(http.StatusOK)
	} else {
		(*w).WriteHeader(http.StatusBadRequest)
		fmt.Fprint(*w, "Cannot respond to friendship request because request does not exist")
	}
}
