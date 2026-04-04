package v1

import (
	"context"
	"fmt"

	"github.com/evrone/go-clean-template/internal/controller/nats_rpc/v1/request"
	"github.com/evrone/go-clean-template/internal/controller/nats_rpc/v1/response"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/pkg/nats/nats_rpc/server"
	"github.com/goccy/go-json"
	"github.com/nats-io/nats.go"
)

func (r *V1) createTask() server.CallHandler {
	return func(msg *nats.Msg) (any, error) {
		userID, data, err := extractUserID(msg, r.j)
		if err != nil {
			return nil, fmt.Errorf("nats_rpc - V1 - createTask - auth: %w", err)
		}

		var req request.CreateTask

		err = json.Unmarshal(data, &req)
		if err != nil {
			return nil, fmt.Errorf("nats_rpc - V1 - createTask - json.Unmarshal: %w", err)
		}

		if err = r.v.Struct(req); err != nil {
			return nil, fmt.Errorf("nats_rpc - V1 - createTask - validation: %w", err)
		}

		task, err := r.tk.Create(context.Background(), userID, req.Title, req.Description)
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - createTask")

			return nil, fmt.Errorf("nats_rpc - V1 - createTask: %w", err)
		}

		return task, nil
	}
}

func (r *V1) getTask() server.CallHandler {
	return func(msg *nats.Msg) (any, error) {
		userID, data, err := extractUserID(msg, r.j)
		if err != nil {
			return nil, fmt.Errorf("nats_rpc - V1 - getTask - auth: %w", err)
		}

		var req request.GetTask

		err = json.Unmarshal(data, &req)
		if err != nil {
			return nil, fmt.Errorf("nats_rpc - V1 - getTask - json.Unmarshal: %w", err)
		}

		if err = r.v.Struct(req); err != nil {
			return nil, fmt.Errorf("nats_rpc - V1 - getTask - validation: %w", err)
		}

		task, err := r.tk.Get(context.Background(), userID, req.ID)
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - getTask")

			return nil, fmt.Errorf("nats_rpc - V1 - getTask: %w", err)
		}

		return task, nil
	}
}

func (r *V1) listTasks() server.CallHandler {
	return func(msg *nats.Msg) (any, error) {
		userID, data, err := extractUserID(msg, r.j)
		if err != nil {
			return nil, fmt.Errorf("nats_rpc - V1 - listTasks - auth: %w", err)
		}

		var req request.ListTasks

		err = json.Unmarshal(data, &req)
		if err != nil {
			return nil, fmt.Errorf("nats_rpc - V1 - listTasks - json.Unmarshal: %w", err)
		}

		if err = r.v.Struct(req); err != nil {
			return nil, fmt.Errorf("nats_rpc - V1 - listTasks - validation: %w", err)
		}

		var status *entity.TaskStatus

		if req.Status != "" {
			s := entity.TaskStatus(req.Status)
			status = &s
		}

		tasks, total, err := r.tk.List(context.Background(), userID, status, req.Limit, req.Offset)
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - listTasks")

			return nil, fmt.Errorf("nats_rpc - V1 - listTasks: %w", err)
		}

		return response.TaskList{Tasks: tasks, Total: total}, nil
	}
}

func (r *V1) updateTask() server.CallHandler {
	return func(msg *nats.Msg) (any, error) {
		userID, data, err := extractUserID(msg, r.j)
		if err != nil {
			return nil, fmt.Errorf("nats_rpc - V1 - updateTask - auth: %w", err)
		}

		var req request.UpdateTask

		err = json.Unmarshal(data, &req)
		if err != nil {
			return nil, fmt.Errorf("nats_rpc - V1 - updateTask - json.Unmarshal: %w", err)
		}

		if err = r.v.Struct(req); err != nil {
			return nil, fmt.Errorf("nats_rpc - V1 - updateTask - validation: %w", err)
		}

		task, err := r.tk.Update(context.Background(), userID, req.ID, req.Title, req.Description)
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - updateTask")

			return nil, fmt.Errorf("nats_rpc - V1 - updateTask: %w", err)
		}

		return task, nil
	}
}

func (r *V1) transitionTask() server.CallHandler {
	return func(msg *nats.Msg) (any, error) {
		userID, data, err := extractUserID(msg, r.j)
		if err != nil {
			return nil, fmt.Errorf("nats_rpc - V1 - transitionTask - auth: %w", err)
		}

		var req request.TransitionTask

		err = json.Unmarshal(data, &req)
		if err != nil {
			return nil, fmt.Errorf("nats_rpc - V1 - transitionTask - json.Unmarshal: %w", err)
		}

		if err = r.v.Struct(req); err != nil {
			return nil, fmt.Errorf("nats_rpc - V1 - transitionTask - validation: %w", err)
		}

		task, err := r.tk.Transition(context.Background(), userID, req.ID, entity.TaskStatus(req.Status))
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - transitionTask")

			return nil, fmt.Errorf("nats_rpc - V1 - transitionTask: %w", err)
		}

		return task, nil
	}
}

func (r *V1) deleteTask() server.CallHandler {
	return func(msg *nats.Msg) (any, error) {
		userID, data, err := extractUserID(msg, r.j)
		if err != nil {
			return nil, fmt.Errorf("nats_rpc - V1 - deleteTask - auth: %w", err)
		}

		var req request.DeleteTask

		err = json.Unmarshal(data, &req)
		if err != nil {
			return nil, fmt.Errorf("nats_rpc - V1 - deleteTask - json.Unmarshal: %w", err)
		}

		if err = r.v.Struct(req); err != nil {
			return nil, fmt.Errorf("nats_rpc - V1 - deleteTask - validation: %w", err)
		}

		err = r.tk.Delete(context.Background(), userID, req.ID)
		if err != nil {
			r.l.Error(err, "nats_rpc - V1 - deleteTask")

			return nil, fmt.Errorf("nats_rpc - V1 - deleteTask: %w", err)
		}

		return response.DeleteStatus{Status: "deleted"}, nil
	}
}
