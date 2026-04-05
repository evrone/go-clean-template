package v1

import (
	"fmt"

	"github.com/evrone/go-clean-template/internal/controller/nats_rpc/v1/request"
	"github.com/evrone/go-clean-template/pkg/jwt"
	"github.com/goccy/go-json"
	"github.com/nats-io/nats.go"
)

func extractUserID(msg *nats.Msg, jwtManager *jwt.Manager) (userID string, data json.RawMessage, err error) {
	var req request.AuthenticatedRequest

	err = json.Unmarshal(msg.Data, &req)
	if err != nil {
		return "", nil, fmt.Errorf("invalid request format: %w", err)
	}

	userID, err = jwtManager.ParseToken(req.Token)
	if err != nil {
		return "", nil, fmt.Errorf("invalid or expired token: %w", err)
	}

	return userID, req.Data, nil
}
