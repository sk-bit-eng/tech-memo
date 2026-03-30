package model

import (
	"time"

	"tech-memo/internal/domain"

	"gorm.io/gorm"
)

type UserModel struct {
	ID        string `gorm:"type:nvarchar(64);primaryKey"`
	Name      string `gorm:"type:nvarchar(100);not null"`
	Email     string `gorm:"type:nvarchar(255);not null;uniqueIndex"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (UserModel) TableName() string { return "users" }

func (u UserModel) ToDomain() *domain.User {
	user := &domain.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
	if u.DeletedAt.Valid {
		user.DeletedAt = &u.DeletedAt.Time
	}
	return user
}

func FromUser(user *domain.User) UserModel {
	return UserModel{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
