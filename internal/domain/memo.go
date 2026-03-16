// internal/domain/memo.go
package domain

import (
	"errors"
	"strings"
	"time"
)

type Memo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	Language  string    `json:"language"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *Memo) Validate() error {
	if strings.TrimSpace(m.Title) == "" {
		return errors.New("title is required")
	}
	if len(m.Title) > 200 {
		return errors.New("title must be 200 characters or less")
	}
	if strings.TrimSpace(m.Content) == "" {
		return errors.New("content is required")
	}
	return nil
}
