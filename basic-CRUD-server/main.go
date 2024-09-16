package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

//create a map to hold users and a mutex for thread safety

var users = make(map[int]User)
var mu sync.Mutex
var nextID = 1

// create handlers for CRUD operatios
func createUser(w http.ResponseWriter, r *http.Request) {
	// The mutex (mu) is used to ensure thread safety when multiple requests are handled by the server concurrently.
	mu.Lock()
	defer mu.Unlock()

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		// reads the request body and attempts to decode the JSON into the User struct.
		// If the input is invalid (e.g., malformed JSON), an error is returned.
		http.Error(w, "Invalid input", http.StatusBadRequest)
		// If there's an error during decoding, the server responds with a 400 Bad Request error
		// and the message "Invalid input."
		return
	}

	user.ID = nextID
	// every user needs unique ID. The nextID variable tracks the next available ID
	nextID++
	// nextID is incremented so that future users get unique IDs.

	users[user.ID] = user
	// The users map acts as an in-memory database.
	// The new user is added to the users map, with the ID as the key.

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
	// w.WriteHeader(http.StatusCreated) sends a 201 Created HTTP status code,
	// indicating that the resource was successfully created.

	// json.NewEncoder(w).Encode(user) serializes the created user struct into JSON
	// and writes it to the response body, which will be sent back to the client.
}

func getUser(w http.ResponseWriter, r *http.Request) {

}

func main() {

}
