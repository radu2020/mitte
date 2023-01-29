# Dating APP Api

Simple API which creates user profiles and allows retrieving other profiles
 and swiping/matching profiles.

## Endpoints

### /user/create: create random user profile
Example:

```curl localhost:3000/user/create```

### /profiles/{id}: retrieve matching profiles
Example:
```
curl -H 'Content-Type: application/json' \
   -d '{"api_key":"secretApiKey"}' \
   -X POST \
   localhost:3000/profiles/1
```

### /swipe: swipe on a profile
Example:
```
curl -H 'Content-Type: application/json' \
   -d '{ "user_id":1, "potential_match_id":2, "preference":true, "api_key":"secretApiKey"}' \
   -X POST \
   localhost:3000/swipe
```

### /login: user login
Example:
```
curl -H 'Content-Type: application/json' \
   -d '{ "email":"Natasha@gmail.com", "password":"Natashas-safe-password"}' \
   -X POST \
   localhost:3000/login
```

## Setup

### Clone the repo
`git clone https://github.com/radu2020/mitte.git`

### Pulling dependencies
```
go get github.com/mattn/go-sqlite3
go get "github.com/gorilla/mux"
```

### Running locally
```
go run main.go
Starting the server on localhost:3000...
```