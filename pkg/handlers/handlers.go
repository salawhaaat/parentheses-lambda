package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/salawhaaat/parentheses-lambda/pkg/parentheses"
)

// Response structure for returning JSON
type Response struct {
	StatusCode int         `json:"statusCode"`
	Body       interface{} `json:"body"`
}

// GenerateHandler is the Lambda function handler
func GenerateHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Extract and validate the 'length' parameter from the query string
	lengthParam := request.QueryStringParameters["n"]
	length, err := validateLengthParam(lengthParam)
	if err != nil {
		return createErrorResponse(http.StatusBadRequest, "Invalid parameter 'n'. It must be a positive integer.")
	}
	sequence := parentheses.Generate(length)
	return createSuccessResponse(sequence)
}

// validateLengthParam converts and validates the 'length' query parameter.
func validateLengthParam(lengthParam string) (int, error) {
	length, err := strconv.Atoi(lengthParam)
	if err != nil || length <= 0 {
		return 0, err
	}
	return length, nil
}

// createErrorResponse generates an API Gateway Proxy response for error scenarios.
func createErrorResponse(statusCode int, message string) (events.APIGatewayProxyResponse, error) {
	response := Response{
		StatusCode: statusCode,
		Body:       message,
	}
	responseBody, err := json.Marshal(response.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Failed to generate the response",
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: response.StatusCode,
		Body:       string(responseBody),
	}, nil
}

// createSuccessResponse generates an API Gateway Proxy response with the generated sequence.
func createSuccessResponse(sequence string) (events.APIGatewayProxyResponse, error) {
	response := Response{
		StatusCode: http.StatusOK,
		Body:       map[string]string{"sequence": sequence},
	}
	responseBody, err := json.Marshal(response.Body)
	if err != nil {
		return createErrorResponse(http.StatusInternalServerError, "Failed to generate the response")
	}
	return events.APIGatewayProxyResponse{
		StatusCode: response.StatusCode,
		Body:       string(responseBody),
	}, nil
}
