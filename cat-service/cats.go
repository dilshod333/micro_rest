package main

import (
	"conn/connection"
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
	http.HandleFunc("/cat", Send)
	http.HandleFunc("/cats", All)
	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}

func Send(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	client := &http.Client{}
	stringURL := "https://catfact.ninja/fact"

	resp, err := client.Get(stringURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var m Message
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		http.Error(w, "Error while decoding", http.StatusInternalServerError)
		return
	}

	db, err := connection.Initialize()
	if err != nil {
		http.Error(w, "Error while connecting to the database", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO cat(fact, length) VALUES($1, $2)", m.Fact, m.Length)
	if err != nil {
		http.Error(w, "Error while inserting data", http.StatusInternalServerError)
		return
	}


	if err != nil {
		http.Error(w, "Error retrieving last insert ID", http.StatusInternalServerError)
		return
	}


	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)


}

func All(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	db, err := connection.Initialize()
	if err != nil {
		http.Error(w, "Error while connecting to the database", http.StatusInternalServerError)
		return
	}

	rows, err := db.Query("SELECT id, fact, length FROM cat")
	if err != nil {
		http.Error(w, "Error while selecting data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		if err = rows.Scan(&m.ID, &m.Fact, &m.Length); err != nil {
			http.Error(w, "Error while scanning", http.StatusInternalServerError)
			return
		}
		messages = append(messages, m)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Error with rows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
