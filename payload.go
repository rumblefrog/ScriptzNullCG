package main

import "net/http"

// Client - HTTP client for sending payload
var Client = &http.Client{
	Timeout: 30000000000,
}

// CreatePayload - Constructs payload for POST form query
func CreatePayload(t *Thread, r *Reply) {
	// request, err := http.NewRequest()
	// PayloadRequest.URL = URL.
}
