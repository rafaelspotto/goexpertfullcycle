package graphql

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
