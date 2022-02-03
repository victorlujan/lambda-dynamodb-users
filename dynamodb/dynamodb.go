package dynamodb

import (
	"lambda-dynamodb-users/types"

	errs "github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const tableName = "Users"

func SaveUser(user types.User) error {
	exitst, _ := GetUser(user.ID)

	if exitst.ID == "" {
		userMap, err := dynamodbattribute.MarshalMap(user)
		if err != nil {
			return err
		}
		dynamodbSession := createDynamoSession()

		input := &dynamodb.PutItemInput{
			Item:      userMap,
			TableName: aws.String(tableName),
		}
		_, err = dynamodbSession.PutItem(input)
		if err != nil {
			return err
		}
	}

	return nil

}

func GetUser(id string) (types.User, error) {
	dynamodbSession := createDynamoSession()

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(tableName),
	}
	result, err := dynamodbSession.GetItem(input)
	if err != nil {
		return types.User{}, err
	}

	var user types.User
	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return types.User{}, err
	}
	if user.ID == "" {
		return types.User{}, errs.New("User not found")
	}

	return user, nil
}

func QueryUser(email string) (types.User, error) {
	dynamodbSession := createDynamoSession()

	input := &dynamodb.QueryInput{
		TableName: aws.String(tableName),
		IndexName: aws.String("email-index"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":email": {
				S: aws.String(email),
			},
		},
		KeyConditionExpression: aws.String("email = :email"),
	}
	result, err := dynamodbSession.Query(input)
	if err != nil {
		return types.User{}, err
	}

	var users []types.User
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return types.User{}, err
	}

	if len(users) == 0 {
		return types.User{}, errs.New("User not found")
	}

	return users[0], nil
}




func DeleteUser(id string) error {
	dynamodbSession := createDynamoSession()

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(tableName),
	}
	_, err := dynamodbSession.DeleteItem(input)
	if err != nil {
		return err
	}

	return nil
}

func ScanUsers(user types.User) ([]types.User, error) {
	dynamodbSession := createDynamoSession()

	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":id": {
				S: aws.String(user.ID),
			},
		},
		FilterExpression: aws.String("id = :id"),
	}
	result, err := dynamodbSession.Scan(input)
	if err != nil {
		return []types.User{}, err
	}

	var users []types.User
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return []types.User{}, err
	}

	return users, nil
}

func createDynamoSession() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return dynamodb.New(sess)
}
