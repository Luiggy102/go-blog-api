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

	config := app.Config{}
	// running in a docker (compose) check
	// default local config
	if os.Getenv("DB_URL") == "" && os.Getenv("ADDR") == "" {
		err = json.NewDecoder(f).Decode(&config)
		if err != nil {
			log.Println("Error parsing config.json", err)
		}
	} else { // docker config
		config.Addr = os.Getenv("ADDR")
		config.DatabaseUrl = os.Getenv("DB_URL")
	}

	fmt.Println(config)

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
