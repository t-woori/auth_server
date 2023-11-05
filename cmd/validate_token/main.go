package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"lambda_list/api/validate_token"
)

func main() {
	lambda.Start(validate_token.HandlerValidateAccessToken)
}
