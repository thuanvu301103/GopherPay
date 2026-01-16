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
	// User's email
	Email string `json:"email" example:"john@gopher.com" binding:"required,email"`
	// User's password
	Password string `json:"password" example:"123456789" binding:"required,min=8"`
	// Re-enter password - Must equal password
	ConfirmPassword string `json:"confirmPassword" example:"123456789" binding:"required,eqfield=Password"`
	// User's full name
	FullName string `json:"fullName" example:"John Doe" binding:"required"`
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
