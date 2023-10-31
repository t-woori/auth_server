package infrastructure

import (
	"encoding/json"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"lambda_list/tools"
	"net/http"
	"net/url"
	"os"
)

const (
	AUTH_URL = "https://kauth.kakao.com/oauth/token"
)

var clientKey = os.Getenv("KAKAO_CLIENT_KEY")
var redirectUri = os.Getenv("KAKAO_REDIRECT_URI")

func GenerateToken(authCode string) (*AuthResponse, error) {
	rawResponse, err := requestKakao(authCode)
	if err != nil {
		tools.Logger().Error("failed to get token. caused by", zap.Error(err))
		return nil, errors.Wrapf(err, "failed to get token from code: %s", authCode)
	}
	tools.Logger().Debug("response: ", zap.Int("status", rawResponse.StatusCode))
	if rawResponse.StatusCode != http.StatusOK {
		tools.Logger().Error("failed get token",
			zap.Int("status", rawResponse.StatusCode),
			zap.Any("header", rawResponse.Header),
			zap.Any("body", rawResponse.Body), zap.Error(err))
		return nil, errors.Wrapf(err, "failed to get token from code: %s", authCode)
	}
	defer rawResponse.Body.Close()
	var response AuthResponse
	err = marshalingRawResponse(rawResponse, response)
	tools.Logger().Debug("response: ", zap.Any("response", response))
	if err != nil {
		return nil, err
	}
	return &response, err
}

func requestKakao(authCode string) (*http.Response, error) {
	rawResponse, err := http.PostForm(AUTH_URL, url.Values{
		"client_id":    {clientKey},
		"redirect_uri": {redirectUri},
		"code":         {authCode},
		"grant_type":   {"authorization_code"},
	})

	if err != nil {
		tools.Logger().Error("failed to get token from code", zap.Error(err))
		return nil, errors.Wrapf(err, "failed to get token from code: %s", authCode)
	}
	return rawResponse, nil
}

func marshalingRawResponse[T AuthResponse | ErrorResponse](response *http.Response, unmarshalValue T) error {
	err := json.NewDecoder(response.Body).Decode(&unmarshalValue)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal response: %s", response.Body)
	}
	tools.Logger().Info("response: ", zap.Any("response", unmarshalValue))
	return nil
}
