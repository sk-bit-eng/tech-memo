// internal/adapter/gateway/model/memo_model.go
package model

import (
	"encoding/json"
	"time"

	"tech-memo/internal/domain"

	"gorm.io/gorm"
)

type MemoModel struct {
	ID         string `gorm:"primaryKey"`
	UserID     string `gorm:"not null;index"`
	Title      string `gorm:"not null"`
	Content    string `gorm:"type:nvarchar(512);not null"`
	CategoryID string `gorm:"default:''"`
	Parameters string `gorm:"type:nvarchar(512);default:'[]'"`
	IsPinned   bool   `gorm:"default:false"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (MemoModel) TableName() string { return "memos" }

func (m MemoModel) ToDomain() *domain.Memo {
	var params []domain.Parameter
	_ = json.Unmarshal([]byte(m.Parameters), &params)
	memo := &domain.Memo{
		ID:         m.ID,
		UserID:     m.UserID,
		Title:      m.Title,
		Content:    m.Content,
		CategoryID: m.CategoryID,
		Parameters: params,
		IsPinned:   m.IsPinned,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
	if m.DeletedAt.Valid {
		memo.DeletedAt = &m.DeletedAt.Time
	}
	return memo
}

func FromMemo(memo *domain.Memo) MemoModel {
	params, _ := json.Marshal(memo.Parameters)
	return MemoModel{
		ID:         memo.ID,
		UserID:     memo.UserID,
		Title:      memo.Title,
		Content:    memo.Content,
		CategoryID: memo.CategoryID,
		Parameters: string(params),
		IsPinned:   memo.IsPinned,
		CreatedAt:  memo.CreatedAt,
		UpdatedAt:  memo.UpdatedAt,
	}
}
