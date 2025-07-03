package storage

import (
	"context"
	"errors"
	"main/internal/model"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDynamoDBAPI for testing
type MockDynamoDBAPI struct {
	mock.Mock
}

func (m *MockDynamoDBAPI) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dynamodb.GetItemOutput), args.Error(1)
}

func (m *MockDynamoDBAPI) UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dynamodb.UpdateItemOutput), args.Error(1)
}

func (m *MockDynamoDBAPI) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dynamodb.PutItemOutput), args.Error(1)
}

func TestStorage_GetCount(t *testing.T) {
	tests := []struct {
		name        string
		countName   string
		setupMock   func(*MockDynamoDBAPI)
		expectedVal int
		expectedErr bool
	}{
		{
			name:      "successful get count",
			countName: "visitor_count",
			setupMock: func(m *MockDynamoDBAPI) {
				m.On("GetItem", mock.Anything, mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
					return *input.TableName == "test-table"
				})).Return(&dynamodb.GetItemOutput{
					Item: map[string]types.AttributeValue{
						"ID":    &types.AttributeValueMemberS{Value: "visitor_count"},
						"Count": &types.AttributeValueMemberN{Value: "42"},
					},
				}, nil)
			},
			expectedVal: 42,
			expectedErr: false,
		},
		{
			name:      "item not found",
			countName: "nonexistent_count",
			setupMock: func(m *MockDynamoDBAPI) {
				m.On("GetItem", mock.Anything, mock.Anything).Return(&dynamodb.GetItemOutput{
					Item: nil,
				}, nil)
			},
			expectedVal: 0,
			expectedErr: true,
		},
		{
			name:      "dynamodb error",
			countName: "visitor_count",
			setupMock: func(m *MockDynamoDBAPI) {
				m.On("GetItem", mock.Anything, mock.Anything).Return(
					&dynamodb.GetItemOutput{}, errors.New("dynamodb error"))
			},
			expectedVal: 0,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDynamoDBAPI)
			tt.setupMock(mockDB)

			storage := New(mockDB, "test-table", "test-session-table")
			count, err := storage.GetCount(context.Background(), tt.countName)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedVal, count)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func TestStorage_IncrementCount(t *testing.T) {
	tests := []struct {
		name        string
		countName   string
		setupMock   func(*MockDynamoDBAPI)
		expectedVal int
		expectedErr bool
	}{
		{
			name:      "successful increment",
			countName: "visitor_count",
			setupMock: func(m *MockDynamoDBAPI) {
				m.On("UpdateItem", mock.Anything, mock.MatchedBy(func(input *dynamodb.UpdateItemInput) bool {
					return *input.TableName == "test-table"
				})).Return(&dynamodb.UpdateItemOutput{
					Attributes: map[string]types.AttributeValue{
						"Count": &types.AttributeValueMemberN{Value: "43"},
					},
				}, nil)
			},
			expectedVal: 43,
			expectedErr: false,
		},
		{
			name:      "dynamodb error",
			countName: "visitor_count",
			setupMock: func(m *MockDynamoDBAPI) {
				m.On("UpdateItem", mock.Anything, mock.Anything).Return(
					&dynamodb.UpdateItemOutput{}, errors.New("update failed"))
			},
			expectedVal: 0,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDynamoDBAPI)
			tt.setupMock(mockDB)

			storage := New(mockDB, "test-table", "test-session-table")
			count, err := storage.IncrementCount(context.Background(), tt.countName)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedVal, count)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func TestStorage_DecrementCount(t *testing.T) {
	tests := []struct {
		name        string
		countName   string
		setupMock   func(*MockDynamoDBAPI)
		expectedVal int
		expectedErr bool
	}{
		{
			name:      "successful decrement",
			countName: "like_count",
			setupMock: func(m *MockDynamoDBAPI) {
				m.On("UpdateItem", mock.Anything, mock.MatchedBy(func(input *dynamodb.UpdateItemInput) bool {
					return *input.TableName == "test-table"
				})).Return(&dynamodb.UpdateItemOutput{
					Attributes: map[string]types.AttributeValue{
						"Count": &types.AttributeValueMemberN{Value: "5"},
					},
				}, nil)
			},
			expectedVal: 5,
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDynamoDBAPI)
			tt.setupMock(mockDB)

			storage := New(mockDB, "test-table", "test-session-table")
			count, err := storage.DecrementCount(context.Background(), tt.countName)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedVal, count)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func TestStorage_UserSession(t *testing.T) {
	t.Run("get existing session", func(t *testing.T) {
		mockDB := new(MockDynamoDBAPI)
		futureTime := time.Now().Add(time.Hour)

		mockDB.On("GetItem", mock.Anything, mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
			return *input.TableName == "test-session-table"
		})).Return(&dynamodb.GetItemOutput{
			Item: map[string]types.AttributeValue{
				"SessionID":  &types.AttributeValueMemberS{Value: "test-session"},
				"HasVisited": &types.AttributeValueMemberBOOL{Value: true},
				"HasLiked":   &types.AttributeValueMemberBOOL{Value: false},
				"ExpiresAt":  &types.AttributeValueMemberS{Value: futureTime.Format(time.RFC3339)},
				"CreatedAt":  &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
			},
		}, nil)

		storage := New(mockDB, "test-table", "test-session-table")
		session, err := storage.GetUserSession(context.Background(), "test-session")

		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, "test-session", session.SessionID)
		assert.True(t, session.HasVisited)
		assert.False(t, session.HasLiked)

		mockDB.AssertExpectations(t)
	})

	t.Run("create new session", func(t *testing.T) {
		mockDB := new(MockDynamoDBAPI)

		mockDB.On("PutItem", mock.Anything, mock.MatchedBy(func(input *dynamodb.PutItemInput) bool {
			return *input.TableName == "test-session-table"
		})).Return(&dynamodb.PutItemOutput{}, nil)

		storage := New(mockDB, "test-table", "test-session-table")
		err := storage.CreateUserSession(context.Background(), "new-session")

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("update existing session", func(t *testing.T) {
		mockDB := new(MockDynamoDBAPI)
		session := &model.UserSession{
			SessionID:  "test-session",
			HasVisited: true,
			HasLiked:   true,
		}
		mockDB.On("UpdateItem", mock.Anything, mock.MatchedBy(func(input *dynamodb.UpdateItemInput) bool {
			return *input.TableName == "test-session-table"
		})).Return(&dynamodb.UpdateItemOutput{}, nil)

		storage := New(mockDB, "test-table", "test-session-table")
		err := storage.UpdateUserSession(context.Background(), session)
		assert.NoError(t, err)
	})
}
