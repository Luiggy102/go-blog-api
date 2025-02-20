package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Luiggy102/go-blog-api/database"
	"github.com/Luiggy102/go-blog-api/handlers"
)

type Config struct {
	DatabaseUrl string `json:"database_url"`
	Addr        string `json:"addr"`
}

// Bootstrap: setup and configure application
func Bootstrap(config *Config) (*http.Server, error) {

	// Instantiate mongodb
	mongo, err := database.NewMongoDb(config.DatabaseUrl)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Setup routes
	mux := http.NewServeMux()

	// home
	mux.HandleFunc("GET /", handlers.HomeHandler(mongo))

	// api posts
	mux.HandleFunc("POST /posts", handlers.InsertPostHandler(mongo))
	mux.HandleFunc("GET /posts", handlers.GetPostsHandler(mongo))
	mux.HandleFunc("GET /posts/{id}", handlers.GetPostsbyIdHandler(mongo))
	mux.HandleFunc("PUT /posts/{id}", handlers.UpdatePostHander(mongo))
	mux.HandleFunc("DELETE /posts/{id}", handlers.DeletePostHandler(mongo))

	return &http.Server{
		Addr:    config.Addr,
		Handler: MiddlewareAccessLog(mux),
	}, nil
}

func MiddlewareAccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}
