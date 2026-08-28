package service

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log/slog"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/wesorat/todo/internal/domain"
	"github.com/wesorat/todo/internal/repository"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

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
	repo          repository.AuthRepository
	log           *slog.Logger
	redis         *redis.Client
	signingKey    string
	refreshPapper string
}

func NewAuthService(repo repository.AuthRepository, log *slog.Logger, redis *redis.Client, signingKey, refreshPapper string) *authService {
	return &authService{repo: repo, log: log, signingKey: signingKey, redis: redis, refreshPapper: refreshPapper}
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
	user, err := s.repo.GetUser(name)
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

func (s *authService) SignIn(ctx context.Context, name, password string) (Tokens, error) {
	user, err := s.repo.GetUser(name)
	if err != nil {
		s.log.Error(err.Error())
		return Tokens{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return Tokens{}, ErrInvalidPassword
	}
	refresh, err := s.generateRefreshToken()
	if err != nil {
		s.log.Error(err.Error())
		return Tokens{}, err
	}
	refresh_hash, err := s.generateRefreshHash(refresh)
	if err != nil {
		s.log.Error(err.Error())
		return Tokens{}, err
	}
	expired_at := time.Now().Add((refreshTokenTTL))
	if err := s.repo.SaveRefresh(user.ID, refresh_hash, expired_at); err != nil {
		s.log.Error(err.Error())
		return Tokens{}, err
	}

	s.cacheRefresh(ctx, user.ID, refresh_hash)

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

func (s *authService) generateJWT(user_id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &StandartClaimsWithUserID{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user_id,
	})
	signedToken, err := token.SignedString([]byte(s.signingKey))
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
			return []byte(s.signingKey), nil
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

func (s *authService) RenewalJWT(ctx context.Context, refresh string) (string, error) {
	refresh_hash, err := s.generateRefreshHash(refresh)
	if err != nil {
		s.log.Error(err.Error())
		return "", err
	}
	userID, ok := s.getCachedRefresh(ctx, refresh_hash)
	if !ok {
		userID, err = s.repo.GetUserIDByRefresh(refresh_hash)
		if err != nil {
			return "", err
		}
		s.cacheRefresh(ctx, userID, refresh_hash)
	}
	access, err := s.generateJWT(userID)
	if err != nil {
		s.log.Error(err.Error())
		return "", err
	}
	return access, nil
}

func (s *authService) Logout(ctx context.Context, refresh string) error {
	refresh_hash, err := s.generateRefreshHash(refresh)
	if err != nil {
		return err
	}
	if err := s.repo.RevokeRefreshByHash(refresh_hash); err != nil {
		s.log.Error(err.Error())
		return repository.ErrRefreshTokenNotFound
	}
	s.uncacheRefresh(ctx, refresh_hash)
	return nil
}

func (s *authService) LogoutAll(ctx context.Context, refresh string) error {
	refresh_hash, err := s.generateRefreshHash(refresh)
	if err != nil {
		return nil
	}
	user_id, err := s.repo.GetUserIDByRefresh(refresh_hash)
	if err != nil {
		s.log.Error("not get user_id by refresh")
		return err
	}

	if err := s.repo.RevokeAllRefreshByUserID(user_id); err != nil {
		s.log.Error(err.Error())
		return repository.ErrRefreshTokenNotFound
	}
	s.uncacheAllRefresh(ctx, user_id)
	return nil
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

func (s *authService) generateRefreshHash(refersh string) (string, error) {
	// bytes, err := bcrypt.GenerateFromPassword([]byte(refersh), bcrypt.DefaultCost)
	hash := hmac.New(sha256.New, []byte(s.refreshPapper))
	hash.Write([]byte(refersh))
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func refreshCacheKey(hash string) string {
	return "refresh:" + hash
}

func userSessionsKey(userID int) string {
	return "user_sessions:" + strconv.Itoa(userID)
}

func (s *authService) getCachedRefresh(ctx context.Context, hash string) (int, bool) {
	userID, err := s.redis.Get(ctx, refreshCacheKey(hash)).Int()
	if err != nil {
		return 0, false
	}
	return userID, true
}
func (s *authService) cacheRefresh(ctx context.Context, userID int, hash string) {
	pipe := s.redis.TxPipeline()
	pipe.Set(ctx, refreshCacheKey(hash), userID, refreshTokenTTL)
	pipe.SAdd(ctx, userSessionsKey(userID), hash)
	pipe.Expire(ctx, userSessionsKey(userID), refreshTokenTTL)
	pipe.Exec(ctx)
}

func (s *authService) uncacheRefresh(ctx context.Context, hash string) {
	s.redis.Del(ctx, refreshCacheKey(hash))
}

func (s *authService) uncacheAllRefresh(ctx context.Context, userID int) {
	hashes, err := s.redis.SMembers(ctx, userSessionsKey(userID)).Result()
	if err != nil {
		return
	}
	for _, hash := range hashes {
		s.redis.Del(ctx, refreshCacheKey(hash))
	}
	s.redis.Del(ctx, userSessionsKey(userID))
}
