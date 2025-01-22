package handlers

import (
	"encoding/json"
	"net/http"
)

type homeResponse struct {
	Message string `json:"message"`
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// set the response type
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// send the response
	json.NewEncoder(w).Encode(homeResponse{
		Message: "Welcome to my REST API for Blog posting",
	})
}
