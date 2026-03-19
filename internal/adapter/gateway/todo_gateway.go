package gateway

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	appgateway "tech-memo/internal/application/gateway"
	"tech-memo/internal/domain"
)

var _ appgateway.TodoGateway = (*SQLiteTodoGateway)(nil)

type SQLiteTodoGateway struct {
	db *sql.DB
}

func NewSQLiteTodoGateway(db *sql.DB) *SQLiteTodoGateway {
	return &SQLiteTodoGateway{db: db}
}

func (g *SQLiteTodoGateway) FindByID(id string) (*domain.Todo, error) {
	row := g.db.QueryRow(
		`SELECT id, user_id, title, content, category_id, parameters, is_pinned,
		        due_at, completed_at, created_at, updated_at, deleted_at
		 FROM todos WHERE id = ? AND deleted_at IS NULL`, id)
	return scanTodo(row)
}

func (g *SQLiteTodoGateway) FindByUserID(userID string) ([]*domain.Todo, error) {
	rows, err := g.db.Query(
		`SELECT id, user_id, title, content, category_id, parameters, is_pinned,
		        due_at, completed_at, created_at, updated_at, deleted_at
		 FROM todos WHERE user_id = ? AND deleted_at IS NULL ORDER BY updated_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTodos(rows)
}

func (g *SQLiteTodoGateway) FindByCategory(userID, categoryID string) ([]*domain.Todo, error) {
	rows, err := g.db.Query(
		`SELECT id, user_id, title, content, category_id, parameters, is_pinned,
		        due_at, completed_at, created_at, updated_at, deleted_at
		 FROM todos WHERE user_id = ? AND category_id = ? AND deleted_at IS NULL ORDER BY updated_at DESC`,
		userID, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTodos(rows)
}

func (g *SQLiteTodoGateway) FindPending(userID string) ([]*domain.Todo, error) {
	rows, err := g.db.Query(
		`SELECT id, user_id, title, content, category_id, parameters, is_pinned,
		        due_at, completed_at, created_at, updated_at, deleted_at
		 FROM todos WHERE user_id = ? AND completed_at IS NULL AND deleted_at IS NULL ORDER BY updated_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTodos(rows)
}

func (g *SQLiteTodoGateway) FindCompleted(userID string) ([]*domain.Todo, error) {
	rows, err := g.db.Query(
		`SELECT id, user_id, title, content, category_id, parameters, is_pinned,
		        due_at, completed_at, created_at, updated_at, deleted_at
		 FROM todos WHERE user_id = ? AND completed_at IS NOT NULL AND deleted_at IS NULL ORDER BY updated_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTodos(rows)
}

func (g *SQLiteTodoGateway) Search(userID, query string) ([]*domain.Todo, error) {
	like := "%" + query + "%"
	rows, err := g.db.Query(
		`SELECT id, user_id, title, content, category_id, parameters, is_pinned,
		        due_at, completed_at, created_at, updated_at, deleted_at
		 FROM todos WHERE user_id = ? AND (title LIKE ? OR content LIKE ?) AND deleted_at IS NULL ORDER BY updated_at DESC`,
		userID, like, like)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTodos(rows)
}

func (g *SQLiteTodoGateway) Save(todo *domain.Todo) error {
	params, err := json.Marshal(todo.Parameters)
	if err != nil {
		return fmt.Errorf("marshal parameters: %w", err)
	}
	_, err = g.db.Exec(
		`INSERT INTO todos (id, user_id, title, content, category_id, parameters, is_pinned,
		                    due_at, completed_at, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		todo.ID, todo.UserID, todo.Title, todo.Content, todo.CategoryID,
		string(params), boolToInt(todo.IsPinned),
		todo.DueAt, todo.CompletedAt, todo.CreatedAt, todo.UpdatedAt)
	return err
}

func (g *SQLiteTodoGateway) Update(todo *domain.Todo) error {
	params, err := json.Marshal(todo.Parameters)
	if err != nil {
		return fmt.Errorf("marshal parameters: %w", err)
	}
	_, err = g.db.Exec(
		`UPDATE todos SET title=?, content=?, category_id=?, parameters=?, is_pinned=?,
		                  due_at=?, completed_at=?, updated_at=?
		 WHERE id=? AND deleted_at IS NULL`,
		todo.Title, todo.Content, todo.CategoryID, string(params),
		boolToInt(todo.IsPinned), todo.DueAt, todo.CompletedAt, todo.UpdatedAt, todo.ID)
	return err
}

func (g *SQLiteTodoGateway) Delete(id string) error {
	_, err := g.db.Exec(
		`UPDATE todos SET deleted_at=? WHERE id=? AND deleted_at IS NULL`,
		time.Now(), id)
	return err
}

func scanTodo(row *sql.Row) (*domain.Todo, error) {
	var td domain.Todo
	var paramsJSON string
	var isPinned int
	var dueAt, completedAt, deletedAt sql.NullTime

	err := row.Scan(&td.ID, &td.UserID, &td.Title, &td.Content, &td.CategoryID,
		&paramsJSON, &isPinned, &dueAt, &completedAt, &td.CreatedAt, &td.UpdatedAt, &deletedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(paramsJSON), &td.Parameters); err != nil {
		return nil, fmt.Errorf("unmarshal parameters: %w", err)
	}
	td.IsPinned = isPinned == 1
	if dueAt.Valid {
		td.DueAt = &dueAt.Time
	}
	if completedAt.Valid {
		td.CompletedAt = &completedAt.Time
	}
	if deletedAt.Valid {
		td.DeletedAt = &deletedAt.Time
	}
	return &td, nil
}

func scanTodos(rows *sql.Rows) ([]*domain.Todo, error) {
	var todos []*domain.Todo
	for rows.Next() {
		var td domain.Todo
		var paramsJSON string
		var isPinned int
		var dueAt, completedAt, deletedAt sql.NullTime

		if err := rows.Scan(&td.ID, &td.UserID, &td.Title, &td.Content, &td.CategoryID,
			&paramsJSON, &isPinned, &dueAt, &completedAt, &td.CreatedAt, &td.UpdatedAt, &deletedAt); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(paramsJSON), &td.Parameters); err != nil {
			return nil, fmt.Errorf("unmarshal parameters: %w", err)
		}
		td.IsPinned = isPinned == 1
		if dueAt.Valid {
			td.DueAt = &dueAt.Time
		}
		if completedAt.Valid {
			td.CompletedAt = &completedAt.Time
		}
		if deletedAt.Valid {
			td.DeletedAt = &deletedAt.Time
		}
		todos = append(todos, &td)
	}
	return todos, rows.Err()
}
