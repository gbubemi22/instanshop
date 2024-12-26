package model

import (
	"time"
)

// StatusType defines the type for product status.
type StatusType string

// Enumeration of product statuses.
const (
	StatusPending  StatusType = "pending"
	StatusDeclined StatusType = "declined"
	StatusApproved StatusType = "approved"
)

// Product represents a product in the system.
type Product struct {
	ID          uint       `gorm:"primaryKey"`
	UserID      uint       `gorm:"not null"`
	Name        string     `gorm:"size:255;not null"`
	Description string     `gorm:"type:text;not null"`
	Price       float64    `gorm:"type:decimal(10,2);not null"`
	Status      StatusType `gorm:"type:varchar(10);default:'pending';not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
