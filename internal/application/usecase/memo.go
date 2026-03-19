// internal/application/usecase/memo.go
package usecase

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	appgateway "tech-memo/internal/application/gateway"
	"tech-memo/internal/application/dto"
	"tech-memo/internal/domain"
)

// ---- インターフェース ----

type MemoUseCase interface {
	GetByID(id string) (*domain.Memo, error)
	ListByUser(userID string) ([]*domain.Memo, error)
	ListByCategory(userID, categoryID string) ([]*domain.Memo, error)
	Search(userID, query string) ([]*domain.Memo, error)
	Create(input dto.CreateMemoInput) (*domain.Memo, error)
	Update(input dto.UpdateMemoInput) (*domain.Memo, error)
	Delete(id string) error
	TogglePin(id string) (*domain.Memo, error)
}

// ---- インタラクター（実装）----

type memoInteracter struct {
	gw appgateway.MemoGateway
}

func NewMemoInteracter(gw appgateway.MemoGateway) MemoUseCase {
	return &memoInteracter{gw: gw}
}

func (uc *memoInteracter) GetByID(id string) (*domain.Memo, error) {
	memo, err := uc.gw.FindByID(id)
	if err != nil {
		return nil, err
	}
	if memo == nil {
		return nil, fmt.Errorf("memo not found: %s", id)
	}
	return memo, nil
}

func (uc *memoInteracter) ListByUser(userID string) ([]*domain.Memo, error) {
	return uc.gw.FindByUserID(userID)
}

func (uc *memoInteracter) ListByCategory(userID, categoryID string) ([]*domain.Memo, error) {
	return uc.gw.FindByCategory(userID, categoryID)
}

func (uc *memoInteracter) Search(userID, query string) ([]*domain.Memo, error) {
	return uc.gw.Search(userID, query)
}

func (uc *memoInteracter) Create(input dto.CreateMemoInput) (*domain.Memo, error) {
	now := time.Now()
	memo := &domain.Memo{
		ID:         uuid.New().String(),
		UserID:     input.UserID,
		Title:      input.Title,
		Content:    input.Content,
		CategoryID: input.CategoryID,
		Parameters: input.Parameters,
		IsPinned:   input.IsPinned,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := uc.gw.Save(memo); err != nil {
		return nil, err
	}
	return memo, nil
}

func (uc *memoInteracter) Update(input dto.UpdateMemoInput) (*domain.Memo, error) {
	memo, err := uc.gw.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if memo == nil {
		return nil, fmt.Errorf("memo not found: %s", input.ID)
	}
	memo.Title = input.Title
	memo.Content = input.Content
	memo.CategoryID = input.CategoryID
	memo.Parameters = input.Parameters
	memo.IsPinned = input.IsPinned
	memo.UpdatedAt = time.Now()
	if err := uc.gw.Update(memo); err != nil {
		return nil, err
	}
	return memo, nil
}

func (uc *memoInteracter) Delete(id string) error {
	memo, err := uc.gw.FindByID(id)
	if err != nil {
		return err
	}
	if memo == nil {
		return fmt.Errorf("memo not found: %s", id)
	}
	return uc.gw.Delete(id)
}

func (uc *memoInteracter) TogglePin(id string) (*domain.Memo, error) {
	memo, err := uc.gw.FindByID(id)
	if err != nil {
		return nil, err
	}
	if memo == nil {
		return nil, fmt.Errorf("memo not found: %s", id)
	}
	memo.IsPinned = !memo.IsPinned
	memo.UpdatedAt = time.Now()
	if err := uc.gw.Update(memo); err != nil {
		return nil, err
	}
	return memo, nil
}
