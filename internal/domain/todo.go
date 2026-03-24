package domain

import "time"

type Todo struct {
	ID          string
	UserID      string
	Title       string
	Content     string
	CategoryID  string
	Parameters  []Parameter
	IsPinned    bool
	DueAt       *time.Time
	CompletedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
