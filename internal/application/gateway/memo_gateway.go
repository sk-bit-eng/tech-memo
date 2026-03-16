// internal/application/gateway/memo_gateway.go
package gateway

import "tech-memo/internal/domain"

type MemoGateway interface {
	FindAll() ([]*domain.Memo, error)
	FindByID(id string) (*domain.Memo, error)
	Search(query string) ([]*domain.Memo, error)
	FindByTag(tag string) ([]*domain.Memo, error)
	Save(memo *domain.Memo) error
	Update(memo *domain.Memo) error
	Delete(id string) error
}
