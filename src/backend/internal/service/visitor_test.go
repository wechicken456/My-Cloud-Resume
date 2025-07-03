package service

import (
	"context"
	"errors"
	"main/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStorage is a mock implementation of the storage interface
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) GetCount(ctx context.Context, countName string) (int, error) {
	args := m.Called(ctx, countName)
	return args.Int(0), args.Error(1)
}

func (m *MockStorage) IncrementCount(ctx context.Context, countName string) (int, error) {
	args := m.Called(ctx, countName)
	return args.Int(0), args.Error(1)
}

func (m *MockStorage) DecrementCount(ctx context.Context, countName string) (int, error) {
	args := m.Called(ctx, countName)
	return args.Int(0), args.Error(1)
}

func (m *MockStorage) GetUserSession(ctx context.Context, sessionID string) (*model.UserSession, error) {
	args := m.Called(ctx, sessionID)
	if session, ok := args.Get(0).(*model.UserSession); ok {
		return session, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockStorage) CreateUserSession(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockStorage) UpdateUserSession(ctx context.Context, session *model.UserSession) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func TestVisitorService_GetCount(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*MockStorage)
		expectedCount int
		expectedError bool
		errorMessage  string
	}{
		{
			name: "successful get count",
			mockSetup: func(m *MockStorage) {
				m.On("GetCount", mock.Anything, "visitor").Return(42, nil)
			},
			expectedCount: 42,
			expectedError: false,
		},
		{
			name: "storage error",
			mockSetup: func(m *MockStorage) {
				m.On("GetCount", mock.Anything, "visitor").Return(0, errors.New("storage error"))
			},
			expectedCount: 0,
			expectedError: true,
			errorMessage:  "storage error",
		},
		{
			name: "zero count",
			mockSetup: func(m *MockStorage) {
				m.On("GetCount", mock.Anything, "visitor").Return(0, nil)
			},
			expectedCount: 0,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockStorage{}
			tt.mockSetup(mockStorage)

			service := NewVisitorService(mockStorage)
			count, err := service.GetVisitorCount(context.Background())

			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMessage)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCount, count)
			}

			mockStorage.AssertExpectations(t)
		})
	}
}

func TestVisitorService_IncrementCount(t *testing.T) {
	tests := []struct {
		name                string
		mockSetup           func(*MockStorage)
		expectedCount       int
		expectedIncremented bool
		expectedAction      string
		expectedError       bool
		errorMessage        string
	}{
		{
			name: "successful increment",
			mockSetup: func(m *MockStorage) {
				m.On("GetUserSession", mock.Anything, "visitor").Return(&model.UserSession{SessionID: "visitor", HasVisited: false}, nil)
				m.On("IncrementCount", mock.Anything, "visitor").Return(43, nil)
				m.On("UpdateUserSession", mock.Anything, mock.AnythingOfType("*model.UserSession")).Return(nil)
			},
			expectedCount:       43,
			expectedIncremented: true,
			expectedAction:      "incremented",
			expectedError:       false,
		},
		{
			name: "storage error during increment",
			mockSetup: func(m *MockStorage) {
				m.On("GetUserSession", mock.Anything, "visitor").Return(nil, errors.New("session not found"))
				m.On("IncrementCount", mock.Anything, "visitor").Return(0, errors.New("session not found"))
			},
			expectedCount:       0,
			expectedIncremented: false,
			expectedAction:      "",
			expectedError:       true,
			errorMessage:        "session not found",
		},
		{
			name: "first increment",
			mockSetup: func(m *MockStorage) {
				m.On("GetUserSession", mock.Anything, "visitor").Return(&model.UserSession{SessionID: "visitor", HasVisited: false}, nil)
				m.On("IncrementCount", mock.Anything, "visitor").Return(1, nil)
				m.On("UpdateUserSession", mock.Anything, mock.AnythingOfType("*model.UserSession")).Return(nil)
			},
			expectedCount:       1,
			expectedIncremented: true,
			expectedAction:      "incremented",
			expectedError:       false,
		},
		{
			name: "duplicate increment - already visited",
			mockSetup: func(m *MockStorage) {
				m.On("GetUserSession", mock.Anything, "visitor").Return(&model.UserSession{SessionID: "visitor", HasVisited: true}, nil)
				m.On("GetCount", mock.Anything, "visitor").Return(1, nil)
			},
			expectedCount:       1,
			expectedIncremented: false,
			expectedAction:      "already_counted",
			expectedError:       false,
		},
		{
			name: "session not found - create new session",
			mockSetup: func(m *MockStorage) {
				m.On("GetUserSession", mock.Anything, "visitor").Return(nil, nil)
				m.On("CreateUserSession", mock.Anything, "visitor").Return(nil)
				m.On("IncrementCount", mock.Anything, "visitor").Return(1, nil)
				m.On("UpdateUserSession", mock.Anything, mock.AnythingOfType("*model.UserSession")).Return(nil)
			},
			expectedCount:       1,
			expectedIncremented: true,
			expectedAction:      "incremented",
			expectedError:       false,
		},
		{
			name: "update session error",
			mockSetup: func(m *MockStorage) {
				m.On("GetUserSession", mock.Anything, "visitor").Return(&model.UserSession{SessionID: "visitor", HasVisited: false}, nil)
				m.On("IncrementCount", mock.Anything, "visitor").Return(43, nil)
				m.On("UpdateUserSession", mock.Anything, mock.AnythingOfType("*model.UserSession")).Return(errors.New("update failed"))
			},
			expectedCount:       0,
			expectedIncremented: false,
			expectedAction:      "",
			expectedError:       true,
			errorMessage:        "update failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockStorage{}
			tt.mockSetup(mockStorage)

			service := NewVisitorService(mockStorage)
			count, incremented, action, err := service.IncrementVisitorCount(context.Background(), "visitor")

			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMessage)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCount, count)
				assert.Equal(t, tt.expectedIncremented, incremented)
				assert.Equal(t, tt.expectedAction, action)
			}
		})
	}
}
