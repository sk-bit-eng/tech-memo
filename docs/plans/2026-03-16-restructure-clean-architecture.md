# クリーンアーキテクチャ再構成 実装計画

> **Claude向け:** 必須サブスキル: superpowers:executing-plans を使用してタスクを1つずつ実行すること。

**目標:** 既存のtech-memo Goプロジェクトを `internal/` 配下の階層構造に再編成し、Domain・Application・Adapter・Infrastructure を明確に分離する。

**アーキテクチャ:** 全コードを `internal/` 配下に置き、明示的な層構成とする。`domain`（エンティティ）、`application`（ユースケースIF・インタラクター・ゲートウェイIF）、`adapter`（コントローラー・プレゼンター・ゲートウェイ実装）、`infrastructure`（HTTPサーバー・DI・ミドルウェア）。エントリーポイントは `cmd/main.go` に移動する。

**技術スタック:** Go 1.22、SQLite（`github.com/mattn/go-sqlite3`）、UUID（`github.com/google/uuid`）、標準ライブラリ `net/http`

---

## 各層の責務

| 層 | パッケージパス | 役割 |
|---|---|---|
| Domain | `internal/domain` | エンティティ＋バリデーション。外部依存なし。 |
| Application/usecase | `internal/application/usecase` | ユースケース**インターフェース**（I = Interface の目印） |
| Application/interacter | `internal/application/interacter` | ユースケース**実装**（ゲートウェイIFを呼び出す） |
| Application/gateway | `internal/application/gateway` | リポジトリ/ゲートウェイ**インターフェース** |
| Adapter/controller | `internal/adapter/controller` | HTTPリクエスト解析 → ユースケース呼び出し |
| Adapter/presenter | `internal/adapter/presenter` | JSONレスポンス生成 |
| Adapter/gateway | `internal/adapter/gateway` | ゲートウェイIFのSQLite実装 |
| Infrastructure/api | `internal/infrastructure/api` | HTTP配線: ハンドラー・ルーター・DI |
| Infrastructure/middleware | `internal/infrastructure/middleware` | HTTPミドルウェア（認証stub） |
| Helper/util | `internal/helper/util` | 共通ユーティリティ |
| cmd | `cmd/` | エントリーポイントのみ |

## 依存関係ルール

```
cmd → infrastructure/api → adapter/* → application/* → domain
                                      ↑
                         adapter/gateway が application/gateway（IF）を実装
```

---

## タスク1: `internal/domain/memo.go` の作成

**ファイル:**
- 作成: `internal/domain/memo.go`

**手順1: ファイル作成**（`domain/memo.go` からコピー。パッケージパスのみ変更、内容は同一）

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

**手順2: 構文エラーがないことを確認**

```bash
cd C:/Users/PC1043/Documents/tech-memo
go vet ./internal/domain/...
```

期待値: 出力なし（成功）

---

## タスク2: `internal/application/gateway/memo_gateway.go` の作成

**ファイル:**
- 作成: `internal/application/gateway/memo_gateway.go`

`domain/memo_repository.go` の代替。インターフェース名を `MemoGateway`（旧: `MemoRepository`）に変更し、ゲートウェイパターンに合わせる。

**手順1: ファイル作成**

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

**手順2: 確認**

```bash
go vet ./internal/application/gateway/...
```

期待値: 出力なし

---

## タスク3: `internal/application/usecase/memo_usecase.go` の作成

**ファイル:**
- 作成: `internal/application/usecase/memo_usecase.go`

**新規ファイル** — インタラクターが実装するユースケースインターフェース。

**手順1: ファイル作成**

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

**手順2: 確認**

```bash
go vet ./internal/application/usecase/...
```

期待値: 出力なし

---

## タスク4: `internal/application/interacter/memo_interacter.go` の作成

**ファイル:**
- 作成: `internal/application/interacter/memo_interacter.go`

`usecase/memo_usecase.go` の代替。`application/gateway`（IF）と `domain` にのみ依存する。

**手順1: ファイル作成**

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

**手順2: コンパイル確認**

```bash
go vet ./internal/application/interacter/...
```

期待値: 出力なし

---

## タスク5: `internal/adapter/gateway/sqlite_memo_gateway.go` の作成

**ファイル:**
- 作成: `internal/adapter/gateway/sqlite_memo_gateway.go`

`infrastructure/sqlite_memo_repository.go` の代替。`application/gateway.MemoGateway` を実装する。
注意: `application/gateway` と `adapter/gateway` は同じパッケージ名 `gateway` を使用するため、両方をimportするファイルではエイリアスが必要。

**手順1: ファイル作成**

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

// コンパイル時インターフェース適合確認
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

