package memo

import (
	"fmt"
	"time"

	"tech-memo/internal/application/dto"
	memogtw "tech-memo/internal/application/gateway/memo"
	"tech-memo/internal/domain"

	"github.com/google/uuid"
)

type interactor struct {
	gw memogtw.Repository
}

func NewInteractor(gw memogtw.Repository) UseCase {
	return &interactor{gw: gw}
}

func (uc *interactor) GetByID(id string) (*domain.Memo, error) {
	memo, err := uc.gw.FindByID(id)
	if err != nil {
		return nil, err
	}
	if memo == nil {
		return nil, fmt.Errorf("memo not found: %s", id)
	}
	return memo, nil
}

func (uc *interactor) ListByUser(userID string) ([]*domain.Memo, error) {
	return uc.gw.FindByUserID(userID)
}

func (uc *interactor) ListByCategory(userID, categoryID string) ([]*domain.Memo, error) {
	return uc.gw.FindByCategory(userID, categoryID)
}

func (uc *interactor) Search(userID, query string) ([]*domain.Memo, error) {
	return uc.gw.Search(userID, query)
}

func (uc *interactor) Create(input dto.CreateMemoInput) (*domain.Memo, error) {
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

func (uc *interactor) Update(input dto.UpdateMemoInput) (*domain.Memo, error) {
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

func (uc *interactor) Delete(id string) error {
	memo, err := uc.gw.FindByID(id)
	if err != nil {
		return err
	}
	if memo == nil {
		return fmt.Errorf("memo not found: %s", id)
	}
	return uc.gw.Delete(id)
}

func (uc *interactor) TogglePin(id string) (*domain.Memo, error) {
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
