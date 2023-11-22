package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Gym struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var gyms []Gym

func main() {
	http.HandleFunc("/gyms", getGyms)
	http.HandleFunc("/gyms/add", addGym)

	fmt.Println("Server listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}

func getGyms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gyms)
}

func addGym(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var gym Gym
	err := json.NewDecoder(r.Body).Decode(&gym)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	gyms = append(gyms, gym)
	json.NewEncoder(w).Encode(gym)
}
