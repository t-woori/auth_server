package tools

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"os"
)

func SetSecretsAboutAWS[T Secret](secret *T, secretName string) error {
	Logger().Info("get secrets on AWS Secret Manager", zap.String("secret name", secretName))
	svc, err := createAWSSecretManagerClient()
	if err != nil {
		Logger().Fatal("failed to create aws secret manager client", zap.Error(err))
		return errors.Wrap(err, "failed to create aws secret manager client")
	}
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"),
	}
	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		Logger().Fatal("failed to get secret value", zap.String("AWS_SECRET_NAME", os.Getenv("AWS_SECRET_NAME")), zap.Error(err))
		return errors.Wrap(err, "failed to get secret value")
	}
	err = json.Unmarshal([]byte(*result.SecretString), &secret)
	return nil
}

func createAWSSecretManagerClient() (*secretsmanager.Client, error) {
	region := os.Getenv("AWS_REGION")
	Logger().Info("create aws secret manager client", zap.String("region", region))
	secretConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		Logger().Error("failed to load config", zap.Error(err))
		return nil, err
	}
	return secretsmanager.NewFromConfig(secretConfig), nil
}
