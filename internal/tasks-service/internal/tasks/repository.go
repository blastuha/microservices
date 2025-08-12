package tasks

import (
	"errors"
	"fmt"

	"github.com/your-org/tasks-service/domain"
	"gorm.io/gorm"
)

type TasksRepo interface {
	CreateTask(t *domain.Task) (*domain.Task, error)
	GetAllTasks() ([]*domain.Task, error)
	UpdateTask(t *domain.Task) (*domain.Task, error)
	DeleteTask(id uint32) error
	GetByID(id uint32) (*domain.Task, error)
	ListTasksByUser(userId uint32) ([]*domain.Task, error)
}

type taskRepo struct {
	db *gorm.DB
}

func NewTaskRepo(db *gorm.DB) TasksRepo {
	return &taskRepo{db: db}
}

// CreateTask создает запись и возвращает domain модель
func (r *taskRepo) CreateTask(dm *domain.Task) (*domain.Task, error) {
	ormTask := (&Task{}).toORM(dm)
	if err := r.db.Create(ormTask).Error; err != nil {
		return nil, fmt.Errorf("CreateTask: failed to create task: %w", err)
	}
	return ormTask.toDomain(), nil
}

// GetAllTasks возвращает слайс domain моделей
func (r *taskRepo) GetAllTasks() ([]*domain.Task, error) {
	var ormTasks []Task
	if err := r.db.Find(&ormTasks).Error; err != nil {
		return nil, fmt.Errorf("GetAllTasks: failed to get tasks: %w", err)
	}

	domainTasks := make([]*domain.Task, len(ormTasks))
	for i := range ormTasks {
		domainTasks[i] = ormTasks[i].toDomain()
	}
	return domainTasks, nil
}

// UpdateTask обновляет orm модель на основе domain и возвращает domain
func (r *taskRepo) UpdateTask(dm *domain.Task) (*domain.Task, error) {
	ormTask := (&Task{}).toORM(dm)
	if err := r.db.Save(ormTask).Error; err != nil {
		return nil, fmt.Errorf("UpdateTask: failed to save task: %w", err)
	}
	return ormTask.toDomain(), nil
}

// DeleteTask удаляет запись по id
func (r *taskRepo) DeleteTask(id uint32) error {
	var ormTask Task
	if err := r.db.First(&ormTask, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTaskNotFound
		}
		return fmt.Errorf("DeleteTask: failed to find task: %w", err)
	}

	if err := r.db.Delete(&ormTask).Error; err != nil {
		return fmt.Errorf("DeleteTask: failed to delete task: %w", err)
	}

	return nil
}

// GetByID возвращает domain модель по строковому id
func (r *taskRepo) GetByID(id uint32) (*domain.Task, error) {
	var ormTask Task
	if err := r.db.First(&ormTask, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, fmt.Errorf("GetByID: failed to find task: %w", err)
	}

	return ormTask.toDomain(), nil
}

func (r *taskRepo) ListTasksByUser(userID uint32) ([]*domain.Task, error) {
	var tasks []Task

	if err := r.db.Where("user_id =?", userID).Find(&tasks).Error; err != nil {
		return nil, fmt.Errorf("ListTasksByUser: failed to get tasks: %w", err)
	}

	out := make([]*domain.Task, 0, len(tasks))

	for _, t := range tasks {
		dm := t.toDomain()
		out = append(out, dm)
	}

	return out, nil
}
