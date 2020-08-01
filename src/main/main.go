package main

import (
	"log"
	"net/http"
)

func main() {
	store := InMemoryUsersStore{
		users: map[string]string{},
	}
	server := &UsersServer{store: &store}

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
