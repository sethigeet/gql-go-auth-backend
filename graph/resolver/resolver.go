// Package resolver provide the resolver struct for resolving graphql queries
// and mutations
package resolver

//go:generate go run github.com/99designs/gqlgen
import (
	"net/http"

	"gorm.io/gorm"

	"github.com/go-redis/redis/v8"
	"github.com/sethigeet/gql-go-auth-backend/session"
)

// Resolver is the struct on which all the different resolvers are created
// This is where dependency injection is also done for the resolvers
type Resolver struct {
	DB             *gorm.DB
	RDB            *redis.Client
	Writer         http.ResponseWriter
	Request        *http.Request
	SessionManager session.SessionManager
}
