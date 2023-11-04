package auth_user

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"lambda_list/tools"
	"os"
	"time"
)

func CreateAccessToken(id uuid.UUID) (string, error) {
	secretValue := os.Getenv("ACCESS_SECRET")
	if secretValue == "" {
		return "", errors.New("ACCESS_SECRET is not set")
	}
	secretKey := []byte(secretValue)
	err := setTz("Asia/Seoul")
	if err != nil {
		return "", errors.Wrap(err, "fail set timezone")
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(15 * time.Minute).Unix()})

	// Create a new token with the claims
	tokenString, err := claims.SignedString(secretKey)
	if err != nil {
		tools.Logger().Error("fail signed access token", zap.Error(err))
		return "", errors.Wrap(err, "fail signed access token")
	}
	return tokenString, nil
}

func CreateRefreshToken(id uuid.UUID) (string, error) {
	secretValue := os.Getenv("REFRESH_SECRET")
	if secretValue == "" {
		return "", errors.New("REFRESH_SECRET is not set")
	}
	secretKey := []byte(secretValue)
	err := setTz("Asia/Seoul")
	if err != nil {
		return "", errors.Wrap(err, "fail set timezone")
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(24 * time.Hour).Unix()})
	// Create a new token with the claims
	tokenString, err := claims.SignedString(secretKey)
	if err != nil {
		tools.Logger().Error("fail signed refresh token", zap.Error(err))
		return "", err
	}
	return tokenString, nil
}

func setTz(nameOfTz string) error {
	location, err := time.LoadLocation(nameOfTz)
	if err != nil {
		return err
	}
	time.Local = location
	return nil
}
