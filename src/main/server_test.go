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
	RunTest(t, server, "unused url path: should return no string", "/someUnusedPath", http.StatusNotFound)
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

func TestFriendshipRequest(t *testing.T) {
	store := EmptyUsersStore()
	server := &UsersServer{store: store}

	// Sign up users
	RunSignUpTest(t, server, "sign up a new user", "arnau", "12345678", http.StatusOK)
	RunSignUpTest(t, server, "sign up a new user", "sergi", "12345678", http.StatusOK)
	RunSignUpTest(t, server, "sign up a new user", "berta", "12345678", http.StatusOK)

	// Send friendship request
	RunFriendshipRequestTest(t, server, "request friendship", "arnau", "sergi", "12345678", http.StatusOK)
	RunFriendshipRequestTest(t, server, "request friendship (to user does not exist)", "arnau", "barbara", "12345678", http.StatusBadRequest)
	RunFriendshipRequestTest(t, server, "request friendship (from user does not exist)", "david", "sergi", "12345678", http.StatusUnauthorized)
	RunFriendshipRequestTest(t, server, "request friendship (request is in pending status)", "arnau", "sergi", "12345678", http.StatusBadRequest)
	RunFriendshipRequestTest(t, server, "request friendship (opposite request has already been made)", "sergi", "arnau", "12345678", http.StatusBadRequest)
	// I decided not to accept such requests (they make no sense and complicate the accept request function)

	RunFriendshipRequestTest(t, server, "request friendship (from user has already sent one request)", "arnau", "berta", "12345678", http.StatusOK)
	RunFriendshipRequestTest(t, server, "request friendship (request is in pending status)", "arnau", "berta", "12345678", http.StatusBadRequest)
	RunFriendshipRequestTest(t, server, "request friendship (wrong password)", "berta", "sergi", "wrongPass", http.StatusUnauthorized)

	// Pending requests: arnau->sergi; sergi->arnau; arnau->berta

	// Accept friendship request
	RunRespondToFriendshipTest(t, server, "accept friendship (user does not exist)", "peter", "arnau", "12345678", true, http.StatusUnauthorized)
	RunRespondToFriendshipTest(t, server, "accept friendship (wrong password)", "sergi", "arnau", "wrongPass", true, http.StatusUnauthorized)
	RunRespondToFriendshipTest(t, server, "accept friendship (other user does not exist)", "sergi", "peter", "12345678", true, http.StatusBadRequest)
	RunRespondToFriendshipTest(t, server, "accept friendship (no existing request)", "sergi", "berta", "12345678", true, http.StatusBadRequest)
	RunRespondToFriendshipTest(t, server, "accept friendship (existing request in opposite direction)", "arnau", "berta", "12345678", true, http.StatusBadRequest)
	RunRespondToFriendshipTest(t, server, "accept friendship (OK)", "sergi", "arnau", "12345678", true, http.StatusOK)

	// Requests: arnau->berta
	// Friends: arnau&sergi

	RunListFriends(t, server, "list friends of arnau (sergi)", "arnau", "[sergi]", http.StatusOK)
	RunListFriends(t, server, "list friends of sergi (arnau)", "sergi", "[arnau]", http.StatusOK)
	RunListFriends(t, server, "list friends of berta (none)", "berta", "[]", http.StatusOK)
	RunListFriends(t, server, "list friends of peter (user does not exist)", "peter", "", http.StatusBadRequest)

	RunRespondToFriendshipTest(t, server, "accept friendship (already accepted)", "sergi", "arnau", "12345678", true, http.StatusBadRequest)
	RunRespondToFriendshipTest(t, server, "accept friendship (already accepted in opposite direction)", "arnau", "sergi", "12345678", true, http.StatusBadRequest)

	// Send a request to a person who's already friend
	RunFriendshipRequestTest(t, server, "request friendship (already friends)", "arnau", "sergi", "12345678", http.StatusBadRequest)

	// Decline friendship request
	RunRespondToFriendshipTest(t, server, "decline friendship (user does not exist)", "peter", "arnau", "12345678", false, http.StatusUnauthorized)
	RunRespondToFriendshipTest(t, server, "decline friendship (wrong password)", "berta", "arnau", "wrongPassword", false, http.StatusUnauthorized)
	RunRespondToFriendshipTest(t, server, "decline friendship (other user does not exist)", "berta", "peter", "12345678", false, http.StatusBadRequest)
	RunRespondToFriendshipTest(t, server, "decline friendship (no existing request)", "sergi", "berta", "12345678", false, http.StatusBadRequest)
	RunRespondToFriendshipTest(t, server, "decline friendship (request should have been removed)", "arnau", "sergi", "12345678", false, http.StatusBadRequest)
	RunRespondToFriendshipTest(t, server, "decline friendship (existing request in opposite direction)", "arnau", "berta", "12345678", false, http.StatusBadRequest)
	RunRespondToFriendshipTest(t, server, "decline friendship (OK)", "berta", "arnau", "12345678", false, http.StatusOK) // Not nice, Berta :(

	// Requests: -
	// Friends: arnau&sergi

	RunRespondToFriendshipTest(t, server, "decline friendship (already declined)", "berta", "arnau", "12345678", false, http.StatusBadRequest)
	RunFriendshipRequestTest(t, server, "request friendship (they declined)", "arnau", "berta", "12345678", http.StatusOK)
	// Requests: arnau->berta
	RunRespondToFriendshipTest(t, server, "decline friendship (again, but OK)", "berta", "arnau", "12345678", false, http.StatusOK) // :( :( :(
	// Requests: -
	RunFriendshipRequestTest(t, server, "request friendship (I declined)", "berta", "arnau", "12345678", http.StatusOK)
	// Requests: berta->arnau
	RunRespondToFriendshipTest(t, server, "decline friendship (again, but OK)", "arnau", "berta", "12345678", false, http.StatusOK) // sweet revenge
	// Requests: -

	RunListFriends(t, server, "list friends of arnau after friendship declines (sergi)", "arnau", "[sergi]", http.StatusOK)
	RunListFriends(t, server, "list friends of sergi after friendship declines (arnau)", "sergi", "[arnau]", http.StatusOK)
	RunListFriends(t, server, "list friends of berta after friendship declines (none)", "berta", "[]", http.StatusOK)

	// Note: this test case has gotten absurdly big.... I should better split it into several tests
}

