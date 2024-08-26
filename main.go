package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	// load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// get port
	portString := os.Getenv("PORT")
	port, err := strconv.Atoi(portString)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Getting ready to start server on localhost:%d\n", port)
}
