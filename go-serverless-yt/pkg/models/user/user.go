package user

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/muhafs/go-serverless-yt/pkg/validators"
)

const (
	FailedToFetchRecord     = "failed to fetch record"
	FailedToUnmarshalRecord = "failed to unmarshal record"
	InvalidUserData         = "invalid user data"
	InvalidEmail            = "invalid email"
	CouldNotMarshalItem     = "couldn't marshal item"
	CouldNotDeleteItem      = "couldn't delete item"
	CouldNotDynamoPutItem   = "couldn't dynamo put item"
	UserAlreadyExists       = "user already exists"
	UserNotFound            = "user not found"
)

type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func ListUsers(tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*[]User, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.Scan(input)
	if err != nil {
		return nil, errors.New(FailedToFetchRecord)
	}

	items := new([]User)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, items)
	if err != nil {
		return nil, errors.New(FailedToUnmarshalRecord)
	}

	return items, nil
}

func FetchUser(email string, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*User, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {S: aws.String(email)},
		},
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.GetItem(input)
	if err != nil {
		return nil, errors.New(FailedToFetchRecord)
	}

	item := new(User)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, errors.New(FailedToUnmarshalRecord)
	}

	return item, nil
}
func CreateUser(r events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*User, error) {
	var u User
	if err := json.Unmarshal([]byte(r.Body), &u); err != nil {
		return nil, errors.New(InvalidUserData)
	}

	if !validators.IsEmail(u.Email) {
		return nil, errors.New(InvalidEmail)
	}

	currUser, _ := FetchUser(u.Email, tableName, dynaClient)
	if currUser != nil && len(currUser.Email) != 0 {
		return nil, errors.New(UserAlreadyExists)
	}

	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, errors.New(CouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(CouldNotDynamoPutItem)
	}

	return &u, nil
}

func UpdateUser(r events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*User, error) {
	var u User
	if err := json.Unmarshal([]byte(r.Body), &u); err != nil {
		return nil, errors.New(InvalidUserData)
	}

	if !validators.IsEmail(u.Email) {
		return nil, errors.New(InvalidEmail)
	}

	currUser, _ := FetchUser(u.Email, tableName, dynaClient)
	if currUser != nil && len(currUser.Email) == 0 {
		return nil, errors.New(UserNotFound)
	}

	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, errors.New(CouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(CouldNotDynamoPutItem)
	}

	return &u, nil
}

func DeleteUser(r events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) error {
	email := r.QueryStringParameters["email"]

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := dynaClient.DeleteItem(input)
	if err != nil {
		return errors.New(CouldNotDeleteItem)
	}

	return nil
}
