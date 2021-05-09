package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/sethigeet/gql-go-auth-backend/database"
	"github.com/sethigeet/gql-go-auth-backend/graph/generated"
	"github.com/sethigeet/gql-go-auth-backend/graph/resolver"
	"github.com/sethigeet/gql-go-auth-backend/util"
)

// PORT The default port on which the server should
// run on if the no port is specified in the environment
const PORT = "4000"

func main() {
	var err error
	err = util.LoadEnv(true)
	if err != nil {
		log.Fatalf("Errors while loading the env file: \n%s", err)
		return
	}

	db, err := database.Connect(true)
	if err != nil {
		log.Fatalf("Errors while connecting to the database: \n%s", err)
		return
	}

	resolvers := resolver.Resolver{DB: db}
	server := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers}))

	http.Handle("/api", playground.Handler("GraphQL playground", "/api/query"))
	http.Handle("/api/query", server)

	port := os.Getenv("PORT")
	if port == "" {
		port = PORT
	}

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
