package token_regenerater

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"lambda_list/internal/auth_user"
	"lambda_list/tools"
	"net/http"
)

func HandlerReGenerateToken(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	tokenRequest := RequestReGenerateToken{}
	err := json.Unmarshal([]byte(event.Body), &tokenRequest)
	if err != nil {
		tools.Logger().Fatal("failed to unmarshal", zap.Any("request", event), zap.Error(errors.Cause(err)))
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "not found token",
		}, nil
	}
	studentId, err := auth_user.ValidateRefreshToken(tokenRequest.RawRefreshToken)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusForbidden,
			Body:       "invalidate refresh token"}, nil
	}
	accessToken, refreshToken, err := auth_user.CreateAccessAndRefreshToken(studentId)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "internal error"}, nil
	}
	marshal, _ := json.Marshal(ResponseReGenerateToken{accessToken, refreshToken})
	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(marshal),
	}
	return response, nil
}
