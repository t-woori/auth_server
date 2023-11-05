package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"lambda_list/api/token_regenerater"
)

func main() {
	lambda.Start(token_regenerater.HandlerReGenerateToken)
}
