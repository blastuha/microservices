package user

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type UsersRepo interface {
	GetAllUsers() ([]User, error)
	// GetTasksForUser(id uint) ([]tasksService.Task, error)
	CreateUser(u *User) (*User, error)
	UpdateUser(u *User) (*User, error)
	DeleteUser(id uint32) error
	GetUserByID(id uint32) (*User, error)
}

type usersRepo struct {
	db *gorm.DB
}

func (repo *usersRepo) GetUserByID(id uint32) (*User, error) {
	var u User
	if err := repo.db.First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNoFound
		}
		return nil, err
	}
	return &u, nil
}

// func (repo *usersRepo) GetTasksForUser(id uint) ([]tasksService.Task, error) {
// 	var tasks []tasksService.Task
// 	if err := repo.db.
// 		Where("user_id = ?", uint(id)).
// 		Find(&tasks).
// 		Error; err != nil {
// 		return nil, fmt.Errorf("usersRepo.GetTasksForUser: %w", err)
// 	}

// 	return tasks, nil
// }

func NewUsersRepo(db *gorm.DB) UsersRepo {
	return &usersRepo{db: db}
}

func (repo *usersRepo) GetAllUsers() ([]User, error) {
	var users []User
	if err := repo.db.Preload("Tasks").Find(&users).Error; err != nil {
		return nil, fmt.Errorf("usersRepo.GetAllUsers: %w", err)
	}

	return users, nil
}

func (repo *usersRepo) CreateUser(u *User) (*User, error) {
	if err := repo.db.Create(u).Error; err != nil {
		return nil, fmt.Errorf("usersRepo.CreateUser: %w", err)
	}

	return u, nil
}

func (repo *usersRepo) UpdateUser(u *User) (*User, error) {
	if err := repo.db.Save(u).Error; err != nil {
		return nil, err
	}

	return u, nil
}

func (repo *usersRepo) DeleteUser(id uint32) error {
	var user User

	if err := repo.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNoFound
		}

		return fmt.Errorf("userRepo.DeleteUser: %w", err)
	}

	if err := repo.db.Delete(&user).Error; err != nil {
		return fmt.Errorf("userRepo.DeleteUser: %w", err)
	}

	return nil
}
