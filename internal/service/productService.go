package service

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"instashop/internal/model"
	"instashop/internal/utils"
)

type ProductService struct {
	DB *gorm.DB
}

func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{DB: db}
}

func (s *ProductService) CreateProduct(ctx context.Context, product *model.Product) error {
	// Validate input
	if product.Name == "" || product.Description == "" || product.Price == 0 || product.UserID == 0 {
		utils.NewBadRequestError("name, description, userId, and price are required")
	}

	// Create product
	if err := s.DB.WithContext(ctx).Create(product).Error; err != nil {
		return fmt.Errorf("internal server error: %w", err)
	}

	return nil
}

func (s *ProductService) GetProduct(ctx context.Context, productID uint, userID uint) (*model.Product, error) {
	var product model.Product

	// Fetch the product by productID and UserID
	if err := s.DB.WithContext(ctx).Where("id = ? AND user_id = ?", productID, userID).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewNotFoundError("product not found")
		}
		return nil, fmt.Errorf("internal server error: %w", err)
	}

	return &product, nil
}

func (s *ProductService) GetAllProductsByUserID(ctx context.Context, userID uint) ([]model.Product, error) {
	var products []model.Product

	if err := s.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("internal server error: %w", err)
	}

	return products, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, productID uint, updatedProduct model.Product) (*model.Product, error) {
	var product model.Product

	// Fetch the existing product
	if err := s.DB.WithContext(ctx).First(&product, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewNotFoundError("product not found")
		}
		return nil, fmt.Errorf("internal server error: %w", err)
	}

	// Update the product fields
	product.Name = updatedProduct.Name
	product.Price = updatedProduct.Price

	// Save the updated product
	if err := s.DB.WithContext(ctx).Save(&product).Error; err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return &product, nil
}

func (s *ProductService) DeletePendingProduct(ctx context.Context, productID uint, userID uint) error {
	// Delete the product with the specified productID and userID where status is "pending"
	result := s.DB.WithContext(ctx).Where("id = ? AND user_id = ? AND status = ?", productID, userID, "pending").Delete(&model.Product{})

	// Check for errors
	if result.Error != nil {
		utils.NewBadRequestError("failed to delete product: %w")
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return fmt.Errorf("no pending product found for the given productID and userID")
	}

	return nil
}


