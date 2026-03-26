package todo

import (
	"tech-memo/internal/application/dto"
	"tech-memo/internal/domain"
)

type UseCase interface {
	GetByID(id string) (*domain.Todo, error)
	ListByUser(userID string) ([]*domain.Todo, error)
	ListByCategory(userID, categoryID string) ([]*domain.Todo, error)
	ListPending(userID string) ([]*domain.Todo, error)
	ListCompleted(userID string) ([]*domain.Todo, error)
	Search(userID, query string) ([]*domain.Todo, error)
	Create(input dto.CreateTodoInput) (*domain.Todo, error)
	Update(input dto.UpdateTodoInput) (*domain.Todo, error)
	Complete(id string) error
	Incomplete(id string) error
	Delete(id string) error
	TogglePin(id string) (*domain.Todo, error)
}
