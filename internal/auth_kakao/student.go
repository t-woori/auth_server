package auth_kakao

import "github.com/google/uuid"

type Student struct {
	StudentId    uuid.UUID `json:"student_id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Nickname     string    `json:"nickname"`
}
