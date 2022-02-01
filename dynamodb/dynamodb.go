package dynamodb

import (
	"lambda-dynamodb-users/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const tableName = "Users"

func SaveUser(user types.User) error {
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

	return user, nil
}

func createDynamoSession() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return dynamodb.New(sess)
}
