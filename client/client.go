package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Message struct {
	ID     int     `json:"id"`
	Fact   string  `json:"fact"`
	Length float32 `json:"length"`
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/client", HandleClient)
	router.HandleFunc("/client-all", GetAll)
	log.Println("Client server started at :8081")
	http.ListenAndServe(":8081", router)
}

func HandleClient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Wrong request method", http.StatusMethodNotAllowed)
		return
	}

	url := "http://localhost:8080/cat"
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Error fetching data from server", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var message Message
	err = json.NewDecoder(resp.Body).Decode(&message)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Wrong request method", http.StatusMethodNotAllowed)
		return
	}

	url := "http://localhost:8080/cats"
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Error fetching data from server", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var messages []Message
	err = json.NewDecoder(resp.Body).Decode(&messages)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
