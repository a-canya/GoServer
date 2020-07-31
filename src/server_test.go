package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUsers(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	GetUsersServer(response, request)

	t.Run("returns list fo users in the social network", func(t *testing.T) {
		got := response.Body.String()
		want := "[arnau]"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
