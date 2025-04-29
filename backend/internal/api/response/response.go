package response

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

var (
	responseHeaders = map[string]string{
		"Content-Type": "application/json",
	}
)

func Success(data any) (events.APIGatewayProxyResponse, error) {
	responseBytes, err := json.Marshal(data)
	if err != nil {
		return Error(http.StatusInternalServerError, "Failed to serialize response data")
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    responseHeaders,
		Body:       string(responseBytes),
	}, nil
}

func Error(statusCode int, message string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    responseHeaders,
		Body:       fmt.Sprintf(`{"message":"%s"}`, message),
	}, nil
}
