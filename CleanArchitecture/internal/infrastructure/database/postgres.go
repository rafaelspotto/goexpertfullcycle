package database

import (
	"fmt"
	"os"

	"github.com/rafaelspotto/goexpertfullcycle/cleanarchitecture/internal/core/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresRepository struct {
	DB *gorm.DB
}

func NewPostgresRepository() (*PostgresRepository, error) {
	// Get database configuration from environment variables
	host := getEnv("DB_HOST", "localhost")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "orders")
	port := getEnv("DB_PORT", "5432")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

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

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func (r *PostgresRepository) Save(order *domain.Order) error {
	return r.DB.Create(order).Error
}

func (r *PostgresRepository) List() ([]domain.Order, error) {
	var orders []domain.Order
	err := r.DB.Find(&orders).Error
	return orders, err
}
