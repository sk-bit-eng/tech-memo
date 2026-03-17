# ユースケース・インタラクター統合 実装計画

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** `application/interacter/` を廃止し、インターフェースと実装を `application/usecase/memo.go` に統合する。あわせて `cmd/main.go` を `cmd/api/main.go` へ移動する。

**Architecture:** `application/usecase/memo.go` に `MemoUseCase` インターフェースと `memoInteracter` 実装を同居させる。`di.go` のインポートパスを更新し、旧ファイルを削除する。Goのビルドが通ることで完了を確認する。

**Tech Stack:** Go 1.22、`go build`、`go vet`

---

## タスク一覧

| # | 内容 |
|---|---|
| 1 | `cmd/api/main.go` を作成（`cmd/main.go` から移動） |
| 2 | `application/usecase/memo.go` を新規作成（IF + 実装を統合） |
| 3 | `infrastructure/api/di.go` のインポートを更新 |
| 4 | 旧ファイル・ディレクトリを削除 |
| 5 | ビルド確認・コミット |

---

## タスク1: `cmd/api/main.go` の作成

**Files:**
- 作成: `cmd/api/main.go`
- 削除予定: `cmd/main.go`（タスク4で実施）

**手順1: `cmd/api/` ディレクトリを作成し `main.go` を配置**

```go
// cmd/api/main.go
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

**手順2: 構文確認**

```bash
go vet ./cmd/api/...
```

期待値: 出力なし（エラーなし）

---

## タスク2: `application/usecase/memo.go` を新規作成

**Files:**
- 作成: `internal/application/usecase/memo.go`
- 削除予定: `internal/application/usecase/memo_usecase.go`（タスク4で実施）
- 削除予定: `internal/application/interacter/memo_interacter.go`（タスク4で実施）

**手順1: 統合ファイルを作成**

```go
// internal/application/usecase/memo.go
package usecase

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	appgateway "tech-memo/internal/application/gateway"
	"tech-memo/internal/domain"
)

// ---- インターフェース ----

type MemoUseCase interface {
	GetAll() ([]*domain.Memo, error)
	GetByID(id string) (*domain.Memo, error)
	Search(query string) ([]*domain.Memo, error)
	FindByTag(tag string) ([]*domain.Memo, error)
	Create(title, content string, tags []string, language string) (*domain.Memo, error)
	Update(id, title, content string, tags []string, language string) (*domain.Memo, error)
	Delete(id string) error
}

// ---- インタラクター（実装）----

type memoInteracter struct {
	gw appgateway.MemoGateway
}

func NewMemoInteracter(gw appgateway.MemoGateway) MemoUseCase {
	return &memoInteracter{gw: gw}
}

func (uc *memoInteracter) GetAll() ([]*domain.Memo, error) {
	return uc.gw.FindAll()
}

func (uc *memoInteracter) GetByID(id string) (*domain.Memo, error) {
	memo, err := uc.gw.FindByID(id)
	if err != nil {
		return nil, err
	}
	if memo == nil {
		return nil, fmt.Errorf("memo not found: %s", id)
	}
	return memo, nil
}

func (uc *memoInteracter) Search(query string) ([]*domain.Memo, error) {
	return uc.gw.Search(query)
}

func (uc *memoInteracter) FindByTag(tag string) ([]*domain.Memo, error) {
	return uc.gw.FindByTag(tag)
}

func (uc *memoInteracter) Create(title, content string, tags []string, language string) (*domain.Memo, error) {
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

func (uc *memoInteracter) Update(id, title, content string, tags []string, language string) (*domain.Memo, error) {
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

func (uc *memoInteracter) Delete(id string) error {
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

**手順2: 構文確認**

```bash
go vet ./internal/application/usecase/...
```

期待値: 出力なし（エラーなし）

> ※ この時点では旧ファイル（`memo_usecase.go`）と型名が重複するためビルドエラーになる。次のタスク3・4で解消する。

---

## タスク3: `infrastructure/api/di.go` のインポートを更新

**Files:**
- 修正: `internal/infrastructure/api/di.go`

**手順1: インポートパスを `interacter` → `usecase` に変更**

変更前:
```go
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

変更後:
```go
import (
	"net/http"

	dbgateway "tech-memo/internal/adapter/gateway"
	"tech-memo/internal/adapter/controller"
	"tech-memo/internal/application/usecase"
)

func BuildApp(dbPath string) (http.Handler, error) {
	gw, err := dbgateway.NewSQLiteMemoGateway(dbPath)
	if err != nil {
		return nil, err
	}

	uc := usecase.NewMemoInteracter(gw)
	ctrl := controller.NewMemoController(uc)
	h := NewMemoHandler(ctrl)
	return newRouter(h), nil
}
```

---

## タスク4: 旧ファイル・ディレクトリを削除

**Files:**
- 削除: `internal/application/usecase/memo_usecase.go`
- 削除: `internal/application/interacter/memo_interacter.go`
- 削除: `internal/application/interacter/`（空になったディレクトリ）
- 削除: `cmd/main.go`

**手順1: 削除コマンドを実行**

```bash
rm internal/application/usecase/memo_usecase.go
rm internal/application/interacter/memo_interacter.go
rmdir internal/application/interacter
rm cmd/main.go
```

---

## タスク5: ビルド確認・コミット

**手順1: プロジェクト全体のビルド確認**

```bash
go build ./...
```

期待値: 出力なし（ビルド成功）

**手順2: `go vet` で静的解析**

```bash
go vet ./...
```

期待値: 出力なし（警告なし）

**手順3: 動作確認（オプション）**

```bash
go run ./cmd/api/main.go
```

期待値:
```
Tech Memo API server starting on http://localhost:8080
Database: tech_memo.db
```

**手順4: コミット**

```bash
git add -A
git commit -m "refactor: merge usecase interface and interacter into usecase/memo.go, move cmd to cmd/api"
```
