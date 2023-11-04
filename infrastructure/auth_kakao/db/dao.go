package db

import "github.com/google/uuid"

type Student struct {
	KakaoId      int       `json:"kakao_id"`
	NickName     string    `json:"nickname"`
	StudentId    uuid.UUID `json:"student_id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}
