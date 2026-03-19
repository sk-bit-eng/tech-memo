// internal/application/dto/todo_input.go
package dto

import (
	"time"

	"tech-memo/internal/domain"
)

type CreateTodoInput struct {
	UserID     string
	Title      string
	Content    string
	CategoryID string
	Parameters []domain.Parameter
	IsPinned   bool
	DueAt      *time.Time
}

type UpdateTodoInput struct {
	ID         string
	Title      string
	Content    string
	CategoryID string
	Parameters []domain.Parameter
	IsPinned   bool
	DueAt      *time.Time
}
