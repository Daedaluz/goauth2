package oidc

import (
	"encoding/json"
	"fmt"
	"io"
)

type ErrorCode string

// Standard error codes
const (
	ErrInvalidRequest       = ErrorCode("invalid_request")
	ErrInvalidClient        = ErrorCode("invalid_client")
	ErrInvalidGrant         = ErrorCode("invalid_grant")
	ErrUnauthorizedClient   = ErrorCode("unauthorized_client")
	ErrUnsupportedGrantType = ErrorCode("unsupported_grant_type")
	ErrInvalidScope         = ErrorCode("invalid_scope")
)

// CIBA error codes
const (
	ErrSlowDown             = ErrorCode("slow_down")
	ErrAuthorizationPending = ErrorCode("authorization_pending")
	ErrAuthorizationViewed  = ErrorCode("authorization_viewed")
	ErrExpiredToken         = ErrorCode("expired_token")
	ErrAccessDenied         = ErrorCode("access_denied")
	ErrTransactionFailed    = ErrorCode("transaction_failed")
)

// ErrorResponse represents an error response from the server
type ErrorResponse struct {
	Err              ErrorCode `json:"error"`
	ErrorDescription string    `json:"error_description"`
	ErrorURI         string    `json:"error_uri"`

	Other json.RawMessage
}

func (e ErrorResponse) Error() string {
	if e.ErrorDescription != "" {
		return fmt.Sprintf("%s; %s", e.Err, e.ErrorDescription)
	}
	return string(e.Err)
}

func ParseError(r io.Reader) error {
	var e ErrorResponse
	if err := json.NewDecoder(r).Decode(&e); err != nil {
		return err
	}
	return &e
}
