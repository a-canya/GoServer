package main

import (
	"fmt"
	"net/http"
)

func GetUsersServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "[arnau]")
}
