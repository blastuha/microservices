package grpc

import (
	"context"
	"errors"

	taskspb "github.com/blastuha/test-service-proto/gen/task"
	"github.com/your-org/tasks-service/internal/tasks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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
		return nil, status.Errorf(codes.InvalidArgument, "user id must be > 0")
	}

	if _, err := h.client.GetUser(ctx, req.GetUserId()); err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user with id %d not found", req.GetUserId())
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	dm, err := h.svc.CreateTask(req.GetTitle(), req.GetIsDone(), req.GetUserId())
	if err != nil {
		if errors.Is(err, tasks.ErrInvalidInput) {
			return nil, status.Error(codes.InvalidArgument, "title must not be empty")
		}
		return nil, status.Errorf(codes.Internal, "failed to create task: %v", err)
	}

	response := &taskspb.TaskResponse{Task: &taskspb.Task{Id: dm.ID, Title: dm.Task, IsDone: dm.IsDone, UserId: dm.UserID}}
	return response, nil
}

func (h *Handler) GetTaskList(ctx context.Context, req *emptypb.Empty) (*taskspb.TaskListResponse, error) {
	tasksList, err := h.svc.GetAllTasks()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get list of tasks")
	}

	out := make([]*taskspb.Task, 0, len(tasksList))
	for _, t := range tasksList {
		out = append(out, &taskspb.Task{Id: t.ID, Title: t.Task, IsDone: t.IsDone, UserId: t.UserID})
	}

	response := &taskspb.TaskListResponse{Tasks: out}
	return response, nil
}

func (h *Handler) UpdateTask(ctx context.Context, req *taskspb.TaskUpdateRequest) (*taskspb.TaskResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	id := req.GetId()
	if id == 0 {
		return nil, status.Error(codes.InvalidArgument, "task id must be > 0")
	}

	dm, err := h.svc.UpdateTask(req.GetTitle(), req.GetIsDone(), id)
	if err != nil {
		switch {
		case errors.Is(err, tasks.ErrInvalidInput):
			return nil, status.Error(codes.InvalidArgument, "title must not be empty")
		case errors.Is(err, tasks.ErrTaskNotFound):
			return nil, status.Error(codes.NotFound, "task not found")
		default:
			return nil, status.Errorf(codes.Internal, "failed to update task: %v", err)
		}
	}

	return &taskspb.TaskResponse{
		Task: &taskspb.Task{
			Id:     dm.ID,
			Title:  dm.Task,
			IsDone: dm.IsDone,
			UserId: dm.UserID,
		},
	}, nil
}

func (h *Handler) DeleteTask(ctx context.Context, req *taskspb.TaskDeleteRequest) (*emptypb.Empty, error) {
	id := req.GetId()
	if id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id must be > 0")
	}

	if err := h.svc.DeleteTask(id); err != nil {
		if errors.Is(err, tasks.ErrTaskNotFound) {
			return nil, status.Errorf(codes.NotFound, "task not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete task: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (h *Handler) ListTasksByUser(ctx context.Context, req *taskspb.ListTasksByUserRequest) (*taskspb.TaskListResponse, error) {
	if req.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "id must be > 0")
	}

	tasks, err := h.svc.ListTasksByUser(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get list of tasks by user_id: %d", req.UserId)
	}

	out := make([]*taskspb.Task, 0, len(tasks))

	for _, t := range tasks {
		pbTask := &taskspb.Task{Id: t.ID, IsDone: t.IsDone, UserId: t.UserID, Title: t.Task}
		out = append(out, pbTask)
	}

	return &taskspb.TaskListResponse{Tasks: out}, nil
}
