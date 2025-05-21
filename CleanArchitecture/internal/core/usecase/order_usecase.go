package usecase

import (
	"time"

	"github.com/rafaelspotto/goexpertfullcycle/cleanarchitecture/internal/core/domain"

	"github.com/google/uuid"
)

type OrderUseCase struct {
	OrderRepository domain.OrderRepository
}

func NewOrderUseCase(repository domain.OrderRepository) *OrderUseCase {
	return &OrderUseCase{
		OrderRepository: repository,
	}
}

func (uc *OrderUseCase) Create(price float64, tax float64) (*domain.Order, error) {
	order := &domain.Order{
		ID:         uuid.New().String(),
		Price:      price,
		Tax:        tax,
		FinalPrice: price + tax,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := uc.OrderRepository.Save(order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (uc *OrderUseCase) List() ([]domain.Order, error) {
	return uc.OrderRepository.List()
}
