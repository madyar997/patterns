package main

import (
	"encoding/json"
	"log"
	"net/http"
	"rate-limit/limit"
)

func main() {
	http.Handle("/ping", limit.RateLimiter(endpointHandler))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("There was an error listening on port :8080", err)
	}
}

func endpointHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	message := Message{
		Status: "Successful",
		Body:   "Hello",
	}
	err := json.NewEncoder(writer).Encode(&message)
	if err != nil {
		return
	}
}

type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}
