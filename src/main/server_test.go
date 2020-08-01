package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testOptions struct {
	name string
	want string
	url  string
}

func RunGetUsersTest(t *testing.T, s *UsersServer, name, want string) {
	request, _ := http.NewRequest(http.MethodGet, "/getUsers", nil)
	response := httptest.NewRecorder()

	s.ServeHTTP(response, request)

	t.Run(name, func(t *testing.T) {
		got := response.Body.String()
		want := want

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func RunTest(t *testing.T, s *UsersServer, name, url, want string) {
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	response := httptest.NewRecorder()

	s.ServeHTTP(response, request)

	t.Run(name, func(t *testing.T) {
		got := response.Body.String()
		want := want

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func RunSignUpTest(t *testing.T, s *UsersServer, name, password, want string) {
	requestBody, err := json.Marshal(map[string]string{
		"name": name,
		"pass": password,
	})

	if err != nil {
		log.Fatalln(err)
	}

	request, err := http.NewRequest(http.MethodPost, "/signUp", bytes.NewBuffer(requestBody))

	if err != nil {
		log.Fatalln(err)
	}

	request.Header.Set("Content-type", "application/json")

	response := httptest.NewRecorder()

	s.ServeHTTP(response, request)

	t.Run("sign up a new user", func(t *testing.T) {
		got := response.Body.String()
		want := want

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestGetUsers(t *testing.T) {
	store := InMemoryUsersStore{
		users: map[string]string{
			"arnau": "1234",
		},
	}

	server := &UsersServer{store: &store}

	// Request users
	RunGetUsersTest(t, server, "returns list of users in the social network", "[arnau]")

	// Wrong request
	RunTest(t, server, "unused url path: should return no string", "/someUnusedPath", "")
}

func TestSignUp(t *testing.T) {
	store := EmptyUsersStore()

	server := &UsersServer{store: store}

	// Request users
	RunGetUsersTest(t, server, "list of users at the beginning should be empty", "[]")

	// Sign up new user
	RunSignUpTest(t, server, "arnau", "1234", "ok")

	// Request users
	RunGetUsersTest(t, server, "list of users should include recently created user", "[arnau]")
}
