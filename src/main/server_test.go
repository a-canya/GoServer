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

func RunTest(t *testing.T, op *testOptions) {
	request, _ := http.NewRequest(http.MethodGet, op.url, nil)
	response := httptest.NewRecorder()

	Server(response, request)

	t.Run(op.name, func(t *testing.T) {
		got := response.Body.String()
		want := op.want

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestGetUsers(t *testing.T) {
	// Request users
	RunTest(t, &testOptions{name: "returns list of users in the social network", want: "[arnau]", url: "/getUsers"})

	// Wrong request
	RunTest(t, &testOptions{name: "unused url path: should return no string", want: "", url: "/someUnusedPath"})
}

func TestSignUp(t *testing.T) {
	// Request users
	RunTest(t, &testOptions{name: "list of users at the beginning should be empty", want: "[]", url: "/getUsers"})

	// Sign up new user
	requestBody, err := json.Marshal(map[string]string{
		"name": "arnau",
		"pass": "1234",
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

	Server(response, request)

	t.Run("sign up a new user", func(t *testing.T) {
		got := response.Body.String()
		want := "ok"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	// Request users
	RunTest(t, &testOptions{name: "list of users should include recently created user", want: "[arnau]", url: "/getUsers"})
}