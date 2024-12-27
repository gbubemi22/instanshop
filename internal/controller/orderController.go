package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"instashop/internal/model"
	"instashop/internal/service"
)

type OrderController struct {
	OrderService *service.OrderService
}

func NewOrderController(orderService *service.OrderService) *OrderController {
	return &OrderController{OrderService: orderService}
}

func (ctrl *OrderController) PlaceOrderHandler(c *gin.Context) {
	// Extract userID from the context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	// Convert userIDStr to uint
	userIDUint, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Extract role from the context
	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "Role not found"})
		return
	}

	roleStr, ok := role.(string)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid role"})
		return
	}

	// Check if the user has the required role to create a product
	if roleStr != "user" && roleStr != "editor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to place an order"})
		return
	}

	// Parse the order payload
	var order model.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Set the UserID to the extracted user ID
	order.UserID = uint(userIDUint)

	// Call the service to place the order
	placedOrder, err := ctrl.OrderService.PlaceOrder(c, &order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to place order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"order": placedOrder})
}


// ListOrdersHandler handles fetching all orders for a specific user
func (ctrl *OrderController) ListOrdersHandler(c *gin.Context) {
	// Extract userID from the context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDUint, err := strconv.ParseUint(userID.(string), 10, 32)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}


	// Extract role from the context
	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "Role not found"})
		return
	}

	roleStr, ok := role.(string)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid role"})
		return
	}

	// Check if the user has the required role to create a product
	if roleStr != "user" && roleStr != "editor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to fetch  orders"})
		return
	}

	// Call the service to list orders
	orders, err := ctrl.OrderService.ListOrders(c, uint(userIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

// CancelOrderHandler handles canceling an order
func (ctrl *OrderController) CancelOrderHandler(c *gin.Context) {
	// Extract userID from the context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDUint, err := strconv.ParseUint(userID.(string), 10, 32)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	// Extract role from the context
	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "Role not found"})
		return
	}

	roleStr, ok := role.(string)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid role"})
		return
	}

	// Check if the user has the required role to create a product
	if roleStr != "user" && roleStr != "editor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to cancel an order"})
		return
	}

	// Extract orderID from the URL
	orderIDStr := c.Param("orderID")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	// Call the service to cancel the order
	err = ctrl.OrderService.CancelOrder(c, uint(orderID), uint(userIDUint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order canceled successfully"})
}