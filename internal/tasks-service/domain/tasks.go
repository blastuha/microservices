package domain

type Task struct {
	ID     uint32
	Task   string
	IsDone bool
	UserID uint32
}
