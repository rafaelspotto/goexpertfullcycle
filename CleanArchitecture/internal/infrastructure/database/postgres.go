package database

import (
	"fmt"

	"github.com/rafaelspotto/goexpertfullcycle/cleanarchitecture/internal/core/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresRepository struct {
	DB *gorm.DB
}

func NewPostgresRepository() (*PostgresRepository, error) {
	dsn := "host=localhost user=postgres password=postgres dbname=orders port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&domain.Order{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	return &PostgresRepository{
		DB: db,
	}, nil
}

func (r *PostgresRepository) Save(order *domain.Order) error {
	return r.DB.Create(order).Error
}

func (r *PostgresRepository) List() ([]domain.Order, error) {
	var orders []domain.Order
	err := r.DB.Find(&orders).Error
	return orders, err
}
