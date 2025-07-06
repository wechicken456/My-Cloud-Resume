package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/internal/model"
	"main/internal/service"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

type APIHandler struct {
	sessionService      *service.SessionService
	visitorService      *service.VisitorService
	likesService        *service.LikeService
	contactService      *service.ContactService
	notificationService *service.NotificationService
}

func NewAPIHandler(
	sessionService *service.SessionService,
	visitorService *service.VisitorService,
	likesService *service.LikeService,
	contactService *service.ContactService,
	notificationService *service.NotificationService,
) *APIHandler {
	return &APIHandler{
		sessionService:      sessionService,
		visitorService:      visitorService,
		likesService:        likesService,
		contactService:      contactService,
		notificationService: notificationService,
	}
}

var sessionIDCookieName string = "session_id"

func (h *APIHandler) HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Extract session ID from cookie
	sessionID := h.extractSessionID(req.Headers["cookie"])

	switch req.Resource {
	case "/api/session":
		return h.handleGetSession(ctx, req, sessionID)
	case "/api/getVisitorCount":
		return h.handleGetVisitorCount(ctx, req)
	case "/api/incrementVisitorCount":
		return h.handleIncrementVisitorCount(ctx, req, sessionID)
	case "/api/getLikeCount":
		return h.handleGetLikeCount(ctx, req)
	case "/api/toggleLike":
		return h.handleToggleLike(ctx, req, sessionID)
	case "/api/contact":
		return h.handleContact(ctx, req)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       `{"error": "Not found", "success": false}`,
		}, nil
	}
}

func (h *APIHandler) handleGetSession(ctx context.Context, req events.APIGatewayProxyRequest, sessionID string) (events.APIGatewayProxyResponse, error) {
	var headers map[string]string = map[string]string{}

	headers["Content-Type"] = "application/json"
	headers["Access-Control-Allow-Origin"] = "https://www.pwnph0fun.com"
	headers["Access-Control-Allow-Methods"] = "GET,OPTIONS"
	headers["Access-Control-Allow-Headers"] = "Content-Type,Cookie"
	headers["Access-Control-Allow-Credentials"] = "true"

	// Handle CORS preflight request
	if req.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 204,
			Headers:    headers,
			Body:       "",
		}, nil
	}

	if req.HTTPMethod != "GET" {
		return h.errorResponse(405, "Method not allowed", headers), nil
	}

	session, isNewSession, err := h.sessionService.GetOrCreateSession(ctx, sessionID)
	if err != nil {
		log.Printf("Error getting session: %v", err)
		return h.errorResponse(500, "Session error", headers), nil
	}

	response := map[string]any{
		"has_visited": session.HasVisited,
		"has_liked":   session.HasLiked,
	}

	body, _ := json.Marshal(response)

	// Set session cookie if new session was created
	if isNewSession {
		headers["Set-Cookie"] = fmt.Sprintf("%s=%s; HttpOnly; Secure; SameSite=Strict; Max-Age=86400; Path=/", sessionIDCookieName, session.SessionID)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       string(body),
	}, nil
}

func (h *APIHandler) handleGetVisitorCount(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var headers map[string]string = map[string]string{}

	headers["Content-Type"] = "application/json"
	headers["Access-Control-Allow-Origin"] = "*"
	headers["Access-Control-Allow-Methods"] = "GET,OPTIONS"
	headers["Access-Control-Allow-Headers"] = "Content-Type,Cookie"
	headers["Access-Control-Allow-Credentials"] = "true"

	// Handle CORS preflight request
	if req.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 204,
			Headers:    headers,
			Body:       "",
		}, nil
	}

	if req.HTTPMethod != "GET" {
		return h.errorResponse(405, "Method not allowed", headers), nil
	}

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

func (h *APIHandler) handleIncrementVisitorCount(ctx context.Context, req events.APIGatewayProxyRequest, sessionID string) (events.APIGatewayProxyResponse, error) {
	var headers map[string]string = map[string]string{}

	headers["Content-Type"] = "application/json"
	headers["Access-Control-Allow-Origin"] = "https://www.pwnph0fun.com"
	headers["Access-Control-Allow-Methods"] = "POST,OPTIONS"
	headers["Access-Control-Allow-Headers"] = "Content-Type,Cookie"
	headers["Access-Control-Allow-Credentials"] = "true"

	// Handle CORS preflight request
	if req.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 204,
			Headers:    headers,
			Body:       "",
		}, nil
	}

	if req.HTTPMethod != "POST" {
		return h.errorResponse(405, "Method not allowed", headers), nil
	}

	// Validate session exists before proceeding
	session, err := h.sessionService.ValidateSession(ctx, sessionID)
	if err != nil {
		log.Printf("Error validating session: %v", err)
		return h.errorResponse(500, "Session error", headers), nil
	}

	if session == nil {
		return h.errorResponse(401, "Invalid session", headers), nil
	}

	count, status, err := h.visitorService.IncrementVisitorCount(ctx, session)
	if err != nil {
		log.Printf("Error incrementing count: %v", err)
		return h.errorResponse(500, "Database error", headers), nil
	}

	// Update session if visitor count was incremented
	if status == "incremented" {
		session.HasVisited = true
		err = h.sessionService.UpdateSession(ctx, session)
		if err != nil {
			log.Printf("Error updating session: %v", err)
			// Don't fail the request if session update fails
		}
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

func (h *APIHandler) handleGetLikeCount(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var headers map[string]string = map[string]string{}

	headers["Content-Type"] = "application/json"
	headers["Access-Control-Allow-Origin"] = "*"
	headers["Access-Control-Allow-Methods"] = "GET,OPTIONS"
	headers["Access-Control-Allow-Headers"] = "Content-Type,Cookie"
	headers["Access-Control-Allow-Credentials"] = "true"

	if req.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 204,
			Headers:    headers,
			Body:       "",
		}, nil
	}

	if req.HTTPMethod != "GET" {
		return h.errorResponse(405, "Method not allowed", headers), nil
	}

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

