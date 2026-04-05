package v1

import (
	"context"
	"fmt"

	"github.com/evrone/go-clean-template/internal/controller/nats_rpc/v1/request"
	"github.com/evrone/go-clean-template/internal/controller/nats_rpc/v1/response"
	"github.com/evrone/go-clean-template/pkg/nats/nats_rpc/server"
	"github.com/goccy/go-json"
	"github.com/nats-io/nats.go"
)

func (r *V1) register() server.CallHandler {
	return func(msg *nats.Msg) (any, error) {
		var req request.Register

		err := json.Unmarshal(msg.Data, &req)
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - register")

			return nil, fmt.Errorf("nats_rpc - V1 - register - json.Unmarshal: %w", err)
		}

		if err = r.v.Struct(req); err != nil {
			return nil, fmt.Errorf("nats_rpc - V1 - register - validation: %w", err)
		}

		user, err := r.u.Register(context.Background(), req.Username, req.Email, req.Password)
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - register")

			return nil, fmt.Errorf("nats_rpc - V1 - register: %w", err)
		}

		return user, nil
	}
}

func (r *V1) login() server.CallHandler {
	return func(msg *nats.Msg) (any, error) {
		var req request.Login

		err := json.Unmarshal(msg.Data, &req)
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - login")

			return nil, fmt.Errorf("nats_rpc - V1 - login - json.Unmarshal: %w", err)
		}

		if err = r.v.Struct(req); err != nil {
			return nil, fmt.Errorf("nats_rpc - V1 - login - validation: %w", err)
		}

		token, err := r.u.Login(context.Background(), req.Email, req.Password)
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - login")

			return nil, fmt.Errorf("nats_rpc - V1 - login: %w", err)
		}

		return response.Token{Token: token}, nil
	}
}
