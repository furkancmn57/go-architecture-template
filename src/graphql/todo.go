package graphql

import (
	"github.com/graphql-go/graphql"

	"github.com/furkancmn57/go-base-template/src/models/requests"
	todoservice "github.com/furkancmn57/go-base-template/src/services/todo"
)

var todoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Todo",
	Fields: graphql.Fields{
		"id":          &graphql.Field{Type: graphql.NewNonNull(graphql.ID)},
		"title":       &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"description": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"completed":   &graphql.Field{Type: graphql.NewNonNull(graphql.Boolean)},
		"createdAt":   &graphql.Field{Type: graphql.NewNonNull(graphql.DateTime)},
		"updatedAt":   &graphql.Field{Type: graphql.NewNonNull(graphql.DateTime)},
	},
})

func todoFields(svc *todoservice.Service) graphql.Fields {
	return graphql.Fields{
		"todos": &graphql.Field{
			Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(todoType))),
			Resolve: func(p graphql.ResolveParams) (any, error) {
				rows, appErr := svc.Todos(p.Context)
				if appErr != nil {
					return nil, appErr
				}
				return rows, nil
			},
		},
		"todo": &graphql.Field{
			Type: todoType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
			},
			Resolve: func(p graphql.ResolveParams) (any, error) {
				resp, appErr := svc.TodoById(p.Context, p.Args["id"].(string))
				if appErr != nil {
					return nil, appErr
				}
				return resp, nil
			},
		},
	}
}

func todoMutations(svc *todoservice.Service) graphql.Fields {
	return graphql.Fields{
		"createTodo": &graphql.Field{
			Type: graphql.NewNonNull(todoType),
			Args: graphql.FieldConfigArgument{
				"title":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"description": &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: func(p graphql.ResolveParams) (any, error) {
				req := requests.CreateTodoRequest{Title: p.Args["title"].(string)}
				if d, ok := p.Args["description"].(string); ok {
					req.Description = d
				}
				resp, appErr := svc.Create(p.Context, req)
				if appErr != nil {
					return nil, appErr
				}
				return resp, nil
			},
		},
		"updateTodo": &graphql.Field{
			Type: graphql.NewNonNull(todoType),
			Args: graphql.FieldConfigArgument{
				"id":          &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
				"title":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"description": &graphql.ArgumentConfig{Type: graphql.String},
				"completed":   &graphql.ArgumentConfig{Type: graphql.Boolean},
			},
			Resolve: func(p graphql.ResolveParams) (any, error) {
				req := requests.UpdateTodoRequest{Title: p.Args["title"].(string)}
				if d, ok := p.Args["description"].(string); ok {
					req.Description = d
				}
				if c, ok := p.Args["completed"].(bool); ok {
					req.Completed = c
				}
				resp, appErr := svc.Update(p.Context, p.Args["id"].(string), req)
				if appErr != nil {
					return nil, appErr
				}
				return resp, nil
			},
		},
		"completeTodo": &graphql.Field{
			Type: graphql.NewNonNull(todoType),
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
			},
			Resolve: func(p graphql.ResolveParams) (any, error) {
				resp, appErr := svc.Complete(p.Context, p.Args["id"].(string))
				if appErr != nil {
					return nil, appErr
				}
				return resp, nil
			},
		},
		"deleteTodo": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Boolean),
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
			},
			Resolve: func(p graphql.ResolveParams) (any, error) {
				if appErr := svc.Delete(p.Context, p.Args["id"].(string)); appErr != nil {
					return false, appErr
				}
				return true, nil
			},
		},
	}
}
