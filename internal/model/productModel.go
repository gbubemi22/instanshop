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
	ID          uint       `gorm:"primaryKey" json:"id"`
	UserID      uint       `gorm:"not null" json:"user_id"`
	Name        string     `gorm:"size:255;not null" json:"name"`
	Description string     `gorm:"type:text;not null" json:"description"`
	Price       float64    `gorm:"type:decimal(10,2);not null" json:"price"`
	Status      StatusType `gorm:"type:varchar(10);default:'pending';not null" json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
 }
 
