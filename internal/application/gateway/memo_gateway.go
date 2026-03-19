package gateway

import "tech-memo/internal/domain"

type MemoGateway interface {
	FindByID(id string) (*domain.Memo, error)
	FindByUserID(userID string) ([]*domain.Memo, error)
	FindByCategory(userID, categoryID string) ([]*domain.Memo, error)
	Search(userID, query string) ([]*domain.Memo, error)
	Save(memo *domain.Memo) error
	Update(memo *domain.Memo) error
	Delete(id string) error
}
