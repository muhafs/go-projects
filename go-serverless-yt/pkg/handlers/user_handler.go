package handlers

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/muhafs/go-serverless-yt/pkg/models/user"
	"github.com/muhafs/go-serverless-yt/pkg/utils"
)

const MethodNotAllowed = "method not allowed"

type ErrorBody struct {
	Message *string `json:"error,omiempty"`
}

func GetUser(r events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	email := r.QueryStringParameters["email"]
	if len(email) > 0 {
		user, err := user.FetchUser(email, tableName, dynaClient)
		if err != nil {
			return utils.ApiResponse(
				ErrorBody{
					Message: aws.String(err.Error()),
				},
				http.StatusBadRequest,
			)
		}

		return utils.ApiResponse(user, http.StatusOK)
	}

	users, err := user.ListUsers(tableName, dynaClient)
	if err != nil {
		return utils.ApiResponse(
			ErrorBody{
				Message: aws.String(err.Error()),
			},
			http.StatusBadRequest,
		)
	}

	return utils.ApiResponse(users, http.StatusOK)
}

func CreateUser(r events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	user, err := user.CreateUser(r, tableName, dynaClient)
	if err != nil {
		return utils.ApiResponse(
			ErrorBody{
				Message: aws.String(err.Error()),
			},
			http.StatusBadRequest,
		)
	}

	return utils.ApiResponse(user, http.StatusOK)
}

func UpdateUser(r events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	user, err := user.UpdateUser(r, tableName, dynaClient)
	if err != nil {
		return utils.ApiResponse(
			ErrorBody{
				Message: aws.String(err.Error()),
			},
			http.StatusBadRequest,
		)
	}

	return utils.ApiResponse(user, http.StatusOK)
}
func DeleteUser(r events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	err := user.DeleteUser(r, tableName, dynaClient)
	if err != nil {
		return utils.ApiResponse(
			ErrorBody{
				Message: aws.String(err.Error()),
			},
			http.StatusBadRequest,
		)
	}

	return utils.ApiResponse(nil, http.StatusOK)
}

func UnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return utils.ApiResponse(MethodNotAllowed, http.StatusMethodNotAllowed)
}
