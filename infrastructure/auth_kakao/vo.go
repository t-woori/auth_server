package auth_kakao

type AuthRequest struct {
	ClientId    string `json:"client_id"`
	RedirectUri string `json:"redirect_uri"`
	Code        string `json:"code"`
	GrantType   string `json:"grant_type"`
}

type KakaoInfoResponse struct {
	KakaoId   int `json:"id"`
	ExpiresIn int `json:"expires_in"`
	AppId     int `json:"app_id"`
}

type KakaoUserProfile struct {
	Id           int    `json:"id"`
	ConnectedAt  string `json:"connected_at"`
	KakaoAccount struct {
		NotAccessNickName bool `json:"profile_nickname_needs_agreement"`
	} `json:"kakao_account"`
	Properties struct {
		Nickname string `json:"nickname"`
	} `json:"properties"`
}
