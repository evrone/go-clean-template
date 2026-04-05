package request

import "github.com/goccy/go-json"

// AuthenticatedRequest is the envelope for all authenticated RPC calls.
type AuthenticatedRequest struct {
	Token string          `json:"token" validate:"required"`
	Data  json.RawMessage `json:"data"`
}
