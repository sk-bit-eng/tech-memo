// internal/application/interacter/memo_interacter.go
package interacter

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"tech-memo/internal/application/gateway"
	"tech-memo/internal/domain"
)

type MemoInteracter struct {
	gw gateway.MemoGateway
}

func NewMemoInteracter(gw gateway.MemoGateway) *MemoInteracter {
	return &MemoInteracter{gw: gw}
}

func (uc *MemoInteracter) GetAll() ([]*domain.Memo, error) {
	return uc.gw.FindAll()
}

func (uc *MemoInteracter) GetByID(id string) (*domain.Memo, error) {
	memo, err := uc.gw.FindByID(id)
	if err != nil {
		return nil, err
	}
	if memo == nil {
		return nil, fmt.Errorf("memo not found: %s", id)
	}
	return memo, nil
}

func (uc *MemoInteracter) Search(query string) ([]*domain.Memo, error) {
	return uc.gw.Search(query)
}

func (uc *MemoInteracter) FindByTag(tag string) ([]*domain.Memo, error) {
	return uc.gw.FindByTag(tag)
}

func (uc *MemoInteracter) Create(title, content string, tags []string, language string) (*domain.Memo, error) {
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

func (uc *MemoInteracter) Update(id, title, content string, tags []string, language string) (*domain.Memo, error) {
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

func (uc *MemoInteracter) Delete(id string) error {
	existing, err := uc.gw.FindByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("memo not found: %s", id)
	}
	return uc.gw.Delete(id)
}
