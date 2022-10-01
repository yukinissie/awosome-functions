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

func fetchLanking() ([]Lank, error) {
	// Create the config specifying the Region for the DynamoDB table.
	// If Config.Region is not set the region must come from the shared
	// config or AWS_REGION environment variable.
	awscfg := &aws.Config{}
	awscfg.WithRegion("ap-northeast-1")

	// Create the session that the DynamoDB service will use.
	sess := session.Must(session.NewSession(awscfg))

	// Create the DynamoDB service client to make the query request with.
	svc := dynamodb.New(sess)

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		TableName: aws.String("bears-sandbag-lank"),
	}
	params.Limit = aws.Int64(1000)

	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	lanks := []Lank{}

	// Unmarshal the Items field in the result value to the Item Go type.
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &lanks)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return lanks, nil
}

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{
		"Content-Type":                    "application/json",
		"Access-Control-Allow-Origin":     request.Headers["origin"],
		"Access-Control-Allow-Methods":    "OPTIONS,POST,GET",
		"Access-Control-Allow-Headers":    "Origin,Authorization,Accept,X-Requested-With",
		"Access-Control-Allow-Credential": "true",
	}
	// スコアランキングをフェッチ
	lanks, err := fetchLanking()
	if err != nil {
		log.Println(err.Error())
		return events.APIGatewayProxyResponse{Headers: headers, Body: err.Error(), StatusCode: 500}, err
	}
	// JSON形式のバイナリに変換
	res, err := json.Marshal(lanks)
	if err != nil {
		log.Println(err.Error())
		return events.APIGatewayProxyResponse{Headers: headers, Body: err.Error(), StatusCode: 500}, err
	}
	return events.APIGatewayProxyResponse{Headers: headers, Body: string(res), StatusCode: 200}, nil
}

func main() {
	lambda.Start(HandleRequest)
}

type Lank struct {
	Score int    `json:"score"`
	Name  string `json:"name"`
}
