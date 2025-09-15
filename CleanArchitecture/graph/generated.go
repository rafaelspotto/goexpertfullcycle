package graph

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/rafaelspotto/goexpertfullcycle/cleanarchitecture/internal/interfaces/graphql"
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
	Mutation struct {
		CreateOrder func(childComplexity int, input graphql.CreateOrderInput) int
	}

	Order struct {
		CreatedAt  func(childComplexity int) int
		FinalPrice func(childComplexity int) int
		ID         func(childComplexity int) int
		Price      func(childComplexity int) int
		Tax        func(childComplexity int) int
		UpdatedAt  func(childComplexity int) int
	}

	Query struct {
		Orders func(childComplexity int) int
	}
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
