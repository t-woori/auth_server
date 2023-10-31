package auth_kakao

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	infrastructure "lambda_list/infrastructure/auth_kakao"
	"lambda_list/tools"
)

func HandlerAuthKakao(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	tools.Logger().Info("events: ", zap.Any("events", event))
	res, err := infrastructure.GenerateToken(event.QueryStringParameters["code"])
	if err != nil {
		tools.Logger().Error("failed to get token from code", zap.Error(err))
		marshal, err := json.Marshal(ErrorResponse{
			Error: errors.Unwrap(err).Error(),
		})
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       string(marshal),
		}, nil
	}
	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       res.AccessToken,
	}
	return response, nil
}
