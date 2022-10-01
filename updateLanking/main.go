package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func updateLanking(lank Lank) error {
	// Create the config specifying the Region for the DynamoDB table.
	// If Config.Region is not set the region must come from the shared
	// config or AWS_REGION environment variable.
	awscfg := &aws.Config{}
	awscfg.WithRegion("ap-northeast-1")

	// Create the session that the DynamoDB service will use.
	sess := session.Must(session.NewSession(awscfg))

	// Create the DynamoDB service client to make the query request with.
	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(lank)
	if err != nil {
		return err
	}
	// Build the query input parameters
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("bears-sandbag-lank"),
	}

	// Make the DynamoDB Query API call
	_, err = svc.PutItem(input)
	if err != nil {
		return err
	}
	return nil
}

func ConvertInputDataToStruct(inputs string) (*Lank, error) {
	var req Lank
	err := json.Unmarshal([]byte(inputs), &req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{
		"Content-Type":                    "application/json",
		"Access-Control-Allow-Origin":     request.Headers["origin"],
		"Access-Control-Allow-Methods":    "OPTIONS,POST,GET",
		"Access-Control-Allow-Headers":    "Origin,Authorization,Accept,X-Requested-With",
		"Access-Control-Allow-Credential": "true",
	}
	req, err := ConvertInputDataToStruct(request.Body)
	if err != nil {
		log.Println(err.Error())
		return events.APIGatewayProxyResponse{Headers: headers, Body: err.Error(), StatusCode: 500}, err
	}
	// スコアランキングをフェッチ
	err = updateLanking(*req)
	if err != nil {
		log.Println(err.Error())
		return events.APIGatewayProxyResponse{Headers: headers, Body: err.Error(), StatusCode: 500}, err
	}
	return events.APIGatewayProxyResponse{Headers: headers, Body: "ok", StatusCode: 200}, nil
}

func main() {
	lambda.Start(HandleRequest)
}

type Lank struct {
	Score int    `json:"score"`
	Name  string `json:"name"`
}

type Config struct {
	Table  string // required
	Region string // optional
	Limit  int64  // optional
}
