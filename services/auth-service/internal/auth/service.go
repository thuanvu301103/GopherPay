package auth

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// Detailed Methods

func (s *service) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {

	exists, _ := s.repo.FindByEmail(ctx, req.Email)
	if exists != nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Email:    req.Email,
		Password: string(hashedPassword),
		FullName: req.FullName,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken: "generated-jwt-token-here",
		TokenType:   "Bearer",
		User:        *user,
	}, nil
}

func (s *service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {

	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &AuthResponse{
		AccessToken: "generated-jwt-token-here",
		TokenType:   "Bearer",
		User:        *user,
	}, nil
}
