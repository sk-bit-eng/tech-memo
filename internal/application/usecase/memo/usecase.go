package memo

import (
	memodto "tech-memo/internal/application/dto/memo"
	"tech-memo/internal/domain"
)

type UseCase interface {
	GetByID(id string) (*domain.Memo, error)
	ListByUser(userID string) ([]*domain.Memo, error)
	ListByCategory(userID, categoryID string) ([]*domain.Memo, error)
	Search(userID, query string) ([]*domain.Memo, error)
	Create(input memodto.CreateInput) (*domain.Memo, error)
	Update(input memodto.UpdateInput) (*domain.Memo, error)
	Delete(id string) error
	TogglePin(id string) (*domain.Memo, error)
}
