package graphql

import (
	"fmt"

	gql "github.com/graphql-go/graphql"

	todoservice "github.com/furkancmn57/go-base-template/src/services/todo"
)

// NewSchema builds the root GraphQL schema. Resolvers call the same services
// as REST controllers — no separate business layer for GraphQL.
func NewSchema(todoSvc *todoservice.Service) (gql.Schema, error) {
	schema, err := gql.NewSchema(gql.SchemaConfig{
		Query: gql.NewObject(gql.ObjectConfig{
			Name:   "Query",
			Fields: todoFields(todoSvc),
		}),
		Mutation: gql.NewObject(gql.ObjectConfig{
			Name:   "Mutation",
			Fields: todoMutations(todoSvc),
		}),
	})
	if err != nil {
		return gql.Schema{}, fmt.Errorf("graphql: schema: %w", err)
	}
	return schema, nil
}
