package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Luiggy102/go-blog-api/app"
)

func main() {

	// Read config
	f, err := os.Open("config.json") // might be configurable
	if err != nil {
		log.Println("Error opening config.json", err)
	}

	config := app.Config{
		DatabaseUrl: "mongodb://localhost:27017/xxxx",
		Addr:        "172.26.0.2:8080",
	}

	err = json.NewDecoder(f).Decode(&config)
	if err != nil {
		log.Println("Error parsing config.json", err)
	}

	fmt.Println("Config:", config)

	// Instantiate an application
	s, err := app.Bootstrap(&config)
	if err != nil {
		log.Fatal(err)
		os.Exit(21)
	}

	// Start server
	fmt.Println("Server started at port", s.Addr)
	err = s.ListenAndServe()
	if err != nil {
		log.Fatalln("Error starting API", err.Error())
	}
}
