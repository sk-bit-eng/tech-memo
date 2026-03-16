// internal/adapter/gateway/sqlite_memo_gateway.go
package gateway

import (
	"database/sql"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	appgateway "tech-memo/internal/application/gateway"
	"tech-memo/internal/domain"
)

// compile-time interface check
var _ appgateway.MemoGateway = (*SQLiteMemoGateway)(nil)

type SQLiteMemoGateway struct {
	db *sql.DB
}

func NewSQLiteMemoGateway(dbPath string) (*SQLiteMemoGateway, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	gw := &SQLiteMemoGateway{db: db}
	if err := gw.migrate(); err != nil {
		return nil, err
	}
	return gw, nil
}

func (r *SQLiteMemoGateway) migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS memos (
		id         TEXT PRIMARY KEY,
		title      TEXT NOT NULL,
		content    TEXT NOT NULL,
		tags       TEXT NOT NULL DEFAULT '',
		language   TEXT NOT NULL DEFAULT '',
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);`
	_, err := r.db.Exec(query)
	return err
}

func tagsToString(tags []string) string {
	return strings.Join(tags, ",")
}

func stringToTags(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}

func (r *SQLiteMemoGateway) scanMemo(row *sql.Row) (*domain.Memo, error) {
	var m domain.Memo
	var tagsStr string
	var createdAt, updatedAt string

	err := row.Scan(&m.ID, &m.Title, &m.Content, &tagsStr, &m.Language, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	m.Tags = stringToTags(tagsStr)
	m.CreatedAt, _ = time.Parse(time.RFC3339Nano, createdAt)
	m.UpdatedAt, _ = time.Parse(time.RFC3339Nano, updatedAt)
	return &m, nil
}

func (r *SQLiteMemoGateway) scanMemos(rows *sql.Rows) ([]*domain.Memo, error) {
	var memos []*domain.Memo
	for rows.Next() {
		var m domain.Memo
		var tagsStr string
		var createdAt, updatedAt string

		if err := rows.Scan(&m.ID, &m.Title, &m.Content, &tagsStr, &m.Language, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		m.Tags = stringToTags(tagsStr)
		m.CreatedAt, _ = time.Parse(time.RFC3339Nano, createdAt)
		m.UpdatedAt, _ = time.Parse(time.RFC3339Nano, updatedAt)
		memos = append(memos, &m)
	}
	if memos == nil {
		memos = []*domain.Memo{}
	}
	return memos, rows.Err()
}

func (r *SQLiteMemoGateway) FindAll() ([]*domain.Memo, error) {
	rows, err := r.db.Query(
		`SELECT id, title, content, tags, language, created_at, updated_at FROM memos ORDER BY updated_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanMemos(rows)
}

func (r *SQLiteMemoGateway) FindByID(id string) (*domain.Memo, error) {
	row := r.db.QueryRow(
		`SELECT id, title, content, tags, language, created_at, updated_at FROM memos WHERE id = ?`, id,
	)
	return r.scanMemo(row)
}

func (r *SQLiteMemoGateway) Search(query string) ([]*domain.Memo, error) {
	like := "%" + query + "%"
	rows, err := r.db.Query(
		`SELECT id, title, content, tags, language, created_at, updated_at FROM memos
		 WHERE title LIKE ? OR content LIKE ?
		 ORDER BY updated_at DESC`,
		like, like,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanMemos(rows)
}

func (r *SQLiteMemoGateway) FindByTag(tag string) ([]*domain.Memo, error) {
	rows, err := r.db.Query(
		`SELECT id, title, content, tags, language, created_at, updated_at FROM memos
		 WHERE (',' || tags || ',') LIKE ?
		 ORDER BY updated_at DESC`,
		"%,"+tag+",%",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanMemos(rows)
}

func (r *SQLiteMemoGateway) Save(memo *domain.Memo) error {
	_, err := r.db.Exec(
		`INSERT INTO memos (id, title, content, tags, language, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		memo.ID, memo.Title, memo.Content, tagsToString(memo.Tags), memo.Language,
		memo.CreatedAt.Format(time.RFC3339Nano), memo.UpdatedAt.Format(time.RFC3339Nano),
	)
	return err
}

func (r *SQLiteMemoGateway) Update(memo *domain.Memo) error {
	_, err := r.db.Exec(
		`UPDATE memos SET title=?, content=?, tags=?, language=?, updated_at=? WHERE id=?`,
		memo.Title, memo.Content, tagsToString(memo.Tags), memo.Language,
		memo.UpdatedAt.Format(time.RFC3339Nano), memo.ID,
	)
	return err
}

func (r *SQLiteMemoGateway) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM memos WHERE id=?`, id)
	return err
}
