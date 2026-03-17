// internal/application/usecase/memo.go
package usecase

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	appgateway "tech-memo/internal/application/gateway"
	"tech-memo/internal/domain"
)

// ---- インターフェース ----

type MemoUseCase interface {
	GetAll() ([]*domain.Memo, error)
	GetByID(id string) (*domain.Memo, error)
	Search(query string) ([]*domain.Memo, error)
	FindByTag(tag string) ([]*domain.Memo, error)
	Create(title, content string, tags []string, language string) (*domain.Memo, error)
	Update(id, title, content string, tags []string, language string) (*domain.Memo, error)
	Delete(id string) error
}

// ---- インタラクター（実装）----

type memoInteracter struct {
	gw appgateway.MemoGateway
}

func NewMemoInteracter(gw appgateway.MemoGateway) MemoUseCase {
	return &memoInteracter{gw: gw}
}

func (uc *memoInteracter) GetAll() ([]*domain.Memo, error) {
	return uc.gw.FindAll()
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

func (uc *memoInteracter) Search(query string) ([]*domain.Memo, error) {
	return uc.gw.Search(query)
}

func (uc *memoInteracter) FindByTag(tag string) ([]*domain.Memo, error) {
	return uc.gw.FindByTag(tag)
}

func (uc *memoInteracter) Create(title, content string, tags []string, language string) (*domain.Memo, error) {
	memo := &domain.Memo{
		ID:        uuid.New().String(),
		Title:     title,
		Content:   content,
		Tags:      tags,
		Language:  language,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := memo.Validate(); err != nil {
		return nil, err
	}
	if err := uc.gw.Save(memo); err != nil {
		return nil, err
	}
	return memo, nil
}

func (uc *memoInteracter) Update(id, title, content string, tags []string, language string) (*domain.Memo, error) {
	existing, err := uc.gw.FindByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, fmt.Errorf("memo not found: %s", id)
	}

	existing.Title = title
	existing.Content = content
	existing.Tags = tags
	existing.Language = language
	existing.UpdatedAt = time.Now()

	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := uc.gw.Update(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (uc *memoInteracter) Delete(id string) error {
	existing, err := uc.gw.FindByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("memo not found: %s", id)
	}
	return uc.gw.Delete(id)
}
