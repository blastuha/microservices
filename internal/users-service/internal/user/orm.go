package user

import (
	"github.com/your-org/users-service/domain"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
}

// Методы конвертации между domain и DB моделями
func (db *User) toDomain() *domain.User {
	return &domain.User{
		ID:       uint32(db.ID),
		Email:    db.Email,
		Password: db.Password,
	}
}

func fromDomain(u *domain.User) *User {
	return &User{
		Model: gorm.Model{
			ID: uint(u.ID),
		},
		Email:    u.Email,
		Password: u.Password,
	}
}
