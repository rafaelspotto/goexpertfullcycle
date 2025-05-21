package domain

import (
	"time"
)

type Order struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	Price      float64   `json:"price"`
	Tax        float64   `json:"tax"`
	FinalPrice float64   `json:"final_price"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type OrderRepository interface {
	Save(order *Order) error
	List() ([]Order, error)
}

type OrderUseCase interface {
	Create(price float64, tax float64) (*Order, error)
	List() ([]Order, error)
}
