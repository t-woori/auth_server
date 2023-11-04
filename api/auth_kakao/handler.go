package auth_kakao

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	customErrors "lambda_list/infrastructure/auth_kakao"
	"lambda_list/internal/auth_kakao"
	"lambda_list/tools"
	"net/http"
)

func HandlerAuthKakao(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	kakaoTokens := AuthKakaoRequestBody{}
	err := json.Unmarshal([]byte(event.Body), &kakaoTokens)
	if err != nil {
		tools.Logger().Fatal("failed to unmarshal", zap.Any("request", event), zap.Error(errors.Cause(err)))
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "not found token",
		}, nil
	}
	tools.Logger().Info("kakao tokens", zap.Any("tokens", kakaoTokens))
	res, err := auth_kakao.RegisterStudent(auth_kakao.KakaoTokens{
		AccessToken:  kakaoTokens.AccessToken,
		RefreshToken: kakaoTokens.RefreshToken})
	if err != nil {
		tools.Logger().Error("failed to registered student", zap.Error(errors.Cause(err)))
		if errors.Is(err, customErrors.ErrFailedToGetToken) {
			marshal, _ := json.Marshal(Response{
				Status:  http.StatusBadRequest,
				Message: "can't access nickname"})
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       string(marshal),
			}, nil
		}
		marshal, _ := json.Marshal(Response{
			Status:  500,
			Message: "internal error"})
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       string(marshal),
		}, nil
	}
	marshal, _ := json.Marshal(res)
	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(marshal),
	}
	return response, nil
}
