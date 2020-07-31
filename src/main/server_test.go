package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUsers(t *testing.T) {
	// Request users
	request, _ := http.NewRequest(http.MethodGet, "/getUsers", nil)
	response := httptest.NewRecorder()

	Server(response, request)

	t.Run("returns list fo users in the social network", func(t *testing.T) {
		got := response.Body.String()
		want := "[arnau]"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	// Wrong request
	request2, _ := http.NewRequest(http.MethodGet, "/someUnusedPath", nil)
	response2 := httptest.NewRecorder()

	Server(response2, request2)

	t.Run("unused url path: returns error", func(t *testing.T) {
		got := response2.Body.String()
		want := ""

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

}
