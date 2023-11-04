package auth_kakao

import "net/http"

type Response struct {
	Status  http.ConnState `json:"status"`
	Message string         `json:"message"`
}

type StudentResponse struct {
	Response
	StudentId    string `json:"student_id"`
	Nickname     string `json:"nickname"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthKakaoRequestBody struct {
	AccessToken  string `json:"kakao_access_token"`
	RefreshToken string `json:"kakao_refresh_token"`
}
