package main

// InMemoryUsersStore collects data about users in memory.
type InMemoryUsersStore struct {
	users map[string]string
}

// GetUsers retrieves list of users
func (i *InMemoryUsersStore) GetUsers() []string {
	usernames := GetKeys(&i.users)
	return usernames
}

// AddUser adds a user
// STUB
func (i *InMemoryUsersStore) AddUser(name string, password string) bool {
	if _, alreadyExists := i.users[name]; alreadyExists {
		return false
	}
	i.users[name] = password
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

// EmptyUsersStore returns a new empty InMemoryUsersStore
func EmptyUsersStore() *InMemoryUsersStore {
	store := InMemoryUsersStore{
		users: map[string]string{},
	}
	return &store
}
