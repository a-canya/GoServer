package main

import "log"

// InMemoryUsersStore collects data about users in memory.
type InMemoryUsersStore struct {
	users              map[string]string
	friendshipRequests map[string][]string // friendshipRequests["john0"] == {"peter", "mike5"} means john0 has sent a friendship request to peter and mike5
}

// GetUsers retrieves list of users
func (i *InMemoryUsersStore) GetUsers() []string {
	usernames := GetKeys(&i.users)
	return usernames
}

// AddUser adds a user with given username and password.
// Returns false iff username already exists (in this case no modifications are made)
func (i *InMemoryUsersStore) AddUser(name string, password string) bool {
	if _, alreadyExists := i.users[name]; alreadyExists {
		return false
	}
	i.users[name] = password
	return true
}

// UserExists returns true iff user with name `name` exists
func (i *InMemoryUsersStore) UserExists(name string) bool {
	_, exists := i.users[name]
	return exists
}

// RequestFriendship adds a friendship request from user `from` to user `to`.
// Returns false iff friendship request already exists (in this case no modifications are made)
func (i *InMemoryUsersStore) RequestFriendship(from, to string) bool {
	log.Println("RequestFriendship from", from, "to", to)
	requests, exists := i.friendshipRequests[from]
	log.Println("Current requests from", from, ":", requests)
	if exists {
		if Contains(requests, to) {
			return false
		}
		i.friendshipRequests[from] = append(requests, to)
	} else {
		i.friendshipRequests[from] = []string{to}
	}
	log.Println("Done! Requests =", requests, " and i.friendshipRequests[from] =", i.friendshipRequests[from])
	return true
}

// GetKeys returns a slice of the keys of map m
// thoughts: returning a ptr might be more efficient; implementing this with interfaces would make func more general
func GetKeys(m *map[string]string) []string {
	keys := make([]string, len(*m))

	i := 0
	for k := range *m {
		keys[i] = k
		i++
	}

	return keys
}

// Contains returns true iff slice s contains element e
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// EmptyUsersStore returns a new empty InMemoryUsersStore
func EmptyUsersStore() *InMemoryUsersStore {
	store := InMemoryUsersStore{
		users:              map[string]string{},
		friendshipRequests: map[string][]string{},
	}
	return &store
}
