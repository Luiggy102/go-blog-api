package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Luiggy102/go-blog-api/handlers"
	"github.com/joho/godotenv"
)

func main() {
	// load env vars
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalln("Error loading env vars", err.Error())
	}
	port := os.Getenv("PORT")

	// handler func
	// home
	http.HandleFunc("GET /", handlers.HomeHandler())

	// posts
	http.HandleFunc("POST /posts", handlers.InsertPostHandler())
	http.HandleFunc("GET /posts", handlers.GetPostsHandler())
	http.HandleFunc("GET /posts/{id}", handlers.GetPostsbyIdHandler())
	http.HandleFunc("PUT /posts/{id}", handlers.UpdatePostHander())
	http.HandleFunc("DELETE /posts/{id}", handlers.DeletePostHandler())

	fmt.Println("Server started at port", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatalln("Error starting API", err.Error())
	}
}
