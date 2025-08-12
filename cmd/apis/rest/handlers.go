package main

import "net/http"

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	// Handle GET /v1/status
}

func NewMessageHandler(w http.ResponseWriter, r *http.Request) {
	// Handle POST /v1/message/new/{client_app}
}

func ResultsHandler(w http.ResponseWriter, r *http.Request) {
	// Handle GET /v1/results/{client_app}
}
