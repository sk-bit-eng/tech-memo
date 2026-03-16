# Restructure Clean Architecture Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Reorganize the existing tech-memo Go project into a layered directory structure under `internal/`, with clean separation of Domain, Application, Adapter, and Infrastructure.

**Architecture:** The new structure places all code under `internal/` with explicit layers: `domain` (entities), `application` (use case interfaces + interacters + gateway interfaces), `adapter` (controllers, presenters, gateway implementations), `infrastructure` (HTTP server, DI, middleware). The entry point moves to `cmd/main.go`.

**Tech Stack:** Go 1.22, SQLite (`github.com/mattn/go-sqlite3`), UUID (`github.com/google/uuid`), standard `net/http`

---

## Layer Responsibilities

| Layer | Package path | Role |
|---|---|---|
| Domain | `internal/domain` | Entity + validation. No external dependencies. |
| Application/usecase | `internal/application/usecase` | Use case **interfaces** (I = Interface marker) |
| Application/interacter | `internal/application/interacter` | Use case **implementations** (calls gateway interface) |
| Application/gateway | `internal/application/gateway` | Repository/gateway **interfaces** |
| Adapter/controller | `internal/adapter/controller` | Parse HTTP input → call use case |
| Adapter/presenter | `internal/adapter/presenter` | Write JSON HTTP response |
| Adapter/gateway | `internal/adapter/gateway` | SQLite implementation of gateway interface |
| Infrastructure/api | `internal/infrastructure/api` | HTTP wiring: handler, router, DI |
| Infrastructure/middleware | `internal/infrastructure/middleware` | HTTP middleware (auth stub) |
| Helper/util | `internal/helper/util` | Shared utilities |
| cmd | `cmd/` | Entry point only |

## Dependency Rule

```
cmd → infrastructure/api → adapter/* → application/* → domain
                                      ↑
                         adapter/gateway implements application/gateway (interface)
```

---

## Task 1: Create `internal/domain/memo.go`

**Files:**
- Create: `internal/domain/memo.go`

**Step 1: Create the file** (copy from `domain/memo.go`, change package path only — content is unchanged)

```go
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
```

**Step 2: Verify no syntax errors**

```bash
cd C:/Users/PC1043/Documents/tech-memo
go vet ./internal/domain/...
```

Expected: no output (success)

---

## Task 2: Create `internal/application/gateway/memo_gateway.go`

**Files:**
- Create: `internal/application/gateway/memo_gateway.go`

This replaces `domain/memo_repository.go`. The interface is renamed `MemoGateway` (was `MemoRepository`) to align with the gateway pattern.

**Step 1: Create the file**

```go
// internal/application/gateway/memo_gateway.go
package gateway

import "tech-memo/internal/domain"

type MemoGateway interface {
	FindAll() ([]*domain.Memo, error)
	FindByID(id string) (*domain.Memo, error)
	Search(query string) ([]*domain.Memo, error)
	FindByTag(tag string) ([]*domain.Memo, error)
	Save(memo *domain.Memo) error
	Update(memo *domain.Memo) error
	Delete(id string) error
}
```

**Step 2: Verify**

```bash
go vet ./internal/application/gateway/...
```

Expected: no output

---

## Task 3: Create `internal/application/usecase/memo_usecase.go`

**Files:**
- Create: `internal/application/usecase/memo_usecase.go`

This is a **new** file — the use case interface that the interacter will implement.

**Step 1: Create the file**

```go
// internal/application/usecase/memo_usecase.go
package usecase

import "tech-memo/internal/domain"

type MemoUseCase interface {
	GetAll() ([]*domain.Memo, error)
	GetByID(id string) (*domain.Memo, error)
	Search(query string) ([]*domain.Memo, error)
	FindByTag(tag string) ([]*domain.Memo, error)
	Create(title, content string, tags []string, language string) (*domain.Memo, error)
	Update(id, title, content string, tags []string, language string) (*domain.Memo, error)
	Delete(id string) error
}
```

**Step 2: Verify**

```bash
go vet ./internal/application/usecase/...
```

Expected: no output

---

## Task 4: Create `internal/application/interacter/memo_interacter.go`

**Files:**
- Create: `internal/application/interacter/memo_interacter.go`

This replaces `usecase/memo_usecase.go`. Depends only on `application/gateway` (interface) and `domain`.

**Step 1: Create the file**

```go
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
```

**Step 2: Verify compile**

```bash
go vet ./internal/application/interacter/...
```

Expected: no output

---

## Task 5: Create `internal/adapter/gateway/sqlite_memo_gateway.go`

**Files:**
- Create: `internal/adapter/gateway/sqlite_memo_gateway.go`

This replaces `infrastructure/sqlite_memo_repository.go`. Implements `application/gateway.MemoGateway`.
Note: Both `application/gateway` and `adapter/gateway` use package name `gateway` — import alias is needed in files using both.

**Step 1: Create the file**

```go
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
```

