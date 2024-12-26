package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Username      string         `gorm:"unique;not null" json:"username"`
	Email         string         `gorm:"unique;not null" json:"email"`
	Password      string         `gorm:"not null" json:"password"`
	VerifiedEmail bool           `gorm:"verifiedEmail" json:"verifiedEmail default:false"`
	OtpToken      string         `gorm:"otpToken"`
	ExpiredAt     time.Time      `gorm:"expired_at" json:"expiredAt"`
	Role          string         `gorm:"not null;default:'user'" json:"role"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
