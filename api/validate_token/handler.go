package validate_token

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"lambda_list/internal/auth_user"
	"lambda_list/tools"
	"strings"
)

func parsingRawJwtToken(authorizationHeader string) (string, error) {
	if !strings.HasPrefix(authorizationHeader, "Bearer ") {
		return "", errors.New("authorization header is empty")
	}
	authHeaderParts := strings.Split(authorizationHeader, "Bearer ")
	if len(authHeaderParts) != 2 {
		return "", errors.New("authorization header is empty")
	}
	return authHeaderParts[1], nil
}

func HandlerValidateAccessToken(ctx context.Context, event events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {
	authorizationHeader := event.Headers["authorization"]
	rawJwtToken, err := parsingRawJwtToken(authorizationHeader)
	if err != nil {
		tools.Logger().Error("fail parsing raw jwt token", zap.Error(err))
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false}, nil
	}
	studentId, err := auth_user.ValidateAccessToken(rawJwtToken)
	if err != nil {
		tools.Logger().Error("fail validate student", zap.Error(err))
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false}, nil
	}
	tools.Logger().Info("validated student", zap.String("student_id", studentId.String()))
	authData := make(map[string]any)
	authData["student_id"] = studentId.String()
	tools.Logger().Info("data", zap.Any("data", authData))
	return events.APIGatewayV2CustomAuthorizerSimpleResponse{
		IsAuthorized: true,
		Context:      authData}, nil
}
