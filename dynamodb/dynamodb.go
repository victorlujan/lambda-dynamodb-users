package dynamodb

import(
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"lambda-dynamodb-users/types"
)

const tableName = "Users"

func SaveUser(user types.User) error{
	userMap, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}

	dynamodbSession:= createDynamoSession()

	input:= &dynamodb.PutItemInput{
		Item: userMap,
		TableName: aws.String(tableName),
	}
	_, err = dynamodbSession.PutItem(input)
	if err != nil {
		return err
	}

	return nil


}


func createDynamoSession() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return dynamodb.New(sess)
}