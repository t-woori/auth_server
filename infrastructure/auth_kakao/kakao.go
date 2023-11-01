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
	AUTH_URL     = "https://kauth.kakao.com/oauth/token"
	USER_PROFILE = "https://kapi.kakao.com/v2/user/me"
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
		tools.LoggingHttpResponse(rawResponse, err)
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

func GetUserProfile(accessToken string) (*KakaoUserProfile, error) {
	request, err := http.NewRequest(http.MethodGet, USER_PROFILE, nil)
	if err != nil {
		tools.Logger().Error("failed to create request", zap.Error(err))
		return nil, errors.Wrap(err, "failed to create request")
	}
	request.Header.Add("Authorization", "Bearer "+accessToken)
	request.Header.Add("Content-type", "application/x-www-form-urlencoded")
	response, err := requestHttp(request)
	if err != nil {
		tools.Logger().Error("failed to request", zap.Error(err))
		return nil, errors.Wrap(err, "failed to request")
	}
	defer response.Body.Close()
	tools.LoggingHttpResponse(response, err)
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get user profile")
	}
	userProfile := &KakaoUserProfile{}
	err = marshalingRawResponse(response, userProfile)
	if err != nil {
		tools.Logger().Error("failed to marshaling response", zap.Error(err))
		return nil, errors.Wrap(err, "failed to marshaling response")
	}
	return userProfile, nil
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

func marshalingRawResponse[T AuthResponse | KakaoUserProfile | KakaoInfoResponse](response *http.Response, unmarshalValue *T) error {
	err := json.NewDecoder(response.Body).Decode(&unmarshalValue)
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal response: %s", response.Body)
	}
	return nil
}
