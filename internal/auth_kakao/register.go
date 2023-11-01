package auth_kakao

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	infrastructure "lambda_list/infrastructure/auth_kakao"
	"lambda_list/tools"
)

func RegisterStudent(authCode string) (*Student, error) {
	tokenInfo, err := infrastructure.GenerateToken(authCode)
	if err != nil {
		return nil, err
	}
	tools.Logger().Info("token info: ", zap.Any("token", tokenInfo))
	validatedInfo, err := infrastructure.ValidateToken(tokenInfo.AccessToken)
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate token")
	}
	tools.Logger().Info("validated user", zap.Int("kakaoId", validatedInfo.KakaoId))
	return &Student{
		StudentId:    "id",
		AccessToken:  tokenInfo.AccessToken,
		RefreshToken: tokenInfo.RefreshToken,
	}, nil
}
