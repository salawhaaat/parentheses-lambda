package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/salawhaaat/parentheses-lambda/pkg/handlers"
)

func main() {
	lambda.Start(handlers.GenerateHandler)
}
