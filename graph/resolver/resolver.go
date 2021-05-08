// Package resolver provide the resolver struct for resolving graphql queries
// and mutations
// The Resolver struct also handles dependency injection for the resolvers
package resolver

//go:generate go run github.com/99designs/gqlgen
import "github.com/sethigeet/gql-go-auth-backend/graph/model"

type Resolver struct {
	todos []*model.Todo
	users []*model.User
}
