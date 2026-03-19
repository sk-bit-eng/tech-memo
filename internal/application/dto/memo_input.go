// internal/application/dto/memo_input.go
package dto

import "tech-memo/internal/domain"

type CreateMemoInput struct {
	UserID     string
	Title      string
	Content    string
	CategoryID string
	Parameters []domain.Parameter
	IsPinned   bool
}

type UpdateMemoInput struct {
	ID         string
	Title      string
	Content    string
	CategoryID string
	Parameters []domain.Parameter
	IsPinned   bool
}
