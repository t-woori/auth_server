package auth_kakao

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	infrastructure "lambda_list/infrastructure/auth_kakao"
	"lambda_list/internal/auth_user"
	"lambda_list/tools"
)

func RegisterStudent(kakaoTokens KakaoTokens) (*Student, error) {
	profile, err := infrastructure.GetUserProfile(kakaoTokens.AccessToken)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user profile")
	}
	tools.Logger().Info("validated user", zap.Int("kakaoId", profile.Id),
		zap.Bool("needs", profile.KakaoAccount.NotAccessNickName),
		zap.String("nickname", profile.Properties.Nickname))
	if profile.KakaoAccount.NotAccessNickName {
		return nil, errors.New("failed to get nickname")
	}
	studentDao, err := saveStudent(err, profile, kakaoTokens)
	if err != nil {
		return nil, err
	}
	accessToken, refreshToken, err := auth_user.CreateAccessAndRefreshToken(studentDao.StudentId)
	return &Student{
		StudentId:    studentDao.StudentId,
		Nickname:     studentDao.NickName,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func saveStudent(err error, profile *infrastructure.KakaoUserProfile, tokens KakaoTokens) (*infrastructure.Student, error) {
	studentDao, err := infrastructure.FindByKakaoId(profile.Id)
	studentDao.RefreshToken = tokens.RefreshToken
	studentDao.AccessToken = tokens.AccessToken
	studentDao.NickName = profile.Properties.Nickname
	studentDao.KakaoId = profile.Id
	err = infrastructure.SaveUser(studentDao)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save user")
	}
	return studentDao, nil
}
