package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"lambda_list/tools"
	"os"
	"strconv"
)

func ConnectRDB() (*sql.DB, error) {
	dbURL, err := getDbURL()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get db url")
	}
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		tools.Logger().Error("fail connect db", zap.Error(err))
		return nil, err
	}
	return db, nil
}

func getDbURL() (string, error) {
	rdsSecret := tools.RdsSecrets{}
	err := tools.SetSecretsAboutAWS(&rdsSecret, os.Getenv("AWS_RDS_SECRET_NAME"))
	if err != nil {
		return "", errors.Wrap(err, "failed to get secret value")
	}
	user := rdsSecret.Username
	password := rdsSecret.Password
	host := os.Getenv("RDS_PROXY_HOST")
	port := rdsSecret.Port
	dbName := "db"
	return user + ":" + password + "@tcp(" + host + ":" + strconv.Itoa(port) + ")/" + dbName, nil
}
