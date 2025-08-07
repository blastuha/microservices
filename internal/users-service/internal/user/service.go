package user

import (
	"errors"
	"fmt"

	"github.com/your-org/users-service/domain"
	// "github.com/your-org/users-service/internal/web/users"
	// "task1/internal/tasksService"
	// "task1/internal/web/users"
)

type UsersService interface {
	GetAllUsers() ([]*domain.User, error)
	CreateUser(email string, password string) (*domain.User, error)
	UpdateUser(id uint32, email string, password string) (*domain.User, error)
	DeleteUser(id uint32) error
	GetUserByID(id uint32) (*domain.User, error)
	// GetTasksForUser(id uint) ([]tasksService.Task, error)
}

type usersService struct {
	repo UsersRepo
}

// func (u *usersService) GetTasksForUser(id uint) ([]tasksService.Task, error) {
// 	tasks, err := u.repo.GetTasksForUser(id)
// 	if errors.Is(err, ErrUserNoFound) {
// 		return nil, ErrUserNoFound
// 	}
// 	if err != nil {
// 		return nil, fmt.Errorf("usersService.GetTasksForUser: %w", err)
// 	}
// 	return tasks, nil
// }

func NewUsersService(repo UsersRepo) UsersService {
	return &usersService{repo: repo}
}

func (u *usersService) GetAllUsers() ([]*domain.User, error) {
	userList, err := u.repo.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("usersService.GetAllUsers: %w", err)
	}

	return userList, nil
}

func (u *usersService) CreateUser(email string, password string) (*domain.User, error) {
	// Валидация уже выполнена в gRPC handler
	userToCreate := domain.User{
		Email:    email,
		Password: password,
	}

	createdUser, err := u.repo.CreateUser(&userToCreate)
	if err != nil {
		return nil, fmt.Errorf("usersService.CreateUser: %w", err)
	}

	return createdUser, nil
}

func (u *usersService) UpdateUser(id uint32, email string, password string) (*domain.User, error) {
	existingUser, err := u.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	// Валидация уже выполнена в gRPC handler
	if email != "" {
		existingUser.Email = email
	}
	if password != "" {
		existingUser.Password = password
	}

	updatedUser, err := u.repo.UpdateUser(existingUser)
	if err != nil {
		return nil, fmt.Errorf("usersService.UpdateUser: %w", err)
	}

	return updatedUser, nil
}

func (u *usersService) DeleteUser(id uint32) error {
	err := u.repo.DeleteUser(id)
	if err != nil {
		if errors.Is(err, ErrUserNoFound) {
			return ErrUserNoFound
		}
		return fmt.Errorf("usersService.DeleteUser: %w", err)
	}

	return nil
}

func (u *usersService) GetUserByID(id uint32) (*domain.User, error) {
	user, err := u.repo.GetUserByID(id)
	if err != nil {
		if errors.Is(err, ErrUserNoFound) {
			return nil, ErrUserNoFound
		}
		return nil, fmt.Errorf("usersService.GetUserByID: %w", err)
	}

	return user, nil
}
