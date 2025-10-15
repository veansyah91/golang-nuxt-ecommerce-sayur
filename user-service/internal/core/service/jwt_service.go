package service

import (
	"time"
	"user-service/config"

	"github.com/golang-jwt/jwt/v5"
)

type JwtServiceInterface interface {
	GenerateToken(userId int64) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtService struct {
	secretKey string
	issuer    string
}

func (j *jwtService) GenerateToken(userId int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"iss":     j.issuer,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.secretKey))
}

func (j *jwtService) ValidateToken(encode_token string) (*jwt.Token, error) {
	return jwt.Parse(encode_token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(j.secretKey), nil
	})
}

func NewJwtService(cfg *config.Config) JwtServiceInterface {
	return &jwtService{
		secretKey: cfg.App.JwtSecretKey,
		issuer:    cfg.App.JwtSecretKey,
	}
}
