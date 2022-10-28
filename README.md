# rest-api-task
Implementation of Task 1 rest API

## Server

Server can be run from `server/main.go`

Server runs on localhost:8080 and has two endpoints /users and /apps

Both endpoints accept GET, POST and DELETE requests

GET requests need to have name values in the query parameters

POST requests need to have a Json payload with the specific object fields inside

DELETE requests need to have name values in the query parameters

## Client

By running `client/main.go` these commands can be run: get-age, get-app, cr-user, cr-app, del-user, del-app.

## TODO

- [ ] Add more tests
- [ ] Make writing to disk more robust
- [ ] Improve server logging
