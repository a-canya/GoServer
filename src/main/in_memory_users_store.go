package main

// InMemoryUsersStore collects data about users in memory.
type InMemoryUsersStore struct {
	users              map[string]string
	friendshipRequests map[string][]string // friendshipRequests["john0"] == {"peter", "mike5"} means john0 has sent a friendship request to peter and mike5
	friends            map[string][]string // must be kept symmetric all time, ie Contains(friends["peter"], "mike5") <==> Contains(friends["mike5"], "peter")
}

// GetUsers retrieves a list of all users
func (s *InMemoryUsersStore) GetUsers() []string {
	usernames := GetKeys(&s.users)
	return usernames
}

// AddUser adds a user with given username and password.
// Returns false iff username already exists (in this case no modifications are made)
func (s *InMemoryUsersStore) AddUser(name string, password string) bool {
	if _, alreadyExists := s.users[name]; alreadyExists {
		return false
	}
	s.users[name] = password
	s.friendshipRequests[name] = make([]string, 0)
	s.friends[name] = make([]string, 0)
	return true
}

// UserExists returns true iff user with name `name` exists
func (s *InMemoryUsersStore) UserExists(name string) bool {
	_, exists := s.users[name]
	return exists
}

// RequestFriendship adds a friendship request from user `from` to user `to`.
// Returns false iff friendship request between both users already exists or users are already friends (in this case no modifications are made)
// Precondition: from and to users exist in the DB and have been correctly initialized (ie using AddUser function)
func (s *InMemoryUsersStore) RequestFriendship(from, to string) bool {
	myRequests := s.friendshipRequests[from]
	if Contains(myRequests, to) {
		return false
	}

	theirRequests := s.friendshipRequests[to]
	if Contains(theirRequests, from) {
		return false
	}

	myFriends := s.friends[from]
	if Contains(myFriends, to) {
		return false
	}

	s.friendshipRequests[from] = append(myRequests, to)
	return true
}

// CheckUsersPassword returns true if user existst and has this password
func (s *InMemoryUsersStore) CheckUsersPassword(user, password string) bool {
	storedPassword, exists := s.users[user]
	return exists && storedPassword == password
}

// RespondToFriendshipRequest responds to a friendship request from otherUser made to user
// Returns false iff friendship request does not exist (in this case no modifications are made)
// Precondition: user and otherUser exist in the DB and have been correctly initialized (ie using AddUser function)
func (s *InMemoryUsersStore) RespondToFriendshipRequest(user, otherUser string, acceptRequest bool) bool {
	requests, hasRequests := s.friendshipRequests[otherUser]
	if !hasRequests {
		return false
	}

	myFriends := s.friends[user]
	theirFriends := s.friends[otherUser]

	for i := 0; i < len(requests); i++ {
		if requests[i] == user {
			s.friendshipRequests[otherUser] = Remove(requests, i)
			if acceptRequest {
				s.friends[user] = append(myFriends, otherUser)
				s.friends[otherUser] = append(theirFriends, user)
			}
			return true
		}
	}

	return false
}

// GetFriends returns the list od friends of a given user
// Precondition: user exists in the DB and has been correctly initialized (ie using AddUser function)
func (s *InMemoryUsersStore) GetFriends(user string) []string {
	return s.friends[user]
}

// --- AUXILIARY FUNCTIONS ---

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

// --- INITIALIZER ---

// EmptyUsersStore returns a new empty InMemoryUsersStore
func EmptyUsersStore() *InMemoryUsersStore {
	store := InMemoryUsersStore{
		users:              map[string]string{},
		friendshipRequests: map[string][]string{},
		friends:            map[string][]string{},
	}
	return &store
}
