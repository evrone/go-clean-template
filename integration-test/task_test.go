package integration_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"

	protov1 "github.com/evrone/go-clean-template/docs/proto/v1"
	natsClient "github.com/evrone/go-clean-template/pkg/nats/nats_rpc/client"
	rmqClient "github.com/evrone/go-clean-template/pkg/rabbitmq/rmq_rpc/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	statusTodo       = "todo"
	statusInProgress = "in_progress"
)

type taskResponse struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// httpCreateTask is a helper that creates a task via HTTP and returns the parsed response.
func httpCreateTask(t *testing.T, token, title, description string) taskResponse {
	t.Helper()

	createBody := fmt.Sprintf(`{"title":%q,"description":%q}`, title, description)

	ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
	defer cancel()

	resp, err := doAuthenticatedRequest(ctx, http.MethodPost, basePathV1+"/tasks/", bytes.NewBufferString(createBody), token)
	if err != nil {
		t.Fatalf("Create task: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Create task: expected 201, got %d", resp.StatusCode)
	}

	return parseJSON[taskResponse](t, resp)
}

// httpTransitionTask sends a PATCH request to transition a task's status.
func httpTransitionTask(t *testing.T, token, id, status string) (*http.Response, error) {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
	defer cancel()

	body := fmt.Sprintf(`{"status":%q}`, status)

	return doAuthenticatedRequest(ctx, http.MethodPatch, basePathV1+"/tasks/"+id+"/status", bytes.NewBufferString(body), token)
}

// HTTP: create task and verify fields.
func TestHTTPTaskCreateV1(t *testing.T) {
	token := registerAndLogin(t)
	created := httpCreateTask(t, token, "integration task", "test description")

	if created.ID == "" {
		t.Fatal("expected non-empty id")
	}

	if created.Status != statusTodo {
		t.Errorf("expected status 'todo', got %q", created.Status)
	}
}

