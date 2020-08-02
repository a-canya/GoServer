# GoServer
Backend part of an HTTP API of a mock social network with a few basic features.

I have developed this project as a challenge to apply for a job and as an opportunity to learn Golang and Test-Driven Development methodology. See development_plan.md for more info on the development process.

## How to run
You can compile the code with any go compiler. For instance, in a Unix console go to `/src/main` directory and run `go build` command, then run `./main`. The service will be available at localhost port 5000 (ie you will be able to access it at http://localhost:5000/).

Alternatively you can run the tests directly in an IDE environment. I used VS Code (https://code.visualstudio.com/) with the Go extension (https://marketplace.visualstudio.com/items?itemName=golang.Go) to test all the code.

## API
The HTML server accepts the following requests. Any other request will cause an HTTP status `404 Not Found`.

### GET `/getUsers`
Returns a list of all users in the social network. Unless some problem external to the application happens, this call should always return HTTP status `200 OK`.

### POST `/signUp`
Signs up a new user. Body must contain:
- `user`: username (should be unique and 5-10 alphanumeric characters)
- `pass`: password (should be 8-12 alphanumeric characters)

If preconditions are not met will return HTTP status `400 BadRequest`, otherwise should return `200 OK`.

### POST `/requestFriendship`
Sends a friendship request. Body must contain:
- `user`: username (should exist)
- `pass`: password (should match user's password)
- `userTo`: username of user to whom we want to send the request (should exist; should not be already friend of user and there should not be a pending friendship request between user and userTo)

If username/password validation fails will return HTTP status `401 Unauthorized`, if preconditions are not met will return HTTP status `400 BadRequest`, otherwise should return `200 OK`.

### POST `/respondToFriendshipRequest`
Responds to a friendship request, either accepting or declining. Body must contain:
- `user`: username (should exist)
- `pass`: password (should match user's password)
- `otherUser`: username of the other user (should exist; should have sent us a friendship request which is still pending)
- `acceptRequest`: either "1" or "0" indicating whether the friendship request is accepted or not

If username/password validation fails will return HTTP status `401 Unauthorized`, if preconditions are not met will return HTTP status `400 BadRequest`, otherwise should return `200 OK`.

### GET `/getFriends/`_\<user\>_
Returns a list of friends of _\<user\>_. If _\<user\>_ does not exist, it will return a HTTP status `400 BadRequest`, otherwise should return `200 OK`.
