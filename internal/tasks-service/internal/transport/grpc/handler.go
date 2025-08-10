package grpc

import (
	"context"
	"errors"

	taskspb "github.com/blastuha/test-service-proto/gen/task"
	"github.com/your-org/tasks-service/internal/tasks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	taskspb.UnimplementedTasksServiceServer
	svc    tasks.TasksService
	client Client
}

func NewHandler(svc tasks.TasksService, client Client) *Handler {
	return &Handler{svc: svc, client: client}
}

func (h *Handler) CreateTask(ctx context.Context, req *taskspb.TaskCreateRequest) (*taskspb.TaskResponse, error) {
	if req.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "user id must be > 0")
	}

	if _, err := h.client.GetUser(ctx, req.GetUserId()); err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user with id %d not found", req.GetUserId())
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	dm, err := h.svc.CreateTask(req.Title, req.IsDone, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create task: %v", err)
	}

	response := &taskspb.TaskResponse{Task: &taskspb.Task{Id: dm.ID, Title: dm.Task, IsDone: dm.IsDone}}
	return response, nil
}