**Step 2: Verify (interface check will catch any mismatch)**

```bash
go vet ./internal/adapter/gateway/...
```

Expected: no output

---

## Task 6: Create `internal/adapter/presenter/memo_presenter.go`

**Files:**
- Create: `internal/adapter/presenter/memo_presenter.go`

Extracted from `interface/handler/memo_handler.go` — the JSON write helpers become a standalone package.

**Step 1: Create the file**

```go
// internal/adapter/presenter/memo_presenter.go
package presenter

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, map[string]string{"error": msg})
}
```

**Step 2: Verify**

```bash
go vet ./internal/adapter/presenter/...
```

Expected: no output

---

## Task 7: Create `internal/adapter/controller/memo_controller.go`

**Files:**
- Create: `internal/adapter/controller/memo_controller.go`

The controller parses HTTP input and calls the use case. It returns domain results (or errors) — it does NOT write the HTTP response (that's the presenter's job).

**Step 1: Create the file**

```go
// internal/adapter/controller/memo_controller.go
package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	"tech-memo/internal/application/usecase"
	"tech-memo/internal/domain"
)

type MemoController struct {
	uc usecase.MemoUseCase
}

func NewMemoController(uc usecase.MemoUseCase) *MemoController {
	return &MemoController{uc: uc}
}

func (c *MemoController) List(r *http.Request) ([]*domain.Memo, error) {
	q := r.URL.Query().Get("q")
	tag := r.URL.Query().Get("tag")
	switch {
	case q != "":
		return c.uc.Search(q)
	case tag != "":
		return c.uc.FindByTag(tag)
	default:
		return c.uc.GetAll()
	}
}

func (c *MemoController) GetByID(r *http.Request) (*domain.Memo, error) {
	id := extractID(r.URL.Path)
	return c.uc.GetByID(id)
}

type MemoRequest struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Tags     []string `json:"tags"`
	Language string   `json:"language"`
}

func (c *MemoController) Create(r *http.Request) (*domain.Memo, error) {
	var req MemoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	if req.Tags == nil {
		req.Tags = []string{}
	}
	return c.uc.Create(req.Title, req.Content, req.Tags, req.Language)
}

func (c *MemoController) Update(r *http.Request) (*domain.Memo, error) {
	id := extractID(r.URL.Path)
	var req MemoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	if req.Tags == nil {
		req.Tags = []string{}
	}
	return c.uc.Update(id, req.Title, req.Content, req.Tags, req.Language)
}

func (c *MemoController) Delete(r *http.Request) error {
	id := extractID(r.URL.Path)
	return c.uc.Delete(id)
}

func extractID(path string) string {
	parts := strings.Split(strings.TrimSuffix(path, "/"), "/")
	return parts[len(parts)-1]
}
```

**Step 2: Verify**

```bash
go vet ./internal/adapter/controller/...
```

Expected: no output

---

## Task 8: Create `internal/infrastructure/api/handler.go`

**Files:**
- Create: `internal/infrastructure/api/handler.go`

The HTTP handler orchestrates controller (input) and presenter (output).

**Step 1: Create the file**

```go
// internal/infrastructure/api/handler.go
package api

import (
	"net/http"
	"strings"

	"tech-memo/internal/adapter/controller"
	"tech-memo/internal/adapter/presenter"
)

type MemoHandler struct {
	ctrl *controller.MemoController
}

func NewMemoHandler(ctrl *controller.MemoController) *MemoHandler {
	return &MemoHandler{ctrl: ctrl}
}

func (h *MemoHandler) List(w http.ResponseWriter, r *http.Request) {
	result, err := h.ctrl.List(r)
	if err != nil {
		presenter.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	presenter.WriteJSON(w, http.StatusOK, result)
}

func (h *MemoHandler) Get(w http.ResponseWriter, r *http.Request) {
	memo, err := h.ctrl.GetByID(r)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			presenter.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		presenter.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	presenter.WriteJSON(w, http.StatusOK, memo)
}

func (h *MemoHandler) Create(w http.ResponseWriter, r *http.Request) {
	memo, err := h.ctrl.Create(r)
	if err != nil {
		presenter.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	presenter.WriteJSON(w, http.StatusCreated, memo)
}

func (h *MemoHandler) Update(w http.ResponseWriter, r *http.Request) {
	memo, err := h.ctrl.Update(r)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			presenter.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		presenter.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	presenter.WriteJSON(w, http.StatusOK, memo)
}

func (h *MemoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.ctrl.Delete(r); err != nil {
		if strings.Contains(err.Error(), "not found") {
			presenter.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		presenter.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
```

**Step 2: Verify**

```bash
go vet ./internal/infrastructure/api/...
```

Expected: may fail until route.go and di.go exist in the same package — proceed to next task.

---

## Task 9: Create `internal/infrastructure/api/route.go`

**Files:**
- Create: `internal/infrastructure/api/route.go`

**Step 1: Create the file**