func (h *APIHandler) handleToggleLike(ctx context.Context, req events.APIGatewayProxyRequest, sessionID string) (events.APIGatewayProxyResponse, error) {
	var headers map[string]string = map[string]string{}

	headers["Content-Type"] = "application/json"
	headers["Access-Control-Allow-Origin"] = "https://www.pwnph0fun.com"
	headers["Access-Control-Allow-Methods"] = "POST,OPTIONS"
	headers["Access-Control-Allow-Headers"] = "Content-Type,Cookie"
	headers["Access-Control-Allow-Credentials"] = "true"

	// Handle CORS preflight request
	if req.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 204,
			Headers:    headers,
			Body:       "",
		}, nil
	}

	if req.HTTPMethod != "POST" {
		return h.errorResponse(405, "Method not allowed", headers), nil
	}

	// Validate session exists before proceeding
	session, err := h.sessionService.ValidateSession(ctx, sessionID)
	if err != nil {
		log.Printf("Error validating session: %v", err)
		return h.errorResponse(500, "Session error", headers), nil
	}

	if session == nil {
		return h.errorResponse(401, "Invalid session", headers), nil
	}

	count, action, err := h.likesService.ToggleLike(ctx, session)
	if err != nil {
		log.Printf("Error toggling like: %v", err)
		return h.errorResponse(500, "Database error", headers), nil
	}

	// Update session with new like status
	err = h.sessionService.UpdateSession(ctx, session)
	if err != nil {
		log.Printf("Error updating session: %v", err)
		// Don't fail the request if session update fails
	}

	// Send notification if this is a new like
	if action == "liked" {
		payload := &model.NotificationPayload{
			Type:      "like",
			Data:      map[string]any{},
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

func (h *APIHandler) handleContact(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var headers map[string]string = map[string]string{}

	headers["Content-Type"] = "application/json"
	headers["Access-Control-Allow-Origin"] = "https://www.pwnph0fun.com"
	headers["Access-Control-Allow-Methods"] = "POST,OPTIONS"
	headers["Access-Control-Allow-Headers"] = "Content-Type,Cookie"
	headers["Access-Control-Allow-Credentials"] = "true"

	// Handle CORS preflight request
	if req.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 204,
			Headers:    headers,
			Body:       "",
		}, nil
	}

	if req.HTTPMethod != "POST" {
		return h.errorResponse(405, "Method not allowed", headers), nil
	}

	body := req.Body
	var contactReq model.ContactRequest
	if err := json.Unmarshal([]byte(body), &contactReq); err != nil {
		return h.errorResponse(400, "Invalid request body", headers), nil
	}

	if err := h.contactService.ProcessContactRequest(ctx, &contactReq); err != nil {
		log.Printf("Error processing contact request: %v", err)
		return h.errorResponse(400, fmt.Sprintf("Invalid request %v", err), headers), nil
	}

	// Send notification
	payload := &model.NotificationPayload{
		Type: "contact",
		Data: map[string]any{
			"name":    contactReq.Name,
			"email":   contactReq.Email,
			"message": contactReq.Message,
		},
		Source: "resume-website",
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

func (h *APIHandler) extractSessionID(cookieHeader string) string {
	if cookieHeader == "" {
		return ""
	}
	// Parse cookie header to extract session_id
	cookies := parseCookies(cookieHeader)
	return cookies[sessionIDCookieName]
}

func parseCookies(cookieHeader string) map[string]string {
	cookies := make(map[string]string)
	if cookieHeader == "" {
		return cookies
	}

	parts := strings.SplitSeq(cookieHeader, ";")
	for part := range parts {
		part = strings.TrimSpace(part)
		if eq := strings.Index(part, "="); eq >= 0 {
			key := strings.TrimSpace(part[:eq])
			value := strings.TrimSpace(part[eq+1:])
			cookies[key] = value
		}
	}
	return cookies
}

func (h *APIHandler) errorResponse(statusCode int, message string, headers map[string]string) events.APIGatewayProxyResponse {
	response := model.APIResponse{Error: message, Success: false}
	body, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    headers,
		Body:       string(body),
	}
}
