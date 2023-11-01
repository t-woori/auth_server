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
	profile, err := infrastructure.GetUserProfile(tokenInfo.AccessToken)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user profile")
	}
	tools.Logger().Info("validated user", zap.Int("kakaoId", profile.Id),
		zap.Bool("needs", profile.KakaoAccount.NotAccessNickName),
		zap.String("nickname", profile.Properties.Nickname))
	if profile.KakaoAccount.NotAccessNickName {
		return nil, errors.New("failed to get nickname")
	}
	return &Student{
		StudentId:    "id",
		Nickname:     profile.Properties.Nickname,
		AccessToken:  tokenInfo.AccessToken,
		RefreshToken: tokenInfo.RefreshToken,
	}, nil
}
