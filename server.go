package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/sethigeet/gql-go-auth-backend/graph/generated"
	"github.com/sethigeet/gql-go-auth-backend/graph/resolver"
	"github.com/sethigeet/gql-go-auth-backend/util"
)

// PORT The default port on which the server should
// run on if the no port is specified in the environment
const PORT = "8080"

func main() {
	err := util.LoadEnv(true)
	if err != nil {
		log.Fatalf("Errors while loading the env file: \n%s", err)
		return
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = PORT
	}

	server := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", server)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