```go
// internal/infrastructure/api/route.go
package api

import (
	"encoding/json"
	"net/http"
)

func newRouter(h *MemoHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	mux.HandleFunc("/api/memos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.List(w, r)
		case http.MethodPost:
			h.Create(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/memos/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.Get(w, r)
		case http.MethodPut:
			h.Update(w, r)
		case http.MethodDelete:
			h.Delete(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}
```

---

## Task 10: Create `internal/infrastructure/api/di.go`

**Files:**
- Create: `internal/infrastructure/api/di.go`

The DI container wires all layers together. This is the only file that imports across all layers.

**Step 1: Create the file**

```go
// internal/infrastructure/api/di.go
package api

import (
	"net/http"

	dbgateway "tech-memo/internal/adapter/gateway"
	"tech-memo/internal/adapter/controller"
	"tech-memo/internal/application/interacter"
)

func BuildApp(dbPath string) (http.Handler, error) {
	gw, err := dbgateway.NewSQLiteMemoGateway(dbPath)
	if err != nil {
		return nil, err
	}

	uc := interacter.NewMemoInteracter(gw)
	ctrl := controller.NewMemoController(uc)
	h := NewMemoHandler(ctrl)
	return newRouter(h), nil
}
```

**Step 2: Verify whole api package**

```bash
go vet ./internal/infrastructure/api/...
```

Expected: no output

---

## Task 11: Create `internal/infrastructure/middleware/auth.go`

**Files:**
- Create: `internal/infrastructure/middleware/auth.go`

Placeholder — no auth logic yet, just the package and a stub type.

**Step 1: Create the file**

```go
// internal/infrastructure/middleware/auth.go
package middleware

import "net/http"

// AuthMiddleware is a stub for future authentication middleware.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
```

---

## Task 12: Create `internal/helper/util/util.go`

**Files:**
- Create: `internal/helper/util/util.go`

Placeholder package for shared utilities.

**Step 1: Create the file**

```go
// internal/helper/util/util.go
package util
```

---

## Task 13: Create `cmd/main.go`

**Files:**
- Create: `cmd/main.go`

**Step 1: Create the file**

```go
// cmd/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"tech-memo/internal/infrastructure/api"
)

func main() {
	port := getEnv("PORT", "8080")
	dbPath := getEnv("DB_PATH", "tech_memo.db")

	handler, err := api.BuildApp(dbPath)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Tech Memo API server starting on http://localhost%s", addr)
	log.Printf("Database: %s", dbPath)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
```

---

## Task 14: Verify full build compiles

**Step 1: Run go build**

```bash
cd C:/Users/PC1043/Documents/tech-memo
go build ./...
```

Expected: no errors. If errors occur, fix import paths before proceeding.

**Step 2: Run go vet on all new packages**

```bash
go vet ./cmd/... ./internal/...
```

Expected: no output

---

## Task 15: Delete old files

Only delete after Task 14 passes with no errors.

**Step 1: Remove old top-level packages**

```bash
rm -rf domain/ usecase/ infrastructure/ interface/
rm main.go
```

**Step 2: Verify build still passes**

```bash
go build ./...
```

Expected: no errors

**Step 3: Update Makefile to use new entry point**

Modify `Makefile` — change the build/run target from `go run .` or `go run main.go` to:

```makefile
run:
	go run ./cmd/main.go

build:
	go build -o tech-memo ./cmd/main.go
```

**Step 4: Smoke test**

```bash
# Start the server in one terminal
go run ./cmd/main.go

# In another terminal, verify health check
curl http://localhost:8080/health
# Expected: {"status":"ok"}

# Create a memo
curl -X POST http://localhost:8080/api/memos \
  -H "Content-Type: application/json" \
  -d '{"title":"test","content":"hello world","tags":["go"],"language":"go"}'
# Expected: 201 with memo object

# List all
curl http://localhost:8080/api/memos
# Expected: 200 with array containing the created memo
```

---

## Final Directory Structure

```
tech-memo/
├── cmd/
│   └── main.go
├── internal/
│   ├── adapter/
│   │   ├── controller/
│   │   │   └── memo_controller.go
│   │   ├── gateway/
│   │   │   └── sqlite_memo_gateway.go
│   │   └── presenter/
│   │       └── memo_presenter.go
│   ├── application/
│   │   ├── gateway/
│   │   │   └── memo_gateway.go
│   │   ├── interacter/
│   │   │   └── memo_interacter.go
│   │   └── usecase/
│   │       └── memo_usecase.go
│   ├── domain/
│   │   └── memo.go
│   ├── helper/
│   │   └── util/
│   │       └── util.go
│   └── infrastructure/
│       ├── api/
│       │   ├── di.go
│       │   ├── handler.go
│       │   └── route.go
│       └── middleware/
│           └── auth.go
├── docs/
│   ├── plans/
│   │   └── 2026-03-16-restructure-clean-architecture.md
│   └── spec.md
├── go.mod
└── Makefile
```
