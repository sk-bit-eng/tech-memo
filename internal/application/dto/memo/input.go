package memo

import "tech-memo/internal/domain"

type CreateInput struct {
	UserID     string
	Title      string
	Content    string
	CategoryID string
	Parameters []domain.Parameter
	IsPinned   bool
}

type UpdateInput struct {
	ID         string
	Title      string
	Content    string
	CategoryID string
	Parameters []domain.Parameter
	IsPinned   bool
}
