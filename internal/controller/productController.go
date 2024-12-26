package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"instashop/internal/model"
	"instashop/internal/service"
)

type ProductController struct {
	ProductService *service.ProductService
}

func NewProductController(productService *service.ProductService) *ProductController {
	return &ProductController{ProductService: productService}
}

func (ctrl *ProductController) CreateProduct(c *gin.Context) {
	// Extract userID from the context
	userID, exists := c.Get("userID")
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
	if roleStr != "admin" && roleStr != "editor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to create a product"})
		return
	}

	// Bind the request body to the product struct
	var product model.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Set the UserID in the product
	product.UserID = uint(userIDUint)

	// Create the product using the ProductService
	if err := ctrl.ProductService.CreateProduct(c.Request.Context(), &product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully"})
}

func (ctrl *ProductController) GetProduct(c *gin.Context) {
	// Extract userID from the context
	userID, exists := c.Get("userID")
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
	if roleStr != "admin" && roleStr != "editor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to this product"})
		return
	}

	// Extract productID from URL parameters
	productIDStr := c.Param("productID")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID format"})
		return
	}

	// Fetch the product using the ProductService
	product, err := ctrl.ProductService.GetProduct(c.Request.Context(), uint(productID), uint(userIDUint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}

// GetAllProductsByUserID retrieves all products for a specific user
func (ctrl *ProductController) GetAllProductsByUserID(c *gin.Context) {
	// Extract userID from the context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

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
	if roleStr != "admin" && roleStr != "editor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to this product"})
		return
	}

	// Fetch products using the ProductService
	products, err := ctrl.ProductService.GetAllProductsByUserID(c.Request.Context(), uint(userIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

// UpdateProduct updates an existing product
func (ctrl *ProductController) UpdateProduct(c *gin.Context) {
	// Extract userID from the context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	_, err := strconv.ParseUint(userIDStr, 10, 32)
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
	if roleStr != "admin" && roleStr != "editor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to this product"})
		return
	}

	// Extract productID from URL parameters
	productIDStr := c.Param("productID")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID format"})
		return
	}

	// Bind the request body to the updated product struct
	var updatedProduct model.Product
	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Update the product using the ProductService
	product, err := ctrl.ProductService.UpdateProduct(c.Request.Context(), uint(productID), updatedProduct)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}

// DeletePendingProduct deletes a pending product
func (ctrl *ProductController) DeletePendingProduct(c *gin.Context) {
	// Extract userID from the context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

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
	if roleStr != "admin" && roleStr != "editor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to this product"})
		return
	}

	// Extract productID from URL parameters
	productIDStr := c.Param("productID")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID format"})
		return
	}

	// Delete the pending product using the ProductService
	if err := ctrl.ProductService.DeletePendingProduct(c.Request.Context(), uint(productID), uint(userIDUint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
