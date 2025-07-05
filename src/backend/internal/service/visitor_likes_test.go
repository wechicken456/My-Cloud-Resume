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
				m.On("GetCount", mock.Anything, "visitors").Return(42, nil)
			},
			expectedCount: 42,
			expectedError: false,
		},
		{
			name: "storage error",
			mockSetup: func(m *MockStorage) {
				m.On("GetCount", mock.Anything, "visitors").Return(0, errors.New("storage error"))
			},
			expectedCount: 0,
			expectedError: true,
			errorMessage:  "storage error",
		},
		{
			name: "zero count",
			mockSetup: func(m *MockStorage) {
				m.On("GetCount", mock.Anything, "visitors").Return(0, nil)
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
		name           string
		session        *model.UserSession
		mockSetup      func(*MockStorage)
		expectedCount  int
		expectedAction string
		expectedError  bool
		errorMessage   string
	}{
		{
			name: "successful increment - new visitors",
			session: &model.UserSession{
				SessionID:  "visitors",
				HasVisited: false,
			},
			mockSetup: func(m *MockStorage) {
				m.On("IncrementCount", mock.Anything, "visitors").Return(43, nil)
			},
			expectedCount:  43,
			expectedAction: "incremented",
			expectedError:  false,
		},
		{
			name: "already visited - return current count",
			session: &model.UserSession{
				SessionID:  "visitors",
				HasVisited: true,
			},
			mockSetup: func(m *MockStorage) {
				m.On("GetCount", mock.Anything, "visitors").Return(42, nil)
			},
			expectedCount:  42,
			expectedAction: "already_visited",
			expectedError:  false,
		},
		{
			name: "storage error during increment",
			session: &model.UserSession{
				SessionID:  "visitors",
				HasVisited: false,
			},
			mockSetup: func(m *MockStorage) {
				m.On("IncrementCount", mock.Anything, "visitors").Return(0, errors.New("increment failed"))
			},
			expectedCount:  0,
			expectedAction: "",
			expectedError:  true,
			errorMessage:   "increment failed",
		},
		{
			name: "storage error getting current count",
			session: &model.UserSession{
				SessionID:  "visitors",
				HasVisited: true,
			},
			mockSetup: func(m *MockStorage) {
				m.On("GetCount", mock.Anything, "visitors").Return(0, errors.New("get count failed"))
			},
			expectedCount:  0,
			expectedAction: "already_visited",
			expectedError:  true,
			errorMessage:   "get count failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockStorage{}
			tt.mockSetup(mockStorage)

			service := NewVisitorService(mockStorage)
			count, action, err := service.IncrementVisitorCount(context.Background(), tt.session)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMessage)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCount, count)
				assert.Equal(t, tt.expectedAction, action)
			}

			mockStorage.AssertExpectations(t)
		})
	}
}

func TestLikeService_GetLikeCount(t *testing.T) {
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

			service := NewLikeService(mockStorage)
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

func TestLikeService_ToggleLike(t *testing.T) {
	tests := []struct {
		name           string
		session        *model.UserSession
		mockSetup      func(*MockStorage)
		expectedCount  int
		expectedLiked  bool
		expectedAction string
		expectedError  bool
		errorMessage   string
	}{
		{
			name: "successful like - first time",
			session: &model.UserSession{
				SessionID: "test-session",
				HasLiked:  false,
			},
			mockSetup: func(m *MockStorage) {
				m.On("IncrementCount", mock.Anything, "likes").Return(26, nil)
			},
			expectedCount:  26,
			expectedLiked:  true,
			expectedAction: "liked",
			expectedError:  false,
		},
		{
			name: "successful unlike - toggle off",
			session: &model.UserSession{
				SessionID: "test-session",
				HasLiked:  true,
			},
			mockSetup: func(m *MockStorage) {
				m.On("DecrementCount", mock.Anything, "likes").Return(24, nil)
			},
			expectedCount:  24,
			expectedLiked:  false,
			expectedAction: "unliked",
			expectedError:  false,
		},
		{
			name: "storage error during increment",
			session: &model.UserSession{
				SessionID: "test-session",
				HasLiked:  false,
			},
			mockSetup: func(m *MockStorage) {
				m.On("IncrementCount", mock.Anything, "likes").Return(0, errors.New("increment failed"))
			},
			expectedCount:  0,
			expectedLiked:  false,
			expectedAction: "",
			expectedError:  true,
			errorMessage:   "increment failed",
		},
		{
			name: "storage error during decrement",
			session: &model.UserSession{
				SessionID: "test-session",
				HasLiked:  true,
			},
			mockSetup: func(m *MockStorage) {
				m.On("DecrementCount", mock.Anything, "likes").Return(0, errors.New("decrement failed"))
			},
			expectedCount:  0,
			expectedLiked:  false,
			expectedAction: "",
			expectedError:  true,
			errorMessage:   "decrement failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockStorage{}
			tt.mockSetup(mockStorage)

			service := NewLikeService(mockStorage)
			count, action, err := service.ToggleLike(context.Background(), tt.session)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMessage)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCount, count)
				assert.Equal(t, tt.expectedAction, action)
				// Check that the session's HasLiked field was updated correctly
				assert.Equal(t, tt.expectedLiked, tt.session.HasLiked)
			}

			mockStorage.AssertExpectations(t)
		})
	}
}
