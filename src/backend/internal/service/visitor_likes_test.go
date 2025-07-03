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

func TestLikesService_GetLikeCount(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*MockStorage)
		expectedCount int
		expectedError bool
		errorMessage  string
	}{
		{
			name: "successful get like count",
			mockSetup: func(m *MockStorage) {
				m.On("GetCount", mock.Anything, "likes").Return(25, nil)
			},
			expectedCount: 25,
			expectedError: false,
		},
		{
			name: "storage error on get likes",
			mockSetup: func(m *MockStorage) {
				m.On("GetCount", mock.Anything, "likes").Return(0, errors.New("storage error"))
			},
			expectedCount: 0,
			expectedError: true,
			errorMessage:  "storage error",
		},
		{
			name: "zero likes count",
			mockSetup: func(m *MockStorage) {
				m.On("GetCount", mock.Anything, "likes").Return(0, nil)
			},
			expectedCount: 0,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockStorage{}
			tt.mockSetup(mockStorage)

			service := NewLikesService(mockStorage)
			count, err := service.GetLikeCount(context.Background())

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

func TestLikesService_ToggleLike(t *testing.T) {
	tests := []struct {
		name           string
		sessionID      string
		mockSetup      func(*MockStorage)
		expectedCount  int
		expectedLiked  bool
		expectedAction string
		expectedError  bool
		errorMessage   string
	}{
		{
			name:      "successful like - first time",
			sessionID: "test-session",
			mockSetup: func(m *MockStorage) {
				m.On("GetUserSession", mock.Anything, "test-session").Return(&model.UserSession{SessionID: "test-session", HasLiked: false}, nil)
				m.On("IncrementCount", mock.Anything, "likes").Return(26, nil)
				m.On("UpdateUserSession", mock.Anything, mock.AnythingOfType("*model.UserSession")).Return(nil)
			},
			expectedCount:  26,
			expectedLiked:  true,
			expectedAction: "liked",
			expectedError:  false,
		},
		{
			name:      "successful unlike - toggle off",
			sessionID: "test-session",
			mockSetup: func(m *MockStorage) {
				m.On("GetUserSession", mock.Anything, "test-session").Return(&model.UserSession{SessionID: "test-session", HasLiked: true}, nil)
				m.On("DecrementCount", mock.Anything, "likes").Return(24, nil)
				m.On("UpdateUserSession", mock.Anything, mock.AnythingOfType("*model.UserSession")).Return(nil)
			},
			expectedCount:  24,
			expectedLiked:  false,
			expectedAction: "unliked",
			expectedError:  false,
		},
		{
			name:      "session not found - create new session and like",
			sessionID: "new-session",
			mockSetup: func(m *MockStorage) {
				m.On("GetUserSession", mock.Anything, "new-session").Return(nil, nil)
				m.On("CreateUserSession", mock.Anything, "new-session").Return(nil)
				m.On("IncrementCount", mock.Anything, "likes").Return(1, nil)
				m.On("UpdateUserSession", mock.Anything, mock.AnythingOfType("*model.UserSession")).Return(nil)
			},
			expectedCount:  1,
			expectedLiked:  true,
			expectedAction: "liked",
			expectedError:  false,
		},
		{
			name:      "storage error getting session",
			sessionID: "test-session",
			mockSetup: func(m *MockStorage) {
				m.On("GetUserSession", mock.Anything, "test-session").Return(nil, errors.New("session error"))
			},
			expectedCount:  0,
			expectedLiked:  false,
			expectedAction: "",
			expectedError:  true,
			errorMessage:   "session error",
		},
		{
			name:      "storage error creating session",
			sessionID: "new-session",
			mockSetup: func(m *MockStorage) {
				m.On("GetUserSession", mock.Anything, "new-session").Return(nil, nil)
				m.On("CreateUserSession", mock.Anything, "new-session").Return(errors.New("create failed"))
			},
			expectedCount:  0,
			expectedLiked:  false,
			expectedAction: "",
			expectedError:  true,
			errorMessage:   "create failed",
		},
		{
			name:      "storage error during increment",
			sessionID: "test-session",
			mockSetup: func(m *MockStorage) {
				m.On("GetUserSession", mock.Anything, "test-session").Return(&model.UserSession{SessionID: "test-session", HasLiked: false}, nil)
				m.On("IncrementCount", mock.Anything, "likes").Return(0, errors.New("increment failed"))
			},
			expectedCount:  0,
			expectedLiked:  false,
			expectedAction: "",
			expectedError:  true,
			errorMessage:   "increment failed",
		},
		{
			name:      "storage error during decrement",
			sessionID: "test-session",
			mockSetup: func(m *MockStorage) {
				m.On("GetUserSession", mock.Anything, "test-session").Return(&model.UserSession{SessionID: "test-session", HasLiked: true}, nil)
				m.On("DecrementCount", mock.Anything, "likes").Return(0, errors.New("decrement failed"))
			},
			expectedCount:  0,
			expectedLiked:  false,
			expectedAction: "",
			expectedError:  true,
			errorMessage:   "decrement failed",
		},
		{
			name:      "storage error updating session",
			sessionID: "test-session",
			mockSetup: func(m *MockStorage) {
				m.On("GetUserSession", mock.Anything, "test-session").Return(&model.UserSession{SessionID: "test-session", HasLiked: false}, nil)
				m.On("IncrementCount", mock.Anything, "likes").Return(26, nil)
				m.On("UpdateUserSession", mock.Anything, mock.AnythingOfType("*model.UserSession")).Return(errors.New("update failed"))
			},
			expectedCount:  26,
			expectedLiked:  true,
			expectedAction: "liked",
			expectedError:  true,
			errorMessage:   "update failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockStorage{}
			tt.mockSetup(mockStorage)

			service := NewLikesService(mockStorage)
			count, liked, action, err := service.ToggleLike(context.Background(), tt.sessionID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMessage)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCount, count)
				assert.Equal(t, tt.expectedLiked, liked)
				assert.Equal(t, tt.expectedAction, action)
			}

			mockStorage.AssertExpectations(t)
		})
	}
}
