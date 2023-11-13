package auth_kakao

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
	"lambda_list/infrastructure/db"
	"lambda_list/tools"
)

func FindByKakaoId(kakaoId int) (*Student, error) {
	dbConnection, err := db.ConnectRDB()
	if err != nil {
		return nil, err
	}
	defer dbConnection.Close()
	tools.Logger().Info("search student", zap.Int("kakaoId", kakaoId))
	studentDao := Student{}
	err = dbConnection.QueryRow("SELECT kakao_id, nickname, id, access_token,refresh_token FROM student WHERE kakao_id = ?", kakaoId).Scan(
		&studentDao.KakaoId, &studentDao.NickName, &studentDao.StudentId, &studentDao.AccessToken, &studentDao.RefreshToken)
	if err != nil {
		tools.Logger().Error("failed to find student", zap.Error(err))
		return &Student{}, err
	}
	tools.Logger().Info("found student", zap.Int("kakaoId", studentDao.KakaoId), zap.String("student name", studentDao.NickName))
	return &studentDao, nil
}

func SaveUser(studentDao *Student) error {
	dbConnection, err := db.ConnectRDB()
	if err != nil {
		return err
	}
	defer dbConnection.Close()
	if studentDao.StudentId != uuid.Nil {
		rawUUID, err := studentDao.StudentId.MarshalBinary()
		tools.Logger().Info("updated student", zap.Int("kakaoId", studentDao.KakaoId), zap.String("student name", studentDao.NickName))
		_, err = dbConnection.Exec("UPDATE student SET access_token = ?, refresh_token = ?, updated_at=now() WHERE id = ? ",
			studentDao.AccessToken, studentDao.RefreshToken, rawUUID)
		return err
	}
	studentDao.StudentId = uuid.New()
	rawUUID, err := studentDao.StudentId.MarshalBinary()
	tools.Logger().Info("insert student", zap.Int("kakaoId", studentDao.KakaoId), zap.String("student name", studentDao.NickName))
	_, err = dbConnection.Exec("INSERT INTO student (kakao_id, nickname, id, access_token, refresh_token,created_at,updated_at) VALUES (?, ?, ?, ?, ?,now(),now())",
		studentDao.KakaoId, studentDao.NickName, rawUUID, studentDao.AccessToken, studentDao.RefreshToken)
	return err
}
