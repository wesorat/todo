package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"example/todo/internal/domain"
	"example/todo/internal/repository"
	"log/slog"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// отзыв токена из бд
// парсинг токена
// хэширование рефреша
// сохранение в редис рефреша

const (
	tokenTTL        = time.Minute * 15
	refreshTokenTTL = time.Minute * 30 * 24
)

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidToken    = errors.New("invalid token")
)

type authService struct {
	repo repository.AuthRepository
	log  *slog.Logger
}

func NewAuthService(repo repository.AuthRepository, log *slog.Logger) *authService {
	return &authService{repo: repo, log: log}
}

func (s *authService) CreateUser(user domain.CreateUser) (int, error) {
	password_hash, err := generatePasswordHash(user.Password)
	if err != nil {
		s.log.Error("The password cannot be empty", slog.Any("err", err))
		return 0, err
	}
	user.Password = password_hash
	id, err := s.repo.CreateUser(user)
	if err != nil {
		s.log.Error("Cannot create user", slog.Any("err", err))
		return 0, err
	}

	return id, nil
}

func (s *authService) GetUser(name, password string) (domain.User, error) {
	user, err := s.repo.GetUser(name, password)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return domain.User{}, repository.ErrUserNotFound
		}
		return domain.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidPassword
	}
	return user, nil

}

func (s *authService) SignIn(name, password string) (Tokens, error) {
	user, err := s.repo.GetUser(name, password)
	if err != nil {
		s.log.Error(err.Error())
		return Tokens{}, err
	}
	refresh, err := s.generateRefreshToken()
	if err != nil {
		s.log.Error(err.Error())
		return Tokens{}, err
	}
	refresh_hash, err := generateRefreshHash(refresh)
	if err != nil {
		s.log.Error(err.Error())
		return Tokens{}, err

	}
	expired_at := time.Now().Add((refreshTokenTTL))
	if err := s.repo.SaveRefresh(user.ID, refresh_hash, expired_at); err != nil {
		s.log.Error(err.Error())
		return Tokens{}, err

	}

	jwt, err := s.generateJWT(user.ID)
	if err != nil {
		s.log.Error(err.Error())
		return Tokens{}, err
	}
	return Tokens{RefreshToken: refresh, JWT: jwt}, nil
}

func (s *authService) generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil

}

type StandartClaimsWithUserID struct {
	jwt.StandardClaims
	UserID int `json:"user_id"`
}

func (s *authService) generateJWT(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &StandartClaimsWithUserID{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userID,
	})
	signingKey := []byte(os.Getenv("signingKey"))
	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (s *authService) ParseJWT(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &StandartClaimsWithUserID{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidToken
			}
			return []byte(os.Getenv("signingKey")), nil
		})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*StandartClaimsWithUserID)
	if !token.Valid || !ok {
		return 0, ErrInvalidToken
	}
	return claims.UserID, nil

}

func generatePasswordHash(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func generateRefreshHash(refersh string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(refersh), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
