package todo

import (
	tododto "tech-memo/internal/application/dto/todo"
	"tech-memo/internal/domain"
)

type UseCase interface {
	GetByID(id string) (*domain.Todo, error)
	ListByUser(userID string) ([]*domain.Todo, error)
	ListByCategory(userID, categoryID string) ([]*domain.Todo, error)
	ListPending(userID string) ([]*domain.Todo, error)
	ListCompleted(userID string) ([]*domain.Todo, error)
	Search(userID, query string) ([]*domain.Todo, error)
	Create(input tododto.CreateInput) (*domain.Todo, error)
	Update(input tododto.UpdateInput) (*domain.Todo, error)
	Complete(id string) error
	Incomplete(id string) error
	Delete(id string) error
	TogglePin(id string) (*domain.Todo, error)
}
