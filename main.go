package main

import (
	"encoding/json"
	"net/http"

	"lambda-dynamodb-users/dynamodb"
	"lambda-dynamodb-users/types"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(Handler)

}

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.HTTPMethod == "POST" {
		var user types.User
		err := json.Unmarshal([]byte(req.Body), &user)
		if err != nil {
			return response("Error de Unmarshal", http.StatusBadRequest), err
		}
		err = dynamodb.SaveUser(user)
		if err != nil {
			return response("Error saving user", http.StatusInternalServerError), err
		}
		return response("User saved", http.StatusOK), nil
	}

	if req.HTTPMethod == "GET" {
		var user types.User
		err:= json.Unmarshal([]byte(req.Body), &user)
		if err != nil {
			return response("Error de Unmarshal", http.StatusBadRequest), err
		}
		id:= user.ID
		user , err = dynamodb.GetUser(id)
		if err != nil {
			return response("Error getting user", http.StatusInternalServerError), err
		}
		userresponse, _ := json.Marshal(user)
		return responseUser(userresponse, http.StatusOK), nil
	} else {
		return response("Method not allowed", http.StatusMethodNotAllowed), nil
	}
}

func responseUser(body []byte, statusCode int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: statusCode,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
	}
}

func response(body string, statusCode int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(body),
		Headers: map[string]string{
			"Acces-Control-Allow-Origin": "*",
		},
	}

}
