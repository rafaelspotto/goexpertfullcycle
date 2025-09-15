package graph

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

// NewExecutableSchema creates an ExecutableSchema from the ResolverRoot interface.
func NewExecutableSchema(cfg Config) graphql.ExecutableSchema {
	return &executableSchema{
		resolvers:  cfg.Resolvers,
		directives: cfg.Directives,
		complexity: cfg.Complexity,
	}
}

type Config struct {
	Resolvers  ResolverRoot
	Directives DirectiveRoot
	Complexity ComplexityRoot
}

type ResolverRoot interface {
	Mutation() *mutationResolver
	Query() *queryResolver
}

type DirectiveRoot struct {
}

type ComplexityRoot struct {
}

type CreateOrderInput struct {
	Price float64 `json:"price"`
	Tax   float64 `json:"tax"`
}

type Order struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"finalPrice"`
	CreatedAt  string  `json:"createdAt"`
	UpdatedAt  string  `json:"updatedAt"`
}

type CreateOrderInput struct {
	Price float64 `json:"price"`
	Tax   float64 `json:"tax"`
}

type Order struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"finalPrice"`
	CreatedAt  string  `json:"createdAt"`
	UpdatedAt  string  `json:"updatedAt"`
}

type mutationResolver struct{ *ResolverRoot }
type queryResolver struct{ *ResolverRoot }

type executableSchema struct {
	resolvers  ResolverRoot
	directives DirectiveRoot
	complexity ComplexityRoot
}

func (e *executableSchema) Query(ctx context.Context, op *graphql.OperationDefinition) *graphql.Response {
	return &graphql.Response{
		Data:   []byte(`{"data":{}}`),
		Errors: nil,
	}
}

func (e *executableSchema) Mutation(ctx context.Context, op *graphql.OperationDefinition) *graphql.Response {
	return &graphql.Response{
		Data:   []byte(`{"data":{}}`),
		Errors: nil,
	}
}

func (e *executableSchema) Subscription(ctx context.Context, op *graphql.OperationDefinition) func() *graphql.Response {
	return func() *graphql.Response {
		return &graphql.Response{
			Data:   []byte(`{"data":{}}`),
			Errors: nil,
		}
	}
}

func (e *executableSchema) Complexity(ctx context.Context, typeName, field string, childComplexity int, args map[string]any) (int, bool) {
	return 0, false
}

func (e *executableSchema) Exec(ctx context.Context) graphql.ResponseHandler {
	return func(ctx context.Context) *graphql.Response {
		return &graphql.Response{
			Data:   []byte(`{"data":{}}`),
			Errors: nil,
		}
	}
}
