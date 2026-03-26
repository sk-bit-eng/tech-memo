package todo

import "tech-memo/internal/domain"

type Repository interface {
	FindByID(id string) (*domain.Todo, error)
	FindByUserID(userID string) ([]*domain.Todo, error)
	FindByCategory(userID, categoryID string) ([]*domain.Todo, error)
	FindPending(userID string) ([]*domain.Todo, error)
	FindCompleted(userID string) ([]*domain.Todo, error)
	Search(userID, query string) ([]*domain.Todo, error)
	Save(todo *domain.Todo) error
	Update(todo *domain.Todo) error
	Delete(id string) error
}
