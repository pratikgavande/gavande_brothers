package auth

import (
	"errors"
	"strings"
	"time"

	"billing-app/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const tokenLifetime = 15 * time.Minute

var ErrInvalidCredentials = errors.New("invalid email or password")

type Service struct {
	db        *gorm.DB
	jwtSecret []byte
}

type Claims struct {
	UserID uint64 `json:"user_id"`
	jwt.RegisteredClaims
}

type LoginResult struct {
	AccessToken string
	ExpiresAt   time.Time
	User        models.User
}

func NewService(db *gorm.DB, jwtSecret []byte) *Service {
	return &Service{db: db, jwtSecret: jwtSecret}
}

func (s *Service) Login(email, password string) (LoginResult, error) {
	var user models.User
	err := s.db.Preload("Role").
		Where("email = ?", strings.ToLower(strings.TrimSpace(email))).
		First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return LoginResult{}, ErrInvalidCredentials
	}
	if err != nil {
		return LoginResult{}, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil || !user.IsActive {
		return LoginResult{}, ErrInvalidCredentials
	}

	now := time.Now()
	expiresAt := now.Add(tokenLifetime)
	claims := Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "billing-app",
			Subject:   "user",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{
		AccessToken: signedToken,
		ExpiresAt:   expiresAt,
		User:        user,
	}, nil
}
