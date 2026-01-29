package auth

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Data Entities
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Email     string    `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	FullName  string    `gorm:"size:100" json:"fullName"`
	IsActive  bool      `gorm:"index;default:false" json:"isActive"`
	IsDeleted bool      `gorm:"index;default:false" json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	var err error
	u.ID, err = uuid.NewV7()
	return err
}

type EmailVerification struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"not null;index"`
	Token     string    `gorm:"not null; index"`
	CreatedAt time.Time
	ExpiresAt time.Time `gorm:"not null"`

	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (e *EmailVerification) BeforeCreate(tx *gorm.DB) error {
	var err error
	e.ID, err = uuid.NewV7()
	return err
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

type RegisterResponse struct {
	// Successfully registered message
	Message string `json:"message" example:"Registration successful. Please check your email" binding:"required"`
	// User's information
	User User `json:"user" binding:"required"`
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
