package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/securecookie"
	"gorm.io/gorm"

	"github.com/sethigeet/gql-go-auth-backend/database"
	"github.com/sethigeet/gql-go-auth-backend/graph/generated"
	"github.com/sethigeet/gql-go-auth-backend/graph/resolver"
	"github.com/sethigeet/gql-go-auth-backend/session"
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

	db, rdb, err := database.Connect(true)
	if err != nil {
		log.Fatalf("Errors while connecting to the database: \n%s", err)
		return
	}

	server := getServer(db, rdb)
	http.Handle("/api", playground.Handler("GraphQL playground", "/api/query"))
	http.HandleFunc("/api/query", server)

	port := os.Getenv("PORT")
	if port == "" {
		port = PORT
	}

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getServer(db *gorm.DB, rdb *redis.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resolvers := resolver.Resolver{
			DB:      db,
			Writer:  w,
			Request: r,
			SessionManager: session.SessionManager{
				RDB:     rdb,
				Writer:  w,
				Request: r,
				CookieSecurer: securecookie.New(
					[]byte(os.Getenv("SESSION_SECRET_HASH")),
					[]byte(os.Getenv("SESSION_SECRET_BLOCK")),
				),
			},
		}
		server := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolvers}))
		server.ServeHTTP(w, r)
	}
}
