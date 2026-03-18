package domain

import (
	"time"
)

type Category struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// 後々、管理画面を作成し、カテゴリの管理を可能にするが現状はマスタDBを直操作する
