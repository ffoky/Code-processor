package types

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// PostTaskHandlerRequest represents task post request
// swagger:param parameters PostHandler
type PostUserRegistrationHandlerRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func CreatePostUserRegistrationHandlerRequest(r *http.Request) (*PostUserRegistrationHandlerRequest, error) {
	var req PostUserRegistrationHandlerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("error while decoding json: %v", err)
	}
	return &req, nil
}

type PostUserLoginHandlerRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func CreatePostUserLoginHandlerRequest(r *http.Request) (*PostUserLoginHandlerRequest, error) {
	var req PostUserLoginHandlerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("error while decoding json: %v", err)
	}
	return &req, nil
}

// PostUserRegistrationHandlerResponse represents response with created user tid
// swagger:response PostUserRegistrationHandlerResponse
type PostUserRegistrationHandlerResponse struct {
	StatusOK int `json:"status-ok"`
}

type PostUserLoginHandlerResponse struct {
	SessionID string `json:"session-id"`
}
