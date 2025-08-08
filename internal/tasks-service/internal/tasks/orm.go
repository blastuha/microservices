package tasks

import (
	"github.com/your-org/tasks-service/domain"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Task   string `gorm:"type:varchar(255);not null" json:"task"`
	IsDone bool   `gorm:"default:false" json:"is_done"`
}

func (t *Task) toDomain() *domain.Task {
	return &domain.Task{
		ID:     uint32(t.ID),
		Task:   t.Task,
		IsDone: t.IsDone,
	}
}

func (t *Task) toORM(dm *domain.Task) *Task {
	return &Task{Model: gorm.Model{ID: uint(dm.ID)}, Task: dm.Task, IsDone: dm.IsDone}
}
