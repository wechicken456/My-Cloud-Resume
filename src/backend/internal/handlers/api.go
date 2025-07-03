package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"main/internal/model"
	"main/internal/service"

	"github.com/aws/aws-lambda-go/events"
)

type APIHandler struct {
	visitorService      *service.VisitorService
	likesService        *service.LikesService
	contactService      *service.ContactService
	notificationService *service.NotificationService
}

func NewAPIHandler(
	visitorService *service.VisitorService,
	likesService *service.LikesService,
	contactService *service.ContactService,
	notificationService *service.NotificationService,
) *APIHandler {
	return &APIHandler{
		visitorService:      visitorService,
		likesService:        likesService,
		contactService:      contactService,
		notificationService: notificationService,
	}
}

var session_id_string string = "session_id"

func (h *APIHandler) HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{
		"Content-Type":                     "application/json",
		"Access-Control-Allow-Origin":      "https://www.pwnph0fun.com",
		"Access-Control-Allow-Methods":     "GET, POST, OPTIONS",
		"Access-Control-Allow-Headers":     "Content-Type, Cookie",
		"Access-Control-Allow-Credentials": "true",
	}

	if req.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    headers,
			Body:       "",
		}, nil
	}

	// Extract session ID from cookie or create new one
	sessionID := h.extractSessionID(req.Headers["Cookie"])
	if sessionID == "" {
		sessionID = generateSessionID()
	}

	// Set session cookie in response
	sessionCookie := fmt.Sprintf("session_id=%s; HttpOnly; Secure; SameSite=Strict; Max-Age=86400; Path=/", sessionID)
	headers["Set-Cookie"] = sessionCookie

	switch {
	case req.HTTPMethod == "GET" && req.Resource == "/api/getVisitorCount":
		return h.handleGetVisitorCount(ctx, headers)
	case req.HTTPMethod == "POST" && req.Resource == "/api/incrementVisitorCount":
		return h.handleIncrementVisitorCount(ctx, sessionID, headers)
	case req.HTTPMethod == "GET" && req.Resource == "/api/getLikeCount":
		return h.handleGetLikeCount(ctx, headers)
	case req.HTTPMethod == "POST" && req.Resource == "/api/toggleLike":
		return h.handleToggleLike(ctx, sessionID, headers)
	case req.HTTPMethod == "POST" && req.Resource == "/api/contact":
		return h.handleContact(ctx, req.Body, headers)
	case req.HTTPMethod == "GET" && req.Resource == "/api/session":
		return h.handleGetSession(ctx, sessionID, headers)
	default:
		headers["Set-Cookie"] = ""
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Headers:    headers,
			Body:       `{"error": "Not found", "success": false}`,
		}, nil
	}
}

func (h *APIHandler) extractSessionID(cookieHeader string) string {
	if cookieHeader == "" {
		return ""
	}

	// Parse cookie header to extract session_id
	cookies := parseCookies(cookieHeader)
	return cookies[session_id_string]
}

func parseCookies(cookieHeader string) map[string]string {
	cookies := make(map[string]string)
	if cookieHeader == "" {
		return cookies
	}

	parts := strings.Split(cookieHeader, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if eq := strings.Index(part, "="); eq >= 0 {
			key := strings.TrimSpace(part[:eq])
			value := strings.TrimSpace(part[eq+1:])
			cookies[key] = value
		}
	}
	return cookies
}

func generateSessionID() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (h *APIHandler) handleGetVisitorCount(ctx context.Context, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	count, err := h.visitorService.GetVisitorCount(ctx)
	if err != nil {
		log.Printf("Error getting count: %v", err)
		return h.errorResponse(500, "Database error", headers), nil
	}

	response := model.APIResponse{Count: count, Success: true}
	body, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       string(body),
	}, nil
}

func (h *APIHandler) handleIncrementVisitorCount(ctx context.Context, sessionID string, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	count, _, status, err := h.visitorService.IncrementVisitorCount(ctx, sessionID)
	if err != nil {
		log.Printf("Error incrementing count: %v", err)
		return h.errorResponse(500, "Database error", headers), nil
	}

	response := model.APIResponse{
		Count:   count,
		Success: true,
		Message: status,
	}

	body, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       string(body),
	}, nil
}

func (h *APIHandler) handleGetLikeCount(ctx context.Context, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	count, err := h.likesService.GetLikeCount(ctx)
	if err != nil {
		log.Printf("Error getting likes: %v", err)
		return h.errorResponse(500, "Database error", headers), nil
	}

	response := model.APIResponse{Count: count, Success: true}
	body, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       string(body),
	}, nil
}

func (h *APIHandler) handleToggleLike(ctx context.Context, sessionID string, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	count, _, action, err := h.likesService.ToggleLike(ctx, sessionID)
	if err != nil {
		log.Printf("Error toggling like: %v", err)
		return h.errorResponse(500, "Database error", headers), nil
	}

	// Send notification if this is a new like
	if action == "liked" {
		payload := &model.NotificationPayload{
			Type:      "like",
			Data:      map[string]interface{}{},
			Source:    "resume-website",
			Timestamp: time.Now(),
		}
		go h.notificationService.SendEmailNotification(context.Background(), payload)
	}

	response := model.APIResponse{
		Count:   count,
		Success: true,
		Message: action,
	}

	body, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       string(body),
	}, nil
}

func (h *APIHandler) handleGetSession(ctx context.Context, sessionID string, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	session, err := h.visitorService.GetSessionStatus(ctx, sessionID)
	if err != nil {
		log.Printf("Error getting session: %v", err)
		return h.errorResponse(500, "Database error", headers), nil
	}

	response := model.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			session_id_string: sessionID,
			"visited":         session.HasVisited,
			"liked":           session.HasLiked,
		},
	}

	body, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       string(body),
	}, nil
}

func (h *APIHandler) handleContact(ctx context.Context, body string, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	var contactReq model.ContactRequest
	if err := json.Unmarshal([]byte(body), &contactReq); err != nil {
		return h.errorResponse(400, "Invalid request body", headers), nil
	}

	if err := h.contactService.ProcessContactRequest(ctx, &contactReq); err != nil {
		log.Printf("Error processing contact request: %v", err)
		return h.errorResponse(400, "Invalid request", headers), nil
	}

	// Send notification
	payload := &model.NotificationPayload{
		Type: "contact",
		Data: map[string]interface{}{
			"name":    contactReq.Name,
			"email":   contactReq.Email,
			"message": contactReq.Message,
		},
		Source:    "resume-website",
		Timestamp: contactReq.Timestamp,
	}
	go h.notificationService.SendEmailNotification(context.Background(), payload)
	go h.notificationService.SendSMSNotification(context.Background(), payload)

	response := model.APIResponse{Success: true, Message: "Message sent successfully"}
	responseBody, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       string(responseBody),
	}, nil
}

func (h *APIHandler) errorResponse(statusCode int, message string, headers map[string]string) events.APIGatewayProxyResponse {
	headers["Set-Cookie"] = "" // Clear session cookie on error
	response := model.APIResponse{Error: message, Success: false}
	body, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    headers,
		Body:       string(body),
	}
}
