package tasks

import (
	"fmt"
	"strings"

	"github.com/your-org/tasks-service/domain"
)

type tasksService struct {
	repo TasksRepo
}

type TasksService interface {
	CreateTask(task string, isDone bool, userID uint32) (*domain.Task, error)
	GetAllTasks() ([]*domain.Task, error)
	UpdateTask(task string, isDone bool, id uint32) (*domain.Task, error)
	DeleteTask(id uint32) error
}

func NewTasksService(r TasksRepo) TasksService {
	return &tasksService{repo: r}
}

func (s *tasksService) GetAllTasks() ([]*domain.Task, error) {
	return s.repo.GetAllTasks()
}

func (s *tasksService) CreateTask(task string, isDone bool, userID uint32) (*domain.Task, error) {
	if strings.TrimSpace(task) == "" {
		return nil, ErrInvalidInput
	}

	taskToCreate := domain.Task{Task: task, IsDone: isDone, UserID: userID}

	createdTask, err := s.repo.CreateTask(&taskToCreate)
	if err != nil {
		return nil, fmt.Errorf("CreateTask: failed to create the task: %w", err)
	}
	return createdTask, nil
}

func (s *tasksService) UpdateTask(task string, isDone bool, id uint32) (*domain.Task, error) {
	if strings.TrimSpace(task) == "" {
		return nil, ErrInvalidInput
	}

	dm, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	dm.Task = task
	dm.IsDone = isDone

	return s.repo.UpdateTask(dm)
}

func (s *tasksService) DeleteTask(id uint32) error {
	return s.repo.DeleteTask(id)
}
