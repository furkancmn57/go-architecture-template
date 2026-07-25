package extensions

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	gql "github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

// RegisterGraphQL mounts the GraphQL HTTP endpoint (and GraphiQL UI) at /graphql.
// Call only when config.GraphQL.Enabled is true — REST remains the default transport.
func RegisterGraphQL(app *fiber.App, schema gql.Schema) {
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})
	app.All("/graphql", adaptor.HTTPHandler(h))
}
