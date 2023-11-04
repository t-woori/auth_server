package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"lambda_list/api/auth_kakao"
)

func main() {
	lambda.Start(auth_kakao.HandlerAuthKakao)
}
