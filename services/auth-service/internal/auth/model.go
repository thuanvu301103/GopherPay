package auth

import (
	"time"

	"gorm.io/gorm"
)

// Data Entitíe
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"unique;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	FullName  string         `gorm:"size:100" json:"fullName"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Dtos
type RegisterRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password"`
	FullName        string `json:"fullName" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	User        User   `json:"user"`
}
