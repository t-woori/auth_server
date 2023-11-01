package auth_kakao

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"lambda_list/internal/auth_kakao"
	"lambda_list/tools"
)

func HandlerAuthKakao(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	tools.Logger().Info("events: ", zap.Any("events", event))
	res, err := auth_kakao.RegisterStudent(event.QueryStringParameters["code"])
	if err != nil {
		tools.Logger().Error("failed to registered student", zap.Error(err))
		marshal, err := json.Marshal(ErrorResponse{
			Response: Response{
				Status:  500,
				Message: "failed to registered student",
			},
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
	marshal, err := json.Marshal(res)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "failed to marshal response",
		}, err
	}
	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(marshal),
	}
	return response, nil
}
