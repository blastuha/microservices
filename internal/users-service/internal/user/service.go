package user

import (
	"errors"
	"fmt"
)

type UsersService interface {
	GetAllUsers() ([]User, error)
	CreateUser(email, password string) (*User, error)
	UpdateUser(id uint32, email, password *string) (*User, error)
	DeleteUser(id uint32) error
	GetUserByID(id uint32) (*User, error)
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

func (u *usersService) GetAllUsers() ([]User, error) {
	userList, err := u.repo.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("usersService.GetAllUsers: %w", err)
	}

	return userList, nil
}

func (u *usersService) CreateUser(email, password string) (*User, error) {
	// Валидация уже выполнена в gRPC handler
	userToCreate := User{
		Email:    email,
		Password: password,
	}

	createdUser, err := u.repo.CreateUser(&userToCreate)
	if err != nil {
		return nil, fmt.Errorf("usersService.CreateUser: %w", err)
	}

	return createdUser, nil
}

func (u *usersService) UpdateUser(id uint32, email, password *string) (*User, error) {
	existingUser, err := u.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	// Проверяем на nil email и password
	// НО не проверяем на пустоту - это уже сделал handler
	if email != nil {
		existingUser.Email = *email
	}
	if password != nil {
		existingUser.Password = *password
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

func (u *usersService) GetUserByID(id uint32) (*User, error) {
	user, err := u.repo.GetUserByID(id)
	if err != nil {
		if errors.Is(err, ErrUserNoFound) {
			return nil, ErrUserNoFound
		}
		return nil, fmt.Errorf("usersService.GetUserByID: %w", err)
	}

	return user, nil
}
