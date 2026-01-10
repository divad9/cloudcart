package models

import "time"

// CartItem represents a single item in the cart
type CartItem struct {
	ProductID   int     `json:"product_id"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Subtotal    float64 `json:"subtotal"`
	AddedAt     string  `json:"added_at"`
}

// Cart represents a user's shopping cart
type Cart struct {
	UserID     string     `json:"user_id"`
	Items      []CartItem `json:"items"`
	TotalItems int        `json:"total_items"`
	TotalPrice float64    `json:"total_price"`
	UpdatedAt  string     `json:"updated_at"`
}

// AddItemRequest represents the request to add an item
type AddItemRequest struct {
	ProductID int `json:"product_id" binding:"required"`
	Quantity  int `json:"quantity" binding:"required,min=1"`
}

// UpdateItemRequest represents the request to update item quantity
type UpdateItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=0"`
}

// NewCart creates a new empty cart
func NewCart(userID string) *Cart {
	return &Cart{
		UserID:     userID,
		Items:      []CartItem{},
		TotalItems: 0,
		TotalPrice: 0,
		UpdatedAt:  time.Now().Format(time.RFC3339),
	}
}

// CalculateTotals recalculates cart totals
func (c *Cart) CalculateTotals() {
	c.TotalItems = 0
	c.TotalPrice = 0

	for _, item := range c.Items {
		c.TotalItems += item.Quantity
		c.TotalPrice += item.Subtotal
	}

	c.UpdatedAt = time.Now().Format(time.RFC3339)
}