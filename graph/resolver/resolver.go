// Package resolver provide the resolver struct for resolving graphql queries
// and mutations
// The Resolver struct also handles dependency injection for the resolvers
package resolver

//go:generate go run github.com/99designs/gqlgen
import (
	"net/http"

	"gorm.io/gorm"

	"github.com/sethigeet/gql-go-auth-backend/session"
)

type Resolver struct {
	DB             *gorm.DB
	Writer         http.ResponseWriter
	Request        *http.Request
	SessionManager session.SessionManager
}
