package graphql

import (
	"time"

	"github.com/rafaelspotto/goexpertfullcycle/cleanarchitecture/internal/core/usecase"
)

type Resolver struct {
	OrderUseCase *usecase.OrderUseCase
}

func NewResolver(useCase *usecase.OrderUseCase) *Resolver {
	return &Resolver{
		OrderUseCase: useCase,
	}
}

type Order struct {
	ID         string    `json:"id"`
	Price      float64   `json:"price"`
	Tax        float64   `json:"tax"`
	FinalPrice float64   `json:"finalPrice"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type CreateOrderInput struct {
	Price float64 `json:"price"`
	Tax   float64 `json:"tax"`
}

func (r *Resolver) Orders() ([]*Order, error) {
	orders, err := r.OrderUseCase.List()
	if err != nil {
		return nil, err
	}

	var result []*Order
	for _, order := range orders {
		result = append(result, &Order{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
			CreatedAt:  order.CreatedAt,
			UpdatedAt:  order.UpdatedAt,
		})
	}

	return result, nil
}

func (r *Resolver) CreateOrder(input CreateOrderInput) (*Order, error) {
	order, err := r.OrderUseCase.Create(input.Price, input.Tax)
	if err != nil {
		return nil, err
	}

	return &Order{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}, nil
}
