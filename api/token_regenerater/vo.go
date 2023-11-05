package token_regenerater

type RequestReGenerateToken struct {
	RawRefreshToken string `json:"refresh_token"`
}

type ErrResponseReGenerateToken struct {
	Message string `json:"message"`
}

type ResponseReGenerateToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
