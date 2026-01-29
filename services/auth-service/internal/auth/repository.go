package auth

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id uint) (*User, error)

	CreateEmailVerification(ctx context.Context, emailVerification *EmailVerification) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Detailed Methods

func (r *repository) CreateUser(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) FindByID(ctx context.Context, id uint) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).First(&user, id).Error
	return &user, err
}

func (r *repository) CreateEmailVerification(ctx context.Context, emailVerification *EmailVerification) error {
	return r.db.WithContext(ctx).Create(emailVerification).Error
}
