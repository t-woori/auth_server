package auth_user

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"lambda_list/infrastructure/auth_user"
	"lambda_list/tools"
	"os"
	"time"
)

type customTokenClaim struct {
	Id uuid.UUID `json:"id"`
	jwt.StandardClaims
}

func ValidateAccessToken(rawAccessToken string) (uuid.UUID, error) {
	secret := tools.CustomTokenSecrets{}
	err := tools.SetSecretsAboutAWS(&secret, os.Getenv("AWS_TOKEN_SECRET_NAME"))
	if err != nil {
		return uuid.Nil, errors.New("ACCESS_SECRET is not set")
	}
	err = setTz("Asia/Seoul")
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "fail set timezone")
	}
	return validateToken(rawAccessToken, secret.AccessTokenSecret)
}

func ValidateRefreshToken(rawRefreshToken string) (uuid.UUID, error) {
	secret := tools.CustomTokenSecrets{}
	err := tools.SetSecretsAboutAWS(&secret, os.Getenv("AWS_TOKEN_SECRET_NAME"))
	if err != nil {
		return uuid.Nil, errors.New("ACCESS_SECRET is not set")
	}
	err = setTz("Asia/Seoul")
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "fail set timezone")
	}
	return validateToken(rawRefreshToken, secret.RefreshTokenSecret)
}

func CreateAccessAndRefreshToken(id uuid.UUID) (string, string, error) {
	secret := tools.CustomTokenSecrets{}
	err := tools.SetSecretsAboutAWS(&secret, os.Getenv("AWS_TOKEN_SECRET_NAME"))
	if err != nil {
		return "", "", err
	}
	err = setTz("Asia/Seoul")
	if err != nil {
		return "", "", errors.Wrap(err, "fail set timezone")
	}
	rawAccessToken, err := createAccessToken(id, secret.AccessTokenSecret)
	if err != nil {
		return "", "", errors.Wrap(err, "fail create access token")
	}
	rawRefreshToken, err := createRefreshToken(id, secret.RefreshTokenSecret)
	return rawAccessToken, rawRefreshToken, err
}

func validateToken(rawToken string, secret string) (uuid.UUID, error) {
	claim := customTokenClaim{}
	_, err := jwt.ParseWithClaims(rawToken, &claim, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			tools.Logger().Error("unexpected signing method", zap.String("method", token.Header["alg"].(string)))
			return nil, errors.New("unexpected signing method")
		}
		if ok := token.Claims.Valid(); ok != nil {
			return nil, errors.New("invalid claims")
		}
		return []byte(secret), nil
	})
	if err != nil {
		tools.Logger().Error("fail parse access token", zap.Error(err))
		return uuid.Nil, errors.Wrap(err, "fail parse access token")
	}
	dbStudentId, err := auth_user.FindByStudentId(claim.Id)
	if err != nil {
		tools.Logger().Info("not found student", zap.String("stuent_id", claim.Id.String()), zap.Error(err))
		return uuid.Nil, err
	}
	tools.Logger().Info("found student", zap.String("student_id", dbStudentId.String()))
	return dbStudentId, nil

}

func createAccessToken(id uuid.UUID, secret string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS512, customTokenClaim{
		Id: id,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Add(1 * time.Millisecond).Unix(),
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix()}})
	tokenString, err := claims.SignedString([]byte(secret))
	result, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	tools.Logger().Info("result", zap.Any("result", result), zap.Error(err))
	if err != nil {
		tools.Logger().Error("fail signed access token", zap.Error(err))
		return "", errors.Wrap(err, "fail signed access token")
	}
	return tokenString, nil
}

func createRefreshToken(id uuid.UUID, secret string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS512, customTokenClaim{
		Id: id,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Add(1 * time.Millisecond).Unix(),
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix()}})
	tokenString, err := claims.SignedString([]byte(secret))
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
