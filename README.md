# GoServer
Backend part of an HTTP API of a mock social network with a few basic features.

I have developed this project as a challenge to apply for a job and as an opportunity to learn Golang and Test-Driven Development methodology. See development_plan.md for more info on the development process.

## API
The HTML server accepts different requests.

### GET `/getUsers`
Returns a list of all users in the social network.

### POST `/signUp`
Signs up a new user. Body must contain:
- `user`: username (should be unique and 5-10 alphanumeric characters)
- `pass`: password (should be 8-12 alphanumeric characters)
If preconditions are not met will return HTTP status BadRequest, otherwise should return OK.

### POST `/requestFriendship`
Sends a friendship request. Body must contain:
- `user`: username (should exist)
- `pass`: password (should match user's password)
- `userTo`: username of user to whom we want to send the request (should exist; should not be already friend of user and there should not be a pending friendship request between user and userTo)
If username/password validation fails will return HTTP status Unauthorized, if preconditions are not met will return HTTP status BadRequest, otherwise should return OK.

### POST `/respondToFriendshipRequest`
Responds to a friendship request, either accepting or declining. Body must contain:
- `user`: username (should exist)
- `pass`: password (should match user's password)
- `otherUser`: username of the other user (should exist; should have sent us a friendship request which is still pending)
- `acceptRequest`: either "1" or "0" indicating whether the friendship request is accepted or not
If username/password validation fails will return HTTP status Unauthorized, if preconditions are not met will return HTTP status BadRequest, otherwise should return OK.

### GET `/getFriends/`_\<user\>_
Returns a list of friends of _\<user\>_. If _\<user\>_ does not exist, it will return a HTTP status BadRequest.
