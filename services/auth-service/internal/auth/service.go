package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/thuanvu301103/auth-service/internal/config"
	"github.com/thuanvu301103/auth-service/internal/infrastructure/kafka"
	"github.com/thuanvu301103/auth-service/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error)
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
}

type service struct {
	repo          Repository
	kafkaProducer *kafka.Producer
	cfg           config.Config
}

func NewService(repo Repository, kafkaProducer *kafka.Producer, cfg config.Config) Service {
	return &service{
		repo:          repo,
		kafkaProducer: kafkaProducer,
		cfg:           cfg,
	}
}

// Detailed Methods

func (s *service) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {

	// 1. Check email exist
	exists, _ := s.repo.FindByEmail(ctx, req.Email)
	if exists != nil {
		return nil, errors.New("user with this email already exists")
	}

	// 2. Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 3. Create User
	user := &User{
		Email:    req.Email,
		Password: string(hashedPassword),
		FullName: req.FullName,
		// IsActive and IsDeleted is false by default
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	// 4. Create Verification Token pair
	rawToken, hashedToken, err := utils.GenerateVerificationToken()
	if err != nil {
		return nil, err
	}

	// 5. Store and send Token
	verification := &EmailVerification{
		UserID:    user.ID,
		Token:     hashedToken,
		ExpiresAt: time.Now().Add(time.Duration(s.cfg.EmailVerifyExpireTime) * time.Minute),
	}

	if err := s.repo.CreateEmailVerification(ctx, verification); err != nil {
		return nil, err
	}

	payload := map[string]string{
		"email":  user.Email,
		"userId": user.ID.String(),
		"token":  rawToken,
	}
	byteData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	s.kafkaProducer.Publish(ctx, s.cfg.KafkaEmailVerifyTopic, []byte(user.Email), byteData)

	// 6. Return response
	return &RegisterResponse{
		Message: fmt.Sprintf(
			"Registration successful. Please check your email for verification code. The code will be valid within %d minutes.",
			s.cfg.EmailVerifyExpireTime,
		),
		User: *user,
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
