package domain

import "time"

type Memo struct {
	ID         string
	UserID     string
	Title      string
	Content    string
	CategoryID string
	Parameters []Parameter
	IsPinned   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}
