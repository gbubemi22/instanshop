package service

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"instashop/internal/model"
	"instashop/internal/utils"
)

type OrderService struct {
	DB *gorm.DB
}

func NewOderService(db *gorm.DB) *OrderService {
	return &OrderService{DB: db}
}

func (s *OrderService) PlaceOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
	// Validate the order
	if len(order.Products) == 0 {
		utils.NewBadRequestError("order must contain at least one product")
	}

	// Begin a transaction
	tx := s.DB.WithContext(ctx).Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Save the order
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Save the relationship in the many-to-many table
	for _, product := range order.Products {
		err := tx.Model(order).Association("Products").Append(&product)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	fmt.Printf("Order: %w", order)

	return order, nil
}

// ListOrders retrieves all orders for a specific user
func (s *OrderService) ListOrders(ctx context.Context, userID uint) ([]model.Order, error) {
	var orders []model.Order
	// Preload Products if you need associated products in the response
	if err := s.DB.WithContext(ctx).
		Preload("Products").
		Where("user_id = ?", userID).
		Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve orders: %w", err)
	}
	return orders, nil
}

// CancelOrder cancels an order if it is still in Pending status
func (s *OrderService) CancelOrder(ctx context.Context, orderID uint, userID uint) error {
	var order model.Order
	if err := s.DB.WithContext(ctx).First(&order, orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("order not found")
		}
		return fmt.Errorf("failed to retrieve order: %w", err)
	}

	if order.UserID != userID {
		return errors.New("you do not have permission to cancel this order")
	}

	if order.Status != "pending" {
		utils.NewBadRequestError("only pending orders can be canceled")
	}

	order.Status = "canceled"
	if err := s.DB.WithContext(ctx).Save(&order).Error; err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}
	return nil
}

// UpdateOrderStatus updates the status of an order
func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID uint, userID uint, status string) error {
	var order model.Order
	if err := s.DB.WithContext(ctx).First(&order, orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("order not found")
		}
		return fmt.Errorf("failed to retrieve order: %w", err)
	}

	// Optional: Check user permissions or admin privileges
	if order.UserID != userID {
		return errors.New("you do not have permission to update this order")
	}

	// Update the status
	order.Status = model.OrderStatusType(status)
	if err := s.DB.WithContext(ctx).Save(&order).Error; err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	return nil
}
