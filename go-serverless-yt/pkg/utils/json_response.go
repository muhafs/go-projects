package utils

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

func ApiResponse(data interface{}, code int) (*events.APIGatewayProxyResponse, error) {
	r := events.APIGatewayProxyResponse{Headers: map[string]string{"Content-Type": "application/json"}}
	r.StatusCode = code

	stringBody, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Couldn't parse data from json:", err)
	}
	r.Body = string(stringBody)

	return &r, nil
}
