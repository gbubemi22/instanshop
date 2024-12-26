package model

import (
	"time"
)

type OrderStatusType string

// Enumeration of order statuses.
const (
	OrderStatusPending  OrderStatusType = "pending"
	OrderStatusDeclined OrderStatusType = "declined"
	OrderStatusApproved OrderStatusType = "approved"
)

type Order struct {
	ID        uint            `json:"id" gorm:"primaryKey"`
	UserID    uint            `json:"user_id"`
	Products  []Product       `json:"products" gorm:"many2many:order_products;"`
	Status    OrderStatusType `json:"status" gorm:"type:varchar(10);default:'pending';not null"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
