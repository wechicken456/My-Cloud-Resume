package storage

import (
	"context"
	"fmt"
	"log"
	"main/internal/model"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBAPI interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

type StorageInterface interface {
	GetCount(ctx context.Context, countName string) (int, error)
	IncrementCount(ctx context.Context, countName string) (int, error)
	DecrementCount(ctx context.Context, countName string) (int, error)
	GetUserSession(ctx context.Context, sessionID string) (*model.UserSession, error)
	CreateUserSession(ctx context.Context, sessionID string) error
	UpdateUserSession(ctx context.Context, session *model.UserSession) error
}

type Storage struct {
	client       DynamoDBAPI
	tableName    string
	sessionTable string
}

func New(client DynamoDBAPI, tableName, sessionTable string) *Storage {
	return &Storage{
		client:       client,
		tableName:    tableName,
		sessionTable: sessionTable,
	}
}

func (s *Storage) GetClient() DynamoDBAPI {
	return s.client
}

func (s *Storage) GetTableName() string {
	return s.tableName
}

func (s *Storage) GetSessionTable() string {
	return s.sessionTable
}

// retrieve count from DynamoDB and return json: {"count": ret} if successful
func (s *Storage) GetCount(ctx context.Context, countName string) (int, error) {
	// required argument for UpdateItemInput
	key := map[string]types.AttributeValue{
		"ID": &types.AttributeValueMemberS{Value: countName}, // Value is the name of the ID that we set for the counter in DynamoDB
	}

	response, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{TableName: &s.tableName, Key: key})
	if err != nil {
		return 0, err
	}

	// Check if item exists
	if response.Item == nil {
		return 0, fmt.Errorf("no item found with ID %s", countName)
	}

	var vc model.Count
	err = attributevalue.UnmarshalMap(response.Item, &vc)
	if err != nil {
		log.Printf("Couldn't unmarshal response. Here's why: %v\n", err)
		return 0, err
	}

	return vc.Count, nil
}

// use the DynamoDB client's UpdateItem() to increment the counter
// return the new incremented count as json: {"count": ret} if successful
func (s *Storage) IncrementCount(ctx context.Context, countName string) (int, error) {

	updateInput := dynamodb.UpdateItemInput{
		TableName:                 &s.tableName,
		Key:                       map[string]types.AttributeValue{"ID": &types.AttributeValueMemberS{Value: countName}}, // Value is the name of the ID that we set for the counter in DynamoDB
		UpdateExpression:          aws.String("SET #C = #C + :val"),
		ExpressionAttributeNames:  map[string]string{"#C": "Count"},
		ExpressionAttributeValues: map[string]types.AttributeValue{":val": &types.AttributeValueMemberN{Value: "1"}},
		ReturnValues:              types.ReturnValueUpdatedNew, // Returns only the updated attributes, as they appear after theUpdateItem operation
	}
	result, err := s.client.UpdateItem(ctx, &updateInput)
	if err != nil {
		return 0, fmt.Errorf("failed to increment Count: %v", err)
	}

	var newCount int
	err = attributevalue.Unmarshal(result.Attributes["Count"], &newCount)
	if err != nil {
		return 0, err
	}
	return newCount, nil
}

func (s *Storage) DecrementCount(ctx context.Context, countName string) (int, error) {
	updateInput := dynamodb.UpdateItemInput{
		TableName:                &s.tableName,
		Key:                      map[string]types.AttributeValue{"ID": &types.AttributeValueMemberS{Value: countName}},
		UpdateExpression:         aws.String("SET #C = if_not_exists(#C, :zero) - :val"),
		ConditionExpression:      aws.String("#C > :zero"), // Prevent negative counts
		ExpressionAttributeNames: map[string]string{"#C": "Count"},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":val":  &types.AttributeValueMemberN{Value: "1"},
			":zero": &types.AttributeValueMemberN{Value: "0"},
		},
		ReturnValues: types.ReturnValueUpdatedNew,
	}
	result, err := s.client.UpdateItem(context.Background(), &updateInput)
	if err != nil {
		// If condition fails (count would go negative), return current count
		if err.Error() == "ConditionalCheckFailedException" {
			return s.GetCount(ctx, countName)
		}
		return 0, fmt.Errorf("failed to decrement Count: %v", err)
	}

	var newCount int
	err = attributevalue.Unmarshal(result.Attributes["Count"], &newCount)
	if err != nil {
		return 0, err
	}
	return newCount, nil
}

func (s *Storage) GetUserSession(ctx context.Context, sessionID string) (*model.UserSession, error) {
	key := map[string]types.AttributeValue{
		"SessionID": &types.AttributeValueMemberS{Value: sessionID},
	}

	response, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &s.sessionTable,
		Key:       key,
	})
	if err != nil {
		return nil, err
	}

	if response.Item == nil {
		return nil, nil
	}

	var session model.UserSession
	err = attributevalue.UnmarshalMap(response.Item, &session)
	if err != nil {
		return nil, err
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		return nil, nil
	}

	return &session, nil
}

func (s *Storage) CreateUserSession(ctx context.Context, sessionID string) error {
	session := model.UserSession{
		SessionID:  sessionID,
		HasVisited: true,
		HasLiked:   false,
		ExpiresAt:  time.Now().Add(24 * time.Hour),
		CreatedAt:  time.Now(),
	}

	item, err := attributevalue.MarshalMap(session)
	if err != nil {
		return err
	}

	_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &s.sessionTable,
		Item:      item,
	})
	return err
}

func (s *Storage) UpdateUserSession(ctx context.Context, session *model.UserSession) error {
	update := expression.Set(expression.Name("HasVisited"), expression.Value(session.HasVisited))
	update.Set(expression.Name("HasLiked"), expression.Value(session.HasLiked))
	update.Set(expression.Name("UpdatedAt"), expression.Value(time.Now()))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return fmt.Errorf("failed to build expression: %v", err)
	}

	_, err = s.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 &s.sessionTable,
		Key:                       map[string]types.AttributeValue{"SessionID": &types.AttributeValueMemberS{Value: session.SessionID}},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ReturnValues:              types.ReturnValueUpdatedNew,
	})

	return err
}