func RunGetUsersTest(t *testing.T, s *UsersServer, name, want string) {
	request, _ := http.NewRequest(http.MethodGet, "/getUsers", nil)
	response := httptest.NewRecorder()

	s.ServeHTTP(response, request)

	t.Run(name, func(t *testing.T) {
		got := response.Body.String()
		want := want

		// ToDo check if users are sent in different order
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func RunTest(t *testing.T, s *UsersServer, name, url string, expectedHTTPStatus int) {
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	response := httptest.NewRecorder()

	s.ServeHTTP(response, request)

	t.Run(name, func(t *testing.T) {
		gotStatus := response.Code
		wantStatus := expectedHTTPStatus
		ok := AssertStatus(t, gotStatus, wantStatus)

		if !ok {
			gotBody := response.Body.String()
			t.Errorf("Got body: %q", gotBody)
		}
	})
}

func RunSignUpTest(t *testing.T, s *UsersServer, testName, username, password string, expectedHTTPStatus int) {
	requestBody, err := json.Marshal(map[string]string{
		"user": username,
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

func RunFriendshipRequestTest(t *testing.T, s *UsersServer, testName, userFrom, userTo, password string, expectedHTTPStatus int) {
	requestBody, err := json.Marshal(map[string]string{
		"user":   userFrom,
		"userTo": userTo,
		"pass":   password,
	})

	if err != nil {
		log.Fatalln(err)
	}

	request, err := http.NewRequest(http.MethodPost, "/requestFriendship", bytes.NewBuffer(requestBody))

	if err != nil {
		log.Fatalln(err)
	}

	request.Header.Set("Content-type", "appliaction/json")

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

func RunRespondToFriendshipTest(t *testing.T, s *UsersServer, testName, user, userFriendshipRequest, password string, acceptRequest bool, expectedHTTPStatus int) {
	accept := "0"
	if acceptRequest {
		accept = "1"
	}

	requestBody, err := json.Marshal(map[string]string{
		"user":                  user,
		"pass":                  password,
		"userFriendshipRequest": userFriendshipRequest,
		"acceptRequest":         accept,
	})

	if err != nil {
		log.Fatalln(err)
	}

	request, err := http.NewRequest(http.MethodPost, "/respondToFriendshipRequest", bytes.NewBuffer(requestBody))

	if err != nil {
		log.Fatalln(err)
	}

	request.Header.Set("Content-type", "appliaction/json")

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

func RunListFriends(t *testing.T, s *UsersServer, testName, user, expectedBody string, expectedHTTPStatus int) {
	request, _ := http.NewRequest(http.MethodGet, "/getFriends/"+user, nil)
	response := httptest.NewRecorder()

	s.ServeHTTP(response, request)

	t.Run(testName, func(t *testing.T) {
		gotStatus := response.Code
		wantStatus := expectedHTTPStatus
		gotBody := response.Body.String()
		wantBody := expectedBody

		if statusOk := AssertStatus(t, gotStatus, wantStatus); !statusOk {
			t.Errorf("Got body: %q", gotBody)
		}

		AssertResponseBody(t, gotBody, wantBody)
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
