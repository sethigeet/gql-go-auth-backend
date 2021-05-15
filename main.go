package main

import (
	"log"
	"net/http"
	"os"

	"github.com/sethigeet/gql-go-auth-backend/server"
)

func main() {
	var err error
	err = server.Initialize()

	if err != nil {
		log.Fatalf("Unable to start the server!\nError: %s", err)
		return
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = server.DefualtPort
	}

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Unable to start the server!\nError: %s", err)
	}

	env := os.Getenv("GO_ENV")
	if env == "development" || env == "" {
		log.Printf("Connect to http://localhost:%s/ for GraphQL playground", port)
	}
}
