package infrastructure

import (
	"encoding/json"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
	"lambda_list/tools"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	AUTH_URL            = "https://kauth.kakao.com/oauth/token"
	VALIDATE_TOPKEN_URL = "https://kapi.kakao.com/v1/user/access_token_info"
)

var clientKey = os.Getenv("KAKAO_CLIENT_KEY")
var redirectUri = os.Getenv("KAKAO_REDIRECT_URI")

func GenerateToken(authCode string) (*AuthResponse, error) {
	rawResponse, err := getTokenByKakao(authCode)
	if err != nil {
		tools.Logger().Error("failed to get token. caused by", zap.Error(err))
		return nil, errors.Wrapf(err, "failed to get token from code: %s", authCode)
	}
	if rawResponse.StatusCode != http.StatusOK {
		loggingHttpResponse(rawResponse, err)
		return nil, errors.Wrapf(err, "failed to get token from code: %s", authCode)
	}
	defer rawResponse.Body.Close()
	response := &AuthResponse{}
	err = marshalingRawResponse(rawResponse, response)
	if err != nil {
		tools.Logger().Error("failed to marshaling response", zap.Error(err))
		return nil, err
	}
	tools.Logger().Info("access token: ", zap.Any("token", response.AccessToken))
	return response, err
}

func ValidateToken(accessToken string) (*KakaoInfoResponse, error) {
	request, err := http.NewRequest(http.MethodGet, VALIDATE_TOPKEN_URL, nil)
	request.Header.Add("Authorization", "Bearer "+accessToken)
	tools.Logger().Info("request: ", zap.Strings("request header", request.Header.Values("Authorization")),
		zap.String("url", request.URL.String()),
		zap.String("method", request.Method))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}
	response, err := requestHttp(request)
	defer response.Body.Close()
	if err != nil {
		return nil, errors.Wrap(err, "failed to request")
	}
	if response.StatusCode != http.StatusOK {
		loggingHttpResponse(response, err)
		return nil, errors.New("failed to validate token")
	}
	responseValidateToken := &KakaoInfoResponse{}
	err = marshalingRawResponse(response, responseValidateToken)
	if err != nil {
		tools.Logger().Error("failed to marshaling response", zap.Error(err))
		return nil, errors.Wrap(err, "failed to marshaling response")
	}
	tools.Logger().Info("token_info: ", zap.Any("token_info", responseValidateToken))
	return responseValidateToken, nil
}

func loggingHttpResponse(rawResponse *http.Response, err error) {
	tools.Logger().Error("failed get token",
		zap.Int("status", rawResponse.StatusCode),
		zap.Any("header", rawResponse.Header),
		zap.String("body", func() string {
			stringBuffer := &strings.Builder{}
			io.Copy(stringBuffer, rawResponse.Body)
			return stringBuffer.String()
		}()), zap.Error(err))
}

func getTokenByKakao(authCode string) (*http.Response, error) {
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

func requestHttp(request *http.Request) (*http.Response, error) {
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		tools.Logger().Error("failed to request", zap.Error(err))
	}
	return response, nil
}

func marshalingRawResponse[T AuthResponse | KakaoInfoResponse](response *http.Response, unmarshalValue *T) error {
	tools.Logger().Info("response kakao: ",
		zap.Int("status", response.StatusCode),
		zap.Any("header", response.Header),
		zap.Any("body", response.Body))
	err := json.NewDecoder(response.Body).Decode(&unmarshalValue)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal response: %s", response.Body)
	}
	return nil
}
