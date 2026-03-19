package gateway

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	appgateway "tech-memo/internal/application/gateway"
	"tech-memo/internal/domain"
)

var _ appgateway.MemoGateway = (*SQLiteMemoGateway)(nil)

type SQLiteMemoGateway struct {
	db *sql.DB
}

func NewSQLiteMemoGateway(db *sql.DB) *SQLiteMemoGateway {
	return &SQLiteMemoGateway{db: db}
}

func (g *SQLiteMemoGateway) FindByID(id string) (*domain.Memo, error) {
	row := g.db.QueryRow(
		`SELECT id, user_id, title, content, category_id, parameters, is_pinned, created_at, updated_at, deleted_at
		 FROM memos WHERE id = ? AND deleted_at IS NULL`, id)
	return scanMemo(row)
}

func (g *SQLiteMemoGateway) FindByUserID(userID string) ([]*domain.Memo, error) {
	rows, err := g.db.Query(
		`SELECT id, user_id, title, content, category_id, parameters, is_pinned, created_at, updated_at, deleted_at
		 FROM memos WHERE user_id = ? AND deleted_at IS NULL ORDER BY updated_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMemos(rows)
}

func (g *SQLiteMemoGateway) FindByCategory(userID, categoryID string) ([]*domain.Memo, error) {
	rows, err := g.db.Query(
		`SELECT id, user_id, title, content, category_id, parameters, is_pinned, created_at, updated_at, deleted_at
		 FROM memos WHERE user_id = ? AND category_id = ? AND deleted_at IS NULL ORDER BY updated_at DESC`,
		userID, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMemos(rows)
}

func (g *SQLiteMemoGateway) Search(userID, query string) ([]*domain.Memo, error) {
	like := "%" + query + "%"
	rows, err := g.db.Query(
		`SELECT id, user_id, title, content, category_id, parameters, is_pinned, created_at, updated_at, deleted_at
		 FROM memos WHERE user_id = ? AND (title LIKE ? OR content LIKE ?) AND deleted_at IS NULL ORDER BY updated_at DESC`,
		userID, like, like)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMemos(rows)
}

func (g *SQLiteMemoGateway) Save(memo *domain.Memo) error {
	params, err := json.Marshal(memo.Parameters)
	if err != nil {
		return fmt.Errorf("marshal parameters: %w", err)
	}
	_, err = g.db.Exec(
		`INSERT INTO memos (id, user_id, title, content, category_id, parameters, is_pinned, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		memo.ID, memo.UserID, memo.Title, memo.Content, memo.CategoryID,
		string(params), boolToInt(memo.IsPinned), memo.CreatedAt, memo.UpdatedAt)
	return err
}

func (g *SQLiteMemoGateway) Update(memo *domain.Memo) error {
	params, err := json.Marshal(memo.Parameters)
	if err != nil {
		return fmt.Errorf("marshal parameters: %w", err)
	}
	_, err = g.db.Exec(
		`UPDATE memos SET title=?, content=?, category_id=?, parameters=?, is_pinned=?, updated_at=?
		 WHERE id=? AND deleted_at IS NULL`,
		memo.Title, memo.Content, memo.CategoryID, string(params),
		boolToInt(memo.IsPinned), memo.UpdatedAt, memo.ID)
	return err
}

func (g *SQLiteMemoGateway) Delete(id string) error {
	_, err := g.db.Exec(
		`UPDATE memos SET deleted_at=? WHERE id=? AND deleted_at IS NULL`,
		time.Now(), id)
	return err
}

func scanMemo(row *sql.Row) (*domain.Memo, error) {
	var m domain.Memo
	var paramsJSON string
	var deletedAt sql.NullTime
	var isPinned int

	err := row.Scan(&m.ID, &m.UserID, &m.Title, &m.Content, &m.CategoryID,
		&paramsJSON, &isPinned, &m.CreatedAt, &m.UpdatedAt, &deletedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(paramsJSON), &m.Parameters); err != nil {
		return nil, fmt.Errorf("unmarshal parameters: %w", err)
	}
	m.IsPinned = isPinned == 1
	if deletedAt.Valid {
		m.DeletedAt = &deletedAt.Time
	}
	return &m, nil
}

func scanMemos(rows *sql.Rows) ([]*domain.Memo, error) {
	var memos []*domain.Memo
	for rows.Next() {
		var m domain.Memo
		var paramsJSON string
		var deletedAt sql.NullTime
		var isPinned int

		if err := rows.Scan(&m.ID, &m.UserID, &m.Title, &m.Content, &m.CategoryID,
			&paramsJSON, &isPinned, &m.CreatedAt, &m.UpdatedAt, &deletedAt); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(paramsJSON), &m.Parameters); err != nil {
			return nil, fmt.Errorf("unmarshal parameters: %w", err)
		}
		m.IsPinned = isPinned == 1
		if deletedAt.Valid {
			m.DeletedAt = &deletedAt.Time
		}
		memos = append(memos, &m)
	}
	return memos, rows.Err()
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
