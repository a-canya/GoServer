package main

// InMemoryUsersStore collects data about users in memory.
type InMemoryUsersStore struct {
	users              map[string]string
	friendshipRequests map[string][]string // friendshipRequests["john0"] == {"peter", "mike5"} means john0 has sent a friendship request to peter and mike5
	friends            map[string][]string // must be kept symmetric all time, ie Contains(friends["peter"], "mike5") <==> Contains(friends["mike5"], "peter")
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
	i.friendshipRequests[name] = make([]string, 0)
	i.friends[name] = make([]string, 0)
	return true
}

// UserExists returns true iff user with name `name` exists
func (i *InMemoryUsersStore) UserExists(name string) bool {
	_, exists := i.users[name]
	return exists
}

// RequestFriendship adds a friendship request from user `from` to user `to`.
// Returns false iff friendship request between both users already exists or users are already friends (in this case no modifications are made)
func (i *InMemoryUsersStore) RequestFriendship(from, to string) bool {
	// log.Println("RequestFriendship from", from, "to", to)
	myRequests := i.friendshipRequests[from]
	// log.Println("Current requests from", from, ":", myRequests)
	if Contains(myRequests, to) {
		return false
	}

	theirRequests := i.friendshipRequests[to]
	if Contains(theirRequests, from) {
		return false
	}

	myFriends := i.friends[from]
	if Contains(myFriends, to) {
		return false
	}

	i.friendshipRequests[from] = append(myRequests, to)
	// log.Println("Done! Requests =", requests, " and i.friendshipRequests[from] =", i.friendshipRequests[from])

	return true
}

// CheckUsersPassword returns true if user existst and has this password
func (i *InMemoryUsersStore) CheckUsersPassword(user, password string) bool {
	storedPassword, exists := i.users[user]
	return exists && storedPassword == password
}

// RespondToFriendshipRequest responds to a friendship request from otherUser made to user
// Returns false iff friendship request does not exist (in this case no modifications are made)
func (i *InMemoryUsersStore) RespondToFriendshipRequest(user, otherUser string, acceptRequest bool) bool {
	requests, hasRequests := i.friendshipRequests[otherUser]
	if !hasRequests {
		return false
	}

	for j := 0; j < len(requests); j++ {
		if requests[j] == user {
			i.friendshipRequests[otherUser] = Remove(i.friendshipRequests[otherUser], j)
			if acceptRequest {
				i.friends[user] = append(i.friends[user], otherUser)
				i.friends[otherUser] = append(i.friends[otherUser], user)
			}
			return true
		}
	}

	return false
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

// Remove returns a slice identical to s except that element at position i is eliminated (order not preserved)
func Remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// EmptyUsersStore returns a new empty InMemoryUsersStore
func EmptyUsersStore() *InMemoryUsersStore {
	store := InMemoryUsersStore{
		users:              map[string]string{},
		friendshipRequests: map[string][]string{},
		friends:            map[string][]string{},
	}
	return &store
}
