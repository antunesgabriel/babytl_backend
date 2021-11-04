package configs

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTConfig struct {
	secret string
	issure string
}

func NewJWT() *JWTConfig {
	return &JWTConfig{
		os.Getenv("APP_SECRET"),
		"babytl_backend",
	}
}

type Clain struct {
	Sum uint `json:"sum"`
	jwt.StandardClaims
}

func (j *JWTConfig) GenerateToken(id uint) (string, error) {
	clain := &Clain{
		id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
			Issuer: j.issure,
			IssuedAt: time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, clain)

	t, err := token.SignedString([]byte(j.secret))

	if err != nil {
		return "", err
	}

	return t, nil
}

func (j *JWTConfig) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, isValid := t.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("invalid token: %s", token)
		}

		return []byte(j.secret), nil
	})

}