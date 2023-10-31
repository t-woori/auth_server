package infrastructure

type AuthRequest struct {
	ClientId    string `json:"client_id"`
	RedirectUri string `json:"redirect_uri"`
	Code        string `json:"code"`
	GrantType   string `json:"grant_type"`
}
type AuthResponse struct {
	TokenType             string `json:"token_type"`
	AccessToken           string `json:"access_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
}

type ErrorResponse struct {
	ErrorCode        string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
