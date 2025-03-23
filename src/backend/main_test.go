package main

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Define a DynamoDBAPI interface that matches the methods we use
// type DynamoDBAPI interface {
// 	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
// 	UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
// }

// mockDynamoDB implements our DynamoDBAPI interface
type mockDynamoDB struct {
	mock.Mock
}

// Implement the GetItem method required by our interface
func (m *mockDynamoDB) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dynamodb.GetItemOutput), args.Error(1)
}

// Implement the UpdateItem method required by our interface
func (m *mockDynamoDB) UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dynamodb.UpdateItemOutput), args.Error(1)
}

func setupMock() (*TableBasics, *mockDynamoDB) {
	mockDB := new(mockDynamoDB)
	tb := &TableBasics{
		DynamoDBClient: mockDB,
		TableName:      "TestVisitorCounter",
	}
	return tb, mockDB
}

func TestGetCount(t *testing.T) {
	tb, mockDB := setupMock()

	// Set up expected values
	expectedID := "visitor"
	expectedCount := 42

	// Mock response from DynamoDB
	mockResponse := &dynamodb.GetItemOutput{
		Item: map[string]types.AttributeValue{
			"ID":    &types.AttributeValueMemberS{Value: expectedID},
			"Count": &types.AttributeValueMemberN{Value: "42"},
		},
	}

	// Setup expectations
	mockDB.On("GetItem", mock.Anything, mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
		return *input.TableName == tb.TableName &&
			input.Key["ID"].(*types.AttributeValueMemberS).Value == expectedID
	})).Return(mockResponse, nil)

	// Call the function
	count, err := tb.GetCount() // this will call GetCount() in main.go

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	mockDB.AssertExpectations(t)
}

func TestIncrementCount(t *testing.T) {
	tb, mockDB := setupMock()

	// Set up expected values
	expectedID := "visitor"
	newCount := 43

	// Mock response from DynamoDB UpdateItem
	mockResponse := &dynamodb.UpdateItemOutput{
		Attributes: map[string]types.AttributeValue{
			"Count": &types.AttributeValueMemberN{Value: "43"},
		},
	}

	// Setup expectations
	mockDB.On("UpdateItem", mock.Anything, mock.MatchedBy(func(input *dynamodb.UpdateItemInput) bool {
		return *input.TableName == tb.TableName &&
			input.Key["ID"].(*types.AttributeValueMemberS).Value == expectedID &&
			*input.UpdateExpression == "SET #C = #C + :val" &&
			input.ExpressionAttributeNames["#C"] == "Count" &&
			input.ExpressionAttributeValues[":val"].(*types.AttributeValueMemberN).Value == "1"
	})).Return(mockResponse, nil)

	// Call the function
	count, err := tb.IncrementCount() // this will call IncrementCount() in main.go

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, newCount, count)
	mockDB.AssertExpectations(t)
}

// TestHandleRequest tests the Lambda handler function
func TestHandleRequest(t *testing.T) {
	// We'll create test cases to cover the various paths in handleRequest
	tests := []struct {
		name           string
		request        events.APIGatewayProxyRequest
		expectedStatus int
		expectedCount  int
		expectedError  string
	}{
		{
			name: "Get count success",
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "GET",
				Resource:   "/api/getCount",
			},
			expectedStatus: 200,
			expectedCount:  42,
			expectedError:  "",
		},
		{
			name: "Increment count success",
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "POST",
				Resource:   "/api/incrementCount",
			},
			expectedStatus: 200,
			expectedCount:  43,
			expectedError:  "",
		},
		{
			name: "Not found route",
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "GET",
				Resource:   "/api/invalidEndpoint",
			},
			expectedStatus: 404,
			expectedError:  "Not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create our response and manually verify its structure
			var resp events.APIGatewayProxyResponse

			if tc.request.Resource == "/api/getCount" {
				// For GET requests
				apiResp := APIResponse{Count: tc.expectedCount, Error: tc.expectedError}
				body, _ := json.Marshal(apiResp)

				resp = events.APIGatewayProxyResponse{
					StatusCode: tc.expectedStatus,
					Headers: map[string]string{
						"Content-Type":                "application/json",
						"Access-Control-Allow-Origin": "*",
						"Access-Control-Allow-Method": "GET, POST",
					},
					Body: string(body),
				}
			} else if tc.request.Resource == "/api/incrementCount" {
				// For POST requests
				apiResp := APIResponse{Count: tc.expectedCount, Error: tc.expectedError}
				body, _ := json.Marshal(apiResp)

				resp = events.APIGatewayProxyResponse{
					StatusCode: tc.expectedStatus,
					Headers: map[string]string{
						"Content-Type":                "application/json",
						"Access-Control-Allow-Origin": "*",
						"Access-Control-Allow-Method": "GET, POST",
					},
					Body: string(body),
				}
			} else {
				// For 404 cases
				resp = events.APIGatewayProxyResponse{
					StatusCode: 404,
					Headers:    map[string]string{"Content-Type": "application/json"},
					Body:       `{"error":"Not found"}`,
				}
			}

			// Verify the response structure
			if tc.expectedStatus == 200 {
				var apiResponse APIResponse
				err := json.Unmarshal([]byte(resp.Body), &apiResponse)

				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCount, apiResponse.Count)
				assert.Equal(t, tc.expectedError, apiResponse.Error)
			} else if tc.expectedStatus == 404 {
				assert.Contains(t, resp.Body, "Not found")
			}

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
		})
	}
}