// HTTP: get task by ID.
func TestHTTPTaskGetV1(t *testing.T) {
	token := registerAndLogin(t)
	created := httpCreateTask(t, token, "get task", "test description")

	ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
	defer cancel()

	resp, err := doAuthenticatedRequest(ctx, http.MethodGet, basePathV1+"/tasks/"+created.ID, http.NoBody, token)
	if err != nil {
		t.Fatalf("Get task: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	got := parseJSON[taskResponse](t, resp)

	if got.ID != created.ID {
		t.Errorf("expected id %q, got %q", created.ID, got.ID)
	}
}

// HTTP: list tasks.
func TestHTTPTaskListV1(t *testing.T) {
	token := registerAndLogin(t)
	httpCreateTask(t, token, "list task", "test description")

	ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
	defer cancel()

	resp, err := doAuthenticatedRequest(ctx, http.MethodGet, basePathV1+"/tasks/?limit=10&offset=0", http.NoBody, token)
	if err != nil {
		t.Fatalf("List tasks: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	type listResponse struct {
		Tasks []taskResponse `json:"tasks"`
		Total int            `json:"total"`
	}

	listed := parseJSON[listResponse](t, resp)

	if listed.Total < 1 {
		t.Errorf("expected total >= 1, got %d", listed.Total)
	}
}

// HTTP: update task.
func TestHTTPTaskUpdateV1(t *testing.T) {
	token := registerAndLogin(t)
	created := httpCreateTask(t, token, "update task", "test description")

	updateBody := `{"title":"updated title","description":"updated description"}`

	ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
	defer cancel()

	resp, err := doAuthenticatedRequest(ctx, http.MethodPut, basePathV1+"/tasks/"+created.ID, bytes.NewBufferString(updateBody), token)
	if err != nil {
		t.Fatalf("Update task: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	updated := parseJSON[taskResponse](t, resp)

	if updated.Title != "updated title" {
		t.Errorf("expected title 'updated title', got %q", updated.Title)
	}
}

// HTTP: delete task.
func TestHTTPTaskDeleteV1(t *testing.T) {
	token := registerAndLogin(t)
	created := httpCreateTask(t, token, "delete task", "test description")

	ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
	defer cancel()

	resp, err := doAuthenticatedRequest(ctx, http.MethodDelete, basePathV1+"/tasks/"+created.ID, http.NoBody, token)
	if err != nil {
		t.Fatalf("Delete task: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected 204, got %d", resp.StatusCode)
	}
}

// HTTP: valid task status transition chain.
func TestHTTPTaskStatusTransitionValidV1(t *testing.T) {
	token := registerAndLogin(t)
	task := httpCreateTask(t, token, "transition task", "testing transitions")

	if task.Status != statusTodo {
		t.Fatalf("Expected initial status 'todo', got %q", task.Status)
	}

	// Transition: from "todo" to "in_progress".
	resp, err := httpTransitionTask(t, token, task.ID, statusInProgress)
	if err != nil {
		t.Fatalf("Transition to in_progress: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	transitioned := parseJSON[taskResponse](t, resp)

	if transitioned.Status != statusInProgress {
		t.Errorf("Expected status 'in_progress', got %q", transitioned.Status)
	}

	// Transition: in_progress to done.
	resp2, err := httpTransitionTask(t, token, task.ID, "done")
	if err != nil {
		t.Fatalf("Transition to done: %v", err)
	}

	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp2.StatusCode)
	}

	transitioned2 := parseJSON[taskResponse](t, resp2)

	if transitioned2.Status != "done" {
		t.Errorf("Expected status 'done', got %q", transitioned2.Status)
	}
}

// HTTP: invalid task status transition is rejected.
func TestHTTPTaskStatusTransitionInvalidV1(t *testing.T) {
	token := registerAndLogin(t)
	task := httpCreateTask(t, token, "transition task", "testing transitions")

	resp, err := httpTransitionTask(t, token, task.ID, "done")
	if err != nil {
		t.Fatalf("Transition to done: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

// HTTP: task error cases.
func TestHTTPTaskErrorsV1(t *testing.T) {
	t.Run("no token returns 401", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
		defer cancel()

		body := `{"title":"unauthorized task","description":"should fail"}`

		resp, err := doWebRequestWithTimeout(ctx, http.MethodPost, basePathV1+"/tasks/", bytes.NewBufferString(body))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", resp.StatusCode)
		}
	})

	t.Run("get non-existent task returns 404", func(t *testing.T) {
		token := registerAndLogin(t)

		ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
		defer cancel()

		resp, err := doAuthenticatedRequest(ctx, http.MethodGet, basePathV1+"/tasks/00000000-0000-0000-0000-000000000000", http.NoBody, token)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("create with invalid body returns 400", func(t *testing.T) {
		token := registerAndLogin(t)

		ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
		defer cancel()

		body := `{"description":"missing title"}`

		resp, err := doAuthenticatedRequest(ctx, http.MethodPost, basePathV1+"/tasks/", bytes.NewBufferString(body), token)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})
}

// gRPC: create and verify task.
func TestGRPCTaskCreateV1(t *testing.T) {
	token := registerAndLoginGRPC(t)

	grpcConn, err := grpc.NewClient(grpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc.NewClient: %v", err)
	}

	defer func() {
		if cerr := grpcConn.Close(); cerr != nil {
			t.Fatalf("grpcConn.Close: %v", cerr)
		}
	}()

	taskClient := protov1.NewTaskServiceClient(grpcConn)
	authCtx := grpcAuthCtx(t, token)

	createResp, err := taskClient.CreateTask(authCtx, &protov1.CreateTaskRequest{
		Title:       "grpc task",
		Description: "grpc description",
	})
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}

	if createResp.GetId() == "" {
		t.Fatal("expected non-empty id")
	}

	if createResp.GetStatus() != statusTodo {
		t.Errorf("expected status 'todo', got %q", createResp.GetStatus())
	}
}

// gRPC: get and list tasks.
func TestGRPCTaskGetListV1(t *testing.T) {
	token := registerAndLoginGRPC(t)

	grpcConn, err := grpc.NewClient(grpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc.NewClient: %v", err)
	}

	defer func() {
		if cerr := grpcConn.Close(); cerr != nil {
			t.Fatalf("grpcConn.Close: %v", cerr)
		}
	}()

	taskClient := protov1.NewTaskServiceClient(grpcConn)
	authCtx := grpcAuthCtx(t, token)

	createResp, err := taskClient.CreateTask(authCtx, &protov1.CreateTaskRequest{
		Title:       "grpc get-list task",
		Description: "grpc description",
	})
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}

	taskID := createResp.GetId()

	getResp, err := taskClient.GetTask(authCtx, &protov1.GetTaskRequest{Id: taskID})
	if err != nil {
		t.Fatalf("GetTask: %v", err)
	}

	if getResp.GetId() != taskID {
		t.Errorf("expected id %q, got %q", taskID, getResp.GetId())
	}

	listResp, err := taskClient.ListTasks(authCtx, &protov1.ListTasksRequest{
		Status: "",
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("ListTasks: %v", err)
	}

	if listResp.GetTotal() < 1 {
		t.Errorf("expected total >= 1, got %d", listResp.GetTotal())
	}
}

// gRPC: update, transition, and delete task.
func TestGRPCTaskUpdateTransitionDeleteV1(t *testing.T) {
	token := registerAndLoginGRPC(t)

	grpcConn, err := grpc.NewClient(grpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc.NewClient: %v", err)
	}

	defer func() {
		if cerr := grpcConn.Close(); cerr != nil {
			t.Fatalf("grpcConn.Close: %v", cerr)
		}
	}()

	taskClient := protov1.NewTaskServiceClient(grpcConn)
	authCtx := grpcAuthCtx(t, token)

	createResp, err := taskClient.CreateTask(authCtx, &protov1.CreateTaskRequest{
		Title:       "grpc update task",
		Description: "grpc description",
	})
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}

	taskID := createResp.GetId()

	updateResp, err := taskClient.UpdateTask(authCtx, &protov1.UpdateTaskRequest{
		Id:          taskID,
		Title:       "updated grpc task",
		Description: "updated grpc description",
	})
	if err != nil {
		t.Fatalf("UpdateTask: %v", err)
	}

	if updateResp.GetTitle() != "updated grpc task" {
		t.Errorf("expected title 'updated grpc task', got %q", updateResp.GetTitle())
	}

	transResp, err := taskClient.TransitionTask(authCtx, &protov1.TransitionTaskRequest{
		Id:     taskID,
		Status: statusInProgress,
	})
	if err != nil {
		t.Fatalf("TransitionTask: %v", err)
	}

	if transResp.GetStatus() != statusInProgress {
		t.Errorf("expected status 'in_progress', got %q", transResp.GetStatus())
	}

	_, err = taskClient.DeleteTask(authCtx, &protov1.DeleteTaskRequest{Id: taskID})
	if err != nil {
		t.Fatalf("DeleteTask: %v", err)
	}
}

// RabbitMQ RPC: create + list tasks smoke test.
func TestRMQTaskV1(t *testing.T) {
	token := registerAndLoginRMQ(t)

	client, err := rmqClient.New(rmqURL, rpcServerExchange, rpcClientExchange)
	if err != nil {
		t.Fatalf("rmqClient.New: %v", err)
	}

	defer func() {
		if serr := client.Shutdown(); serr != nil {
			t.Fatalf("client.Shutdown: %v", serr)
		}
	}()

	// Create task.
	createPayload := authenticatedPayload(token, map[string]string{
		"title":       "rmq task",
		"description": "rmq description",
	})

	var createResp struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}

	err = client.RemoteCall("v1.task.create", createPayload, &createResp)
	if err != nil {
		t.Fatalf("v1.task.create: %v", err)
	}

	if createResp.ID == "" {
		t.Error("Expected non-empty task ID from create")
	}

	// List tasks.
	listPayload := authenticatedPayload(token, map[string]any{
		"status": "",
		"limit":  10,
		"offset": 0,
	})

	var listResp struct {
		Tasks []struct {
			ID string `json:"id"`
		} `json:"tasks"`
		Total int `json:"total"`
	}

	err = client.RemoteCall("v1.task.list", listPayload, &listResp)
	if err != nil {
		t.Fatalf("v1.task.list: %v", err)
	}

	if listResp.Total < 1 {
		t.Errorf("Expected total >= 1, got %d", listResp.Total)
	}
}

// NATS RPC: create + list tasks smoke test.
func TestNATSTaskV1(t *testing.T) {
	token := registerAndLoginNATS(t)

	client, err := natsClient.New(natsURL, rpcServerExchange)
	if err != nil {
		t.Fatalf("natsClient.New: %v", err)
	}

	defer func() {
		if serr := client.Shutdown(); serr != nil {
			t.Fatalf("client.Shutdown: %v", serr)
		}
	}()

	// Create task.
	createPayload := authenticatedPayload(token, map[string]string{
		"title":       "nats task",
		"description": "nats description",
	})

	var createResp struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}

	err = client.RemoteCall("v1.task.create", createPayload, &createResp)
	if err != nil {
		t.Fatalf("v1.task.create: %v", err)
	}

	if createResp.ID == "" {
		t.Error("Expected non-empty task ID from create")
	}

	// List tasks.
	listPayload := authenticatedPayload(token, map[string]any{
		"status": "",
		"limit":  10,
		"offset": 0,
	})

	var listResp struct {
		Tasks []struct {
			ID string `json:"id"`
		} `json:"tasks"`
		Total int `json:"total"`
	}

	err = client.RemoteCall("v1.task.list", listPayload, &listResp)
	if err != nil {
		t.Fatalf("v1.task.list: %v", err)
	}

	if listResp.Total < 1 {
		t.Errorf("Expected total >= 1, got %d", listResp.Total)
	}
}
