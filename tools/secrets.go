package tools

type Secret interface {
	RdsSecrets | CustomTokenSecrets
}

type RdsSecrets struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Port     int    `json:"port"`
}

type CustomTokenSecrets struct {
	AccessTokenSecret  string `json:"access_token_secret"`
	RefreshTokenSecret string `json:"refresh_token_secret"`
}
