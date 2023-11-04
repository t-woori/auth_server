package auth_kakao

import "github.com/google/uuid"

type Student struct {
	StudentId    uuid.UUID
	AccessToken  string
	RefreshToken string
	Nickname     string
}
