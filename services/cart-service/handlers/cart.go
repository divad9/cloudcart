package handlers

import (
	"cart-service/models"
	"cart-service/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const cartKeyPrefix = "cart:"

// GetCart retrieves the user's cart
func GetCart(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	cartKey := fmt.Sprintf("%s%v", cartKeyPrefix, userID)

	// Get cart from Redis
	cartData, err := utils.RedisClient.Get(utils.Ctx, cartKey).Result()
	if err != nil {
		// Cart doesn't exist, return empty cart
		cart := models.NewCart(fmt.Sprintf("%v", userID))
		c.JSON(http.StatusOK, gin.H{"cart": cart})
		return
	}

	// Unmarshal cart data
	var cart models.Cart
	if err := json.Unmarshal([]byte(cartData), &cart); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse cart data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cart": cart})
}

// AddItem adds an item to the cart
func AddItem(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: In production, fetch product details from product-service
	// For now, using mock data
	productName := fmt.Sprintf("Product %d", req.ProductID)
	price := 99.99

	cartKey := fmt.Sprintf("%s%v", cartKeyPrefix, userID)

	// Get existing cart or create new one
	var cart models.Cart
	cartData, err := utils.RedisClient.Get(utils.Ctx, cartKey).Result()
	if err != nil {
		// Create new cart
		cart = *models.NewCart(fmt.Sprintf("%v", userID))
	} else {
		// Parse existing cart
		if err := json.Unmarshal([]byte(cartData), &cart); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse cart data"})
			return
		}
	}

	// Check if item already exists in cart
	itemExists := false
	for i, item := range cart.Items {
		if item.ProductID == req.ProductID {
			// Update quantity
			cart.Items[i].Quantity += req.Quantity
			cart.Items[i].Subtotal = float64(cart.Items[i].Quantity) * cart.Items[i].Price
			itemExists = true
			break
		}
	}

	// Add new item if it doesn't exist
	if !itemExists {
		newItem := models.CartItem{
			ProductID:   req.ProductID,
			ProductName: productName,
			Price:       price,
			Quantity:    req.Quantity,
			Subtotal:    float64(req.Quantity) * price,
			AddedAt:     time.Now().Format(time.RFC3339),
		}
		cart.Items = append(cart.Items, newItem)
	}

	// Recalculate totals
	cart.CalculateTotals()

	// Save cart to Redis
	cartJSON, err := json.Marshal(cart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save cart"})
		return
	}

	// Set cart with 24-hour expiration
	err = utils.RedisClient.Set(utils.Ctx, cartKey, cartJSON, 24*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save cart to Redis"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Item added to cart",
		"cart":    cart,
	})
}

// UpdateItem updates the quantity of an item in the cart
func UpdateItem(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	productID := c.Param("product_id")

	var req models.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cartKey := fmt.Sprintf("%s%v", cartKeyPrefix, userID)

	// Get cart
	var cart models.Cart
	cartData, err := utils.RedisClient.Get(utils.Ctx, cartKey).Result()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	if err := json.Unmarshal([]byte(cartData), &cart); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse cart data"})
		return
	}

	// Find and update item
	itemFound := false
	for i, item := range cart.Items {
		if fmt.Sprintf("%d", item.ProductID) == productID {
			if req.Quantity == 0 {
				// Remove item if quantity is 0
				cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			} else {
				// Update quantity
				cart.Items[i].Quantity = req.Quantity
				cart.Items[i].Subtotal = float64(req.Quantity) * cart.Items[i].Price
			}
			itemFound = true
			break
		}
	}

	if !itemFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found in cart"})
		return
	}

	// Recalculate totals
	cart.CalculateTotals()

	// Save cart
	cartJSON, err := json.Marshal(cart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save cart"})
		return
	}

	err = utils.RedisClient.Set(utils.Ctx, cartKey, cartJSON, 24*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cart updated",
		"cart":    cart,
	})
}

// RemoveItem removes an item from the cart
func RemoveItem(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	productID := c.Param("product_id")
	cartKey := fmt.Sprintf("%s%v", cartKeyPrefix, userID)

	// Get cart
	var cart models.Cart
	cartData, err := utils.RedisClient.Get(utils.Ctx, cartKey).Result()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	if err := json.Unmarshal([]byte(cartData), &cart); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse cart data"})
		return
	}

	// Find and remove item
	itemFound := false
	for i, item := range cart.Items {
		if fmt.Sprintf("%d", item.ProductID) == productID {
			cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			itemFound = true
			break
		}
	}

	if !itemFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found in cart"})
		return
	}

	// Recalculate totals
	cart.CalculateTotals()

	// Save cart
	cartJSON, err := json.Marshal(cart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save cart"})
		return
	}

	err = utils.RedisClient.Set(utils.Ctx, cartKey, cartJSON, 24*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Item removed from cart",
		"cart":    cart,
	})
}

// ClearCart clears all items from the cart
func ClearCart(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	cartKey := fmt.Sprintf("%s%v", cartKeyPrefix, userID)

	// Delete cart from Redis
	err := utils.RedisClient.Del(utils.Ctx, cartKey).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cart cleared successfully",
	})
}