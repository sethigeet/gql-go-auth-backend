// Package server provides a function to start up the server along with all the
// other things such as the db connections, etc.
package server

import (
	"fmt"
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

// DefualtPort The default port on which the server should
// run on if the no port is specified in the environment
const DefualtPort = "4000"

func Initialize() error {
	var err error
	err = util.LoadEnv(true)
	if err != nil {
		return fmt.Errorf("errors while loading the env file: \n%s", err)
	}

	db, rdb, err := database.Connect(true)
	if err != nil {
		return fmt.Errorf("errors while connecting to the database: \n%s", err)
	}

	gqlServer := getServer(db, rdb)

	env := os.Getenv("GO_ENV")
	if env == "development" || env == "" {
		http.Handle("/api", playground.Handler("GraphQL playground", "/api/query"))
	}
	http.HandleFunc("/api/query", gqlServer)

	return nil
}

func getServer(db *gorm.DB, rdb *redis.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resolvers := resolver.Resolver{
			DB:      db,
			Writer:  w,
			Request: r,
			RDB:     rdb,
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
