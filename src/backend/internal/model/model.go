package model

import "time"

type Count struct {
	ID    string `dynamodbav:"ID" json:"id"`
	Count int    `dynamodbav:"Count" json:"count"`
}

type UserSession struct {
	SessionID  string    `dynamodbav:"SessionID" json:"session_id"`
	HasVisited bool      `dynamodbav:"HasVisited" json:"has_visited"`
	HasLiked   bool      `dynamodbav:"HasLiked" json:"has_liked"`
	ExpiresAt  time.Time `dynamodbav:"ExpiresAt" json:"expires_at"`
	CreatedAt  time.Time `dynamodbav:"CreatedAt" json:"created_at"`
	UpdatedAt  time.Time `dynamodbav:"UpdatedAt" json:"updated_at"`
}

type ContactRequest struct {
	ID        string    `json:"id"`
	Name      string    `json:"name" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	Message   string    `json:"message" validate:"required"`
	Recaptcha string    `json:"recaptcha" validate:"required"`
	Timestamp time.Time `json:"timestamp"`
}

type APIResponse struct {
	Count   int                    `json:"count,omitempty"`
	Message string                 `json:"message,omitempty"`
	Error   string                 `json:"error,omitempty"`
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

type NotificationPayload struct {
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	Source    string                 `json:"source"`
	Timestamp time.Time              `json:"timestamp"`
}
