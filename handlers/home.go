package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type homeResponse struct {
	Message string `json:"message"`
}

func HomeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// set the response type
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// send the response
		err := json.NewEncoder(w).Encode(homeResponse{
			Message: "Welcome to my REST API for Blog posting",
		})
		if err != nil {
			log.Fatalln("Error encoding home message response")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
