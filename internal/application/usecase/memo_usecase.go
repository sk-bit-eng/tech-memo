// internal/application/usecase/memo_usecase.go
package usecase

import "tech-memo/internal/domain"

type MemoUseCase interface {
	GetAll() ([]*domain.Memo, error)
	GetByID(id string) (*domain.Memo, error)
	Search(query string) ([]*domain.Memo, error)
	FindByTag(tag string) ([]*domain.Memo, error)
	Create(title, content string, tags []string, language string) (*domain.Memo, error)
	Update(id, title, content string, tags []string, language string) (*domain.Memo, error)
	Delete(id string) error
}
