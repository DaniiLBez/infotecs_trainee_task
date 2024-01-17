package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"infotecs_trainee_task/internal/entity"
	"infotecs_trainee_task/internal/repo"
	"infotecs_trainee_task/pkg/hasher"
	"log/slog"
	"time"
)

type TokenClaims struct {
	jwt.RegisteredClaims
	UserUUID uuid.UUID
}

type AuthService struct {
	userRepo       repo.User
	passwordHasher hasher.PasswordHasher
	signKey        string
	tokenTTL       time.Duration
}

func NewAuthService(
	userRepo repo.User,
	hasher hasher.PasswordHasher,
	signKey string,
	tokenTTL time.Duration,
) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		passwordHasher: hasher,
		signKey:        signKey,
		tokenTTL:       tokenTTL,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, input struct {
	username string
	password string
}) (uuid.UUID, error) {
	user := entity.User{
		Username: input.username,
		Password: s.passwordHasher.Hash(input.password),
	}

	userUUID, err := s.userRepo.CreateUser(ctx, user)

	if err != nil {
		if errors.Is(err, repo.ErrAlreadyExist) {
			return uuid.Nil, ErrUserAlreadyExists
		}
		slog.Error("AuthService.CreateUser", err.Error())
		return uuid.Nil, ErrCannotCreateUser
	}

	return userUUID, nil
}

func (s *AuthService) GenerateToken(ctx context.Context, input struct {
	username string
	password string
}) (string, error) {
	user, err := s.userRepo.GetUserByUsernameAndPassword(ctx, input.username, s.passwordHasher.Hash(input.password))

	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return "", ErrUserNotFound
		}
		slog.Error("AuthService.GenerateToken", err.Error())
		return "", ErrCannotGetUser
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		UserUUID: user.UUID,
	})

	tokenString, err := token.SignedString([]byte(s.signKey))
	if err != nil {
		slog.Error("AuthService.GenerateToken: can not sign key", err.Error())
		return "", ErrCannotSignToken
	}

	return tokenString, nil
}

func (s *AuthService) ParseToken(accessToken string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return uuid.Nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.signKey), nil
	})

	if err != nil {
		return uuid.Nil, ErrCannotParseToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return uuid.Nil, ErrCannotParseToken
	}

	return claims.UserUUID, nil
}
