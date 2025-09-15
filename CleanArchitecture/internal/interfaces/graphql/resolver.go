package graphql

import (
	"context"

	"github.com/rafaelspotto/goexpertfullcycle/cleanarchitecture/internal/core/usecase"
	"github.com/rafaelspotto/goexpertfullcycle/cleanarchitecture/graph"
)

type Resolver struct {
	OrderUseCase *usecase.OrderUseCase
}

func NewResolver(orderUseCase *usecase.OrderUseCase) *Resolver {
	return &Resolver{
		OrderUseCase: orderUseCase,
	}
}

func (r *Resolver) Mutation() *mutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() *queryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateOrder(ctx context.Context, input graph.CreateOrderInput) (*graph.Order, error) {
	order, err := r.OrderUseCase.Create(input.Price, input.Tax)
	if err != nil {
		return nil, err
	}

	return &graph.Order{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
		CreatedAt:  order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  order.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Orders(ctx context.Context) ([]*graph.Order, error) {
	orders, err := r.OrderUseCase.List()
	if err != nil {
		return nil, err
	}

	var result []*graph.Order
	for _, order := range orders {
		result = append(result, &graph.Order{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
			CreatedAt:  order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:  order.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return result, nil
}
