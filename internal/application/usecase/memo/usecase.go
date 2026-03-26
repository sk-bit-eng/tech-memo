package memo

import (
	"tech-memo/internal/application/dto"
	"tech-memo/internal/domain"
)

type UseCase interface {
	GetByID(id string) (*domain.Memo, error)
	ListByUser(userID string) ([]*domain.Memo, error)
	ListByCategory(userID, categoryID string) ([]*domain.Memo, error)
	Search(userID, query string) ([]*domain.Memo, error)
	Create(input dto.CreateMemoInput) (*domain.Memo, error)
	Update(input dto.UpdateMemoInput) (*domain.Memo, error)
	Delete(id string) error
	TogglePin(id string) (*domain.Memo, error)
}
