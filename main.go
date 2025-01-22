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
	// dbUrl := os.Getenv("DATABASE_URL")

	// handler func
	// home
	http.HandleFunc("GET /", handlers.HomeHandler())
	http.HandleFunc("POST /posts", handlers.InsertPostHandler())
	// posts

	fmt.Println("Server started at port", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatalln("Error starting API", err.Error())
	}
}
