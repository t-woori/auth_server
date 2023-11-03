package auth_kakao

import "github.com/google/uuid"

type Student struct {
	StudentId    uuid.UUID
	KakaoId      int
	AccessToken  string
	RefreshToken string
	Nickname     string
}
