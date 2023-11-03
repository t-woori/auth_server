package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"lambda_list/tools"
	"os"
	"strconv"
)

type _RdsSecrets struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Port     int    `json:"port"`
}

func FindByKakaoId(kakaoId int) (*Student, error) {
	db, err := connectRDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	tools.Logger().Info("search student", zap.Int("kakaoId", kakaoId))
	dao := Student{}
	err = db.QueryRow("SELECT kakao_id, nickname, student_id, access_token,refresh_token FROM students WHERE kakao_id = ?", kakaoId).Scan(
		&dao.KakaoId, &dao.NickName, &dao.StudentId, &dao.AccessToken, &dao.RefreshToken)
	tools.Logger().Info("found student", zap.Any("student", dao), zap.Error(err))
	if err != nil {
		tools.Logger().Error("failed to find student", zap.Error(err))
		return &Student{}, err
	}
	tools.Logger().Info("found student", zap.Any("student", dao))
	return &dao, nil
}

func SaveUser(studentDao *Student) error {
	db, err := connectRDB()
	if err != nil {
		return err
	}
	defer db.Close()
	if studentDao.StudentId != uuid.Nil {
		tools.Logger().Info("updated student", zap.Any("student", studentDao))
		_, err = db.Exec("UPDATE students SET access_token = ?, refresh_token = ?, updated_at=now() WHERE student_id = ? ",
			studentDao.AccessToken, studentDao.RefreshToken, studentDao.StudentId)
		return err
	}
	studentDao.StudentId = uuid.New()
	tools.Logger().Info("inserted student", zap.Any("student", studentDao))
	_, err = db.Exec("INSERT INTO students (kakao_id, nickname, student_id, access_token, refresh_token,created_at,updated_at) VALUES (?, ?, ?, ?, ?,now(),now())",
		studentDao.KakaoId, studentDao.NickName, studentDao.StudentId, studentDao.AccessToken, studentDao.RefreshToken)
	return err
}

func connectRDB() (*sql.DB, error) {
	dbURL, err := getDbURL()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get db url")
	}
	db, err := sql.Open("mysql", dbURL)
	tools.Logger().Info("connected db", zap.Any("db", db.Stats()))
	if err != nil {
		tools.Logger().Error("fail connect db", zap.Error(err))
		return nil, err
	}
	tools.Logger().Info("connected db", zap.Any("db", db.Stats()))
	return db, nil
}

func getDbURL() (string, error) {
	secrets, err := getValueAboutAWSSecret()
	if err != nil {
		return "", errors.Wrap(err, "failed to get secret value")
	}
	user := secrets.Username
	password := secrets.Password
	host := os.Getenv("RDS_PROXY_HOST")
	port := secrets.Port
	dbName := "db"
	return user + ":" + password + "@tcp(" + host + ":" + strconv.Itoa(port) + ")/" + dbName, nil
}

func getValueAboutAWSSecret() (_RdsSecrets, error) {
	tools.Logger().Info("get secrets on AWS Secret Manager", zap.String("AWS_SECRET_NAME", os.Getenv("AWS_SECRET_NAME")))
	svc, err := createAWSSecretManagerClient()
	if err != nil {
		tools.Logger().Fatal("failed to create aws secret manager client", zap.Error(err))
		return _RdsSecrets{}, errors.Wrap(err, "failed to create aws secret manager client")
	}
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(os.Getenv("AWS_SECRET_NAME")),
		VersionStage: aws.String("AWSCURRENT"),
	}
	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		tools.Logger().Fatal("failed to get secret value", zap.String("AWS_SECRET_NAME", os.Getenv("AWS_SECRET_NAME")), zap.Error(err))
		return _RdsSecrets{}, errors.Wrap(err, "failed to get secret value")
	}
	rdsSecrets := _RdsSecrets{}
	err = json.Unmarshal([]byte(*result.SecretString), &rdsSecrets)
	return rdsSecrets, err
}

func createAWSSecretManagerClient() (*secretsmanager.Client, error) {
	region := os.Getenv("AWS_REGION")
	tools.Logger().Info("create aws secret manager client", zap.String("region", region))
	secretConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		tools.Logger().Error("failed to load config", zap.Error(err))
		return nil, err
	}
	return secretsmanager.NewFromConfig(secretConfig), nil
}
