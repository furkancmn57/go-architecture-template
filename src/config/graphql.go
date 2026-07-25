package config

// GraphQL holds optional GraphQL transport settings.
// Disabled by default — enable with GRAPHQL_ENABLED=true when you want /graphql.
type GraphQL struct {
	Enabled bool `env:"GRAPHQL_ENABLED" envDefault:"false"`
}
