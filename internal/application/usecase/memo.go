// internal/application/usecase/memo.go
package usecase

import (
	"fmt"
	"time"

	"tech-memo/internal/application/dto"
	appgateway "tech-memo/internal/application/gateway"
	"tech-memo/internal/domain"

	"github.com/google/uuid"
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

// ---- インタラクター（実装） ----
type memoInteractor struct {
	gw appgateway.MemoGateway
}

// コンストラクタ
func NewMemoInteractor(gw appgateway.MemoGateway) MemoUseCase {
	return &memoInteractor{gw: gw}
}

// メモの取得、作成、更新、削除などのビジネスロジックを実装
func (uc *memoInteractor) GetByID(id string) (*domain.Memo, error) {
	memo, err := uc.gw.FindByID(id)
	if err != nil {
		return nil, err
	}
	if memo == nil {
		return nil, fmt.Errorf("memo not found: %s", id)
	}
	return memo, nil
}

func (uc *memoInteractor) ListByUser(userID string) ([]*domain.Memo, error) {
	return uc.gw.FindByUserID(userID)
}

func (uc *memoInteractor) ListByCategory(userID, categoryID string) ([]*domain.Memo, error) {
	return uc.gw.FindByCategory(userID, categoryID)
}

func (uc *memoInteractor) Search(userID, query string) ([]*domain.Memo, error) {
	return uc.gw.Search(userID, query)
}

func (uc *memoInteractor) Create(input dto.CreateMemoInput) (*domain.Memo, error) {
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

func (uc *memoInteractor) Update(input dto.UpdateMemoInput) (*domain.Memo, error) {
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

func (uc *memoInteractor) Delete(id string) error {
	memo, err := uc.gw.FindByID(id)
	if err != nil {
		return err
	}
	if memo == nil {
		return fmt.Errorf("memo not found: %s", id)
	}
	return uc.gw.Delete(id)
}

func (uc *memoInteractor) TogglePin(id string) (*domain.Memo, error) {
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
