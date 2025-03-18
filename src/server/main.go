package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

type VisitorCount struct {
	ID    string `dynamodbav:"ID"`
	Count int    `dynamodbav:"Count"`
}

type APIResponse struct {
	Count int    `json:"count"`
	Error string `json:"error"`
}

type TableBasics struct {
	DynamoDBClient *dynamodb.Client
	TableName      string
}

func (cnt *VisitorCount) GetKey() map[string]types.AttributeValue {
	count, err := attributevalue.Marshal(cnt.Count)
	if err != nil {
		panic(err)
	}
	return map[string]types.AttributeValue{"ID": count}
}

// retrieve count from DynamoDB and return json: {"count": ret} if successful
func (tb *TableBasics) GetCount() (int, error) {
	// required argument for UpdateItemInput
	id := "visitor"
	key := map[string]types.AttributeValue{
		"ID": &types.AttributeValueMemberS{Value: id}, // Value is the name of the ID that we set for the counter in DynamoDB
	}

	response, err := tb.DynamoDBClient.GetItem(context.Background(), &dynamodb.GetItemInput{TableName: &tb.TableName, Key: key})
	if err != nil {
		return 0, err
	}
	// Check if item exists
	if response.Item == nil {
		return 0, fmt.Errorf("no item found with ID %s", id)
	}

	var vc VisitorCount
	err = attributevalue.UnmarshalMap(response.Item, &vc)
	if err != nil {
		log.Printf("Couldn't unmarshal response. Here's why: %v\n", err)
		return 0, err
	}

	return vc.Count, nil
}

// use the DynamoDB client's UpdateItem() to increment the counter
// return the new incremented count as json: {"count": ret} if successful
func (tb *TableBasics) IncrementCount() (int, error) {

	updateInput := dynamodb.UpdateItemInput{
		TableName:                 &tb.TableName,
		Key:                       map[string]types.AttributeValue{"ID": &types.AttributeValueMemberS{Value: "visitor"}}, // Value is the name of the ID that we set for the counter in DynamoDB
		UpdateExpression:          aws.String("SET count = count + :val"),
		ExpressionAttributeValues: map[string]types.AttributeValue{":val": &types.AttributeValueMemberN{Value: "1"}},
		ReturnValues:              types.ReturnValueUpdatedNew, // Returns only the updated attributes, as they appear after theUpdateItem operation
	}
	result, err := tb.DynamoDBClient.UpdateItem(context.Background(), &updateInput)
	if err != nil {
		return 0, fmt.Errorf("failed to increment count: %v", err)
	}

	var newCount int
	err = attributevalue.Unmarshal(result.Attributes["Count"], &newCount)
	if err != nil {
		return 0, err
	}
	return newCount, nil
}

func handleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse the input event
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Couldn't load AWS config: %s", err)
	}
	tb := TableBasics{TableName: "VisitorCounter", DynamoDBClient: dynamodb.NewFromConfig(cfg)}

	var count int
	if req.HTTPMethod == "GET" && req.Resource == "/api/getCount" {
		count, err = tb.GetCount()
		if err != nil {
			log.Printf("Error getting count: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       `{"error": "Database error"}`,
			}, err
		}

	} else if req.HTTPMethod == "POST" && req.Resource == "/api/incrementCount" {
		count, err = tb.IncrementCount()
		if err != nil {
			log.Printf("Error incrementing count: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       `{"error": "Database error"}`,
			}, err
		}

	} else {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"error": "Not found"}`,
		}, nil
	}

	body, err := json.Marshal(APIResponse{Count: count})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"error": "Internal server error"}`,
		}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
			"Access-Control-Allow-Method": "GET, POST",
		},
		Body: string(body),
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
