package jwt

import (
	crand "crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

const (
	randMin = 100000
	randMax = 999999
)

type Service struct {
	secret            string
	accessExpiration  time.Duration
	refreshExpiration time.Duration
}

func NewService(secret string, accessExpiration, refreshExpiration time.Duration) *Service {
	return &Service{
		secret:            secret,
		accessExpiration:  accessExpiration,
		refreshExpiration: refreshExpiration,
	}
}

func (s *Service) GenerateToken(userID string) (string, error) {
	now := time.Now()
	claims := &AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessExpiration)),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *Service) GenerateRefreshToken() (string, error) {
	token := make([]byte, 32)
	if _, err := crand.Read(token); err != nil {
		return "", fmt.Errorf("failed to generate refresh token %w", err)
	}

	return base64.URLEncoding.EncodeToString(token), nil
}

func (s *Service) GenerateCode() int {
	return rand.Intn(randMax-randMin) + randMin
}

func (s *Service) GetAccessExpiration() time.Duration {
	return s.accessExpiration
}

func (s *Service) GetRefreshExpiration() time.Duration {
	return s.refreshExpiration
}
