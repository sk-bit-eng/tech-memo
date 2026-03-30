// internal/adapter/gateway/model/todo_model.go
package model

import (
	"encoding/json"
	"time"

	"tech-memo/internal/domain"

	"gorm.io/gorm"
)

type TodoModel struct {
	ID          string `gorm:"primaryKey"`
	UserID      string `gorm:"not null;index"`
	Title       string `gorm:"not null"`
	Content     string `gorm:"type:nvarchar(512);not null"`
	CategoryID  string `gorm:"default:''"`
	Parameters  string `gorm:"type:nvarchar(512);default:'[]'"`
	IsPinned    bool   `gorm:"default:false"`
	DueAt       *time.Time
	CompletedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (TodoModel) TableName() string { return "todos" }

func (m TodoModel) ToDomain() *domain.Todo {
	var params []domain.Parameter
	_ = json.Unmarshal([]byte(m.Parameters), &params)
	todo := &domain.Todo{
		ID:          m.ID,
		UserID:      m.UserID,
		Title:       m.Title,
		Content:     m.Content,
		CategoryID:  m.CategoryID,
		Parameters:  params,
		IsPinned:    m.IsPinned,
		DueAt:       m.DueAt,
		CompletedAt: m.CompletedAt,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
	if m.DeletedAt.Valid {
		todo.DeletedAt = &m.DeletedAt.Time
	}
	return todo
}

func FromTodo(todo *domain.Todo) TodoModel {
	params, _ := json.Marshal(todo.Parameters)
	return TodoModel{
		ID:          todo.ID,
		UserID:      todo.UserID,
		Title:       todo.Title,
		Content:     todo.Content,
		CategoryID:  todo.CategoryID,
		Parameters:  string(params),
		IsPinned:    todo.IsPinned,
		DueAt:       todo.DueAt,
		CompletedAt: todo.CompletedAt,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}
}