**手順2: 確認**（インターフェース適合チェックが不一致を検出する）

```bash
go vet ./internal/adapter/gateway/...
```

期待値: 出力なし

---

## タスク6: `internal/adapter/presenter/memo_presenter.go` の作成

**ファイル:**
- 作成: `internal/adapter/presenter/memo_presenter.go`

`interface/handler/memo_handler.go` から抽出 — JSONレスポンス書き込みヘルパーを独立したパッケージとする。

**手順1: ファイル作成**

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

**手順2: 確認**

```bash
go vet ./internal/adapter/presenter/...
```

期待値: 出力なし

---

## タスク7: `internal/adapter/controller/memo_controller.go` の作成

**ファイル:**
- 作成: `internal/adapter/controller/memo_controller.go`

コントローラーはHTTPリクエストを解析してユースケースを呼び出す。ドメイン結果（またはエラー）を返す。HTTPレスポンスの書き込みはしない（それはプレゼンターの役割）。

**手順1: ファイル作成**

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

**手順2: 確認**

```bash
go vet ./internal/adapter/controller/...
```

期待値: 出力なし

---

## タスク8: `internal/infrastructure/api/handler.go` の作成

**ファイル:**
- 作成: `internal/infrastructure/api/handler.go`

HTTPハンドラーはコントローラー（入力）とプレゼンター（出力）を協調させる。

**手順1: ファイル作成**

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

**手順2: 確認**

```bash
go vet ./internal/infrastructure/api/...
```

期待値: route.go と di.go が同パッケージに存在するまで失敗する可能性あり — 次のタスクに進む。

---

## タスク9: `internal/infrastructure/api/route.go` の作成

**ファイル:**
- 作成: `internal/infrastructure/api/route.go`

**手順1: ファイル作成**

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

## タスク10: `internal/infrastructure/api/di.go` の作成

**ファイル:**
- 作成: `internal/infrastructure/api/di.go`

DIコンテナで全層を接続する。全層をまたいでimportするのはこのファイルのみ。

**手順1: ファイル作成**

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

**手順2: apiパッケージ全体の確認**

```bash
go vet ./internal/infrastructure/api/...
```

期待値: 出力なし

---

## タスク11: `internal/infrastructure/middleware/auth.go` の作成

**ファイル:**
- 作成: `internal/infrastructure/middleware/auth.go`

プレースホルダー — 認証ロジックなし。パッケージとstubのみ。

**手順1: ファイル作成**

```go
// internal/infrastructure/middleware/auth.go
package middleware

import "net/http"

// AuthMiddleware は将来の認証ミドルウェアのためのstubです。
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
```

---

## タスク12: `internal/helper/util/util.go` の作成

**ファイル:**
- 作成: `internal/helper/util/util.go`

共通ユーティリティのプレースホルダーパッケージ。

**手順1: ファイル作成**

```go
// internal/helper/util/util.go
package util
```

---

## タスク13: `cmd/main.go` の作成

**ファイル:**
- 作成: `cmd/main.go`

**手順1: ファイル作成**

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

## タスク14: ビルド全体の確認

**手順1: go build 実行**

```bash
cd C:/Users/PC1043/Documents/tech-memo
go build ./...
```

期待値: エラーなし。エラーが出た場合はimportパスを修正してから進む。

**手順2: 全新規パッケージの go vet 実行**

```bash
go vet ./cmd/... ./internal/...
```

期待値: 出力なし

---

## タスク15: 旧ファイルの削除

タスク14がエラーなしで通過してから削除すること。

**手順1: 旧トップレベルパッケージを削除**

```bash
rm -rf domain/ usecase/ infrastructure/ interface/
rm main.go
```

**手順2: ビルドがまだ通ることを確認**

```bash
go build ./...
```

期待値: エラーなし

**手順3: Makefileのエントリーポイントを更新**

`Makefile` を修正 — build/runターゲットを `go run .` または `go run main.go` から以下に変更:

```makefile
run:
	go run ./cmd/main.go

build:
	go build -o tech-memo ./cmd/main.go
```

**手順4: スモークテスト**

```bash
# ターミナル1でサーバー起動
go run ./cmd/main.go

# ターミナル2でヘルスチェック確認
curl http://localhost:8080/health
# 期待値: {"status":"ok"}

# メモ作成
curl -X POST http://localhost:8080/api/memos \
  -H "Content-Type: application/json" \
  -d '{"title":"test","content":"hello world","tags":["go"],"language":"go"}'
# 期待値: 201 + メモオブジェクト

# 全件取得
curl http://localhost:8080/api/memos
# 期待値: 200 + 作成したメモを含む配列
```

---

## 最終ディレクトリ構成

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
