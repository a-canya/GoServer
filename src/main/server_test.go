package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
	RunSignUpTest(t, server, "sign up a new user", "arnau", "12345678", http.StatusOK)

	// Request users
	RunGetUsersTest(t, server, "list of users should include recently created user", "[arnau]")

	// Sign up another new user
	RunSignUpTest(t, server, "sign up another new user", "carla", "Password", http.StatusOK)

	// Sign up already existing user
	RunSignUpTest(t, server, "sign up an already existing user", "arnau", "12345678", http.StatusBadRequest)

	// Sign ups not fulfilling username/password constraints (username 5-10 alphanum characters, password 8-12 alphanum chars)
	RunSignUpTest(t, server, "sign up with username too short", "john", "password", http.StatusBadRequest)
	RunSignUpTest(t, server, "sign up with username too long", "montserrat0", "password", http.StatusBadRequest)
	RunSignUpTest(t, server, "sign up with password too short", "jonathan", "1234567", http.StatusBadRequest)
	RunSignUpTest(t, server, "sign up with password too long", "william", "123456789abcd", http.StatusBadRequest)
	RunSignUpTest(t, server, "sign up with non-alphanumeric username", "arnau!", "12345678", http.StatusBadRequest)
	RunSignUpTest(t, server, "sign up with non-alphanumeric password", "maria", "$€cr€tWörd", http.StatusBadRequest)
	RunSignUpTest(t, server, "sign up with lots of constraints not fulfilled", "I'm the boss!", "¬¬'", http.StatusBadRequest)

	// Request users
	RunGetUsersTest(t, server, "list of users should include both users created", "[arnau carla]")
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
		AssertResponseBody(t, got, want)
	})
}

func RunSignUpTest(t *testing.T, s *UsersServer, testName, username, password string, expectedHTTPStatus int) {
	requestBody, err := json.Marshal(map[string]string{
		"name": username,
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

	t.Run(testName, func(t *testing.T) {
		gotStatus := response.Code
		wantStatus := expectedHTTPStatus
		ok := AssertStatus(t, gotStatus, wantStatus)

		if !ok {
			gotBody := response.Body.String()
			t.Errorf("Got body: %q", gotBody)
		}
	})
}

func AssertStatus(t *testing.T, got, want int) bool {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
		return false
	}
	return true
}

func AssertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}
