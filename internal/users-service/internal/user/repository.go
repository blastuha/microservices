package user

import (
	"errors"
	"fmt"

	"github.com/your-org/users-service/domain"
	"gorm.io/gorm"
)

type UsersRepo interface {
	GetAllUsers() ([]*domain.User, error)
	// GetTasksForUser(id uint) ([]tasksService.Task, error)
	CreateUser(u *domain.User) (*domain.User, error)
	UpdateUser(u *domain.User) (*domain.User, error)
	DeleteUser(id uint32) error
	GetUserByID(id uint32) (*domain.User, error)
}

type usersRepo struct {
	db *gorm.DB
}

func NewUsersRepo(db *gorm.DB) UsersRepo {
	return &usersRepo{db: db}
}

func (repo *usersRepo) GetUserByID(id uint32) (*domain.User, error) {
	var u User
	if err := repo.db.First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNoFound
		}
		return nil, err
	}

	dm := u.toDomain()

	return dm, nil
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

func (repo *usersRepo) GetAllUsers() ([]*domain.User, error) {

	var users []User
	if err := repo.db.Preload("Tasks").Find(&users).Error; err != nil {
		return nil, fmt.Errorf("usersRepo.GetAllUsers: %w", err)
	}

	var dmUsers []*domain.User
	for _, user := range users {
		dmUsers = append(dmUsers, user.toDomain())
	}

	return dmUsers, nil
}

func (repo *usersRepo) CreateUser(u *domain.User) (*domain.User, error) {
	orm := fromDomain(u)

	if err := repo.db.Create(orm).Error; err != nil {
		return nil, fmt.Errorf("usersRepo.CreateUser: %w", err)
	}

	dm := orm.toDomain()

	return dm, nil
}

func (repo *usersRepo) UpdateUser(u *domain.User) (*domain.User, error) {
	orm := fromDomain(u)
	if err := repo.db.Save(orm).Error; err != nil {
		return nil, err
	}

	dm := orm.toDomain()

	return dm, nil
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
