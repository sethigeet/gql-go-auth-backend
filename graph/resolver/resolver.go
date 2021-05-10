// Package resolver provide the resolver struct for resolving graphql queries
// and mutations
// The Resolver struct also handles dependency injection for the resolvers
package resolver

//go:generate go run github.com/99designs/gqlgen
import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Resolver struct {
	DB  *gorm.DB
	RDB *redis.Client
}
