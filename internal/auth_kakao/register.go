package auth_kakao

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	infrastructure "lambda_list/infrastructure/auth_kakao"
	"lambda_list/infrastructure/auth_kakao/db"
	"lambda_list/internal/auth_user"
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
	studentDao, err := saveStudent(err, profile, tokenInfo)
	if err != nil {
		return nil, err
	}
	accessToken, err := auth_user.CreateAccessToken(studentDao.StudentId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create access token")
	}
	refreshToken, err := auth_user.CreateRefreshToken(studentDao.StudentId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create refresh token")
	}

	return &Student{
		StudentId:    studentDao.StudentId,
		Nickname:     studentDao.NickName,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func saveStudent(err error, profile *infrastructure.KakaoUserProfile, tokenInfo *infrastructure.AuthResponse) (*db.Student, error) {
	studentDao, err := db.FindByKakaoId(profile.Id)
	studentDao.RefreshToken = tokenInfo.RefreshToken
	studentDao.AccessToken = tokenInfo.AccessToken
	studentDao.NickName = profile.Properties.Nickname
	studentDao.KakaoId = profile.Id
	err = db.SaveUser(studentDao)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save user")
	}
	return studentDao, nil
}
