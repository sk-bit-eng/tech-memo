# ユースケース・インタラクター統合 設計ドキュメント

**日付:** 2026-03-17

## 概要

クリーンアーキテクチャの `application` 層を整理し、以下2点を変更する。

1. `cmd/main.go` を `cmd/api/main.go` へ移動
2. `application/usecase/`（インターフェース）と `application/interacter/`（実装）を `application/usecase/memo.go` に統合

---

## 変更前の構成

```
cmd/
  main.go

internal/
  application/
    usecase/
      memo_usecase.go     ← インターフェースのみ
    interacter/
      memo_interacter.go  ← 実装のみ（別ディレクトリ）
    gateway/              ← <I> インターフェースのみ
```

## 変更後の構成

```
cmd/
  api/
    main.go               ← 移動（中身は変更なし）

internal/
  application/
    usecase/
      memo.go             ← インターフェース + 実装を同一ファイルに統合
    gateway/              ← <I> インターフェースのみ（変更なし）
```

---

## 設計方針

### なぜ同一ファイルに統合するか

- Go では「インターフェース定義と実装を近くに置く」スタイルが一般的
- `usecase` と `interacter` を分離するメリットが薄い（1対1の関係）
- ファイル数を減らし、把握しやすくする

### ファイル構造

```go
// internal/application/usecase/memo.go

package usecase

// ---- インターフェース（<I> の役割）----

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

// ... メソッド実装
```

### cmd/api/ ディレクトリ

将来のエントリーポイント追加に備える:

```
cmd/
  api/
    main.go     ← REST API サーバー（現在）
  worker/
    main.go     ← バッチ・ワーカー（将来）
```

---

## 影響ファイル

| ファイル | 変更内容 |
|---|---|
| `cmd/api/main.go` | `cmd/main.go` から移動（中身は変更なし） |
| `internal/application/usecase/memo.go` | 新規作成（usecase + interacter を統合） |
| `internal/infrastructure/api/di.go` | import を `interacter` → `usecase` に変更、`NewMemoInteracter` の呼び出しを `usecase.NewMemoInteracter` に変更 |
| `internal/application/usecase/memo_usecase.go` | 削除 |
| `internal/application/interacter/memo_interacter.go` | 削除 |
| `internal/application/interacter/` | ディレクトリ削除 |

---

## 依存関係（変更なし）

```
cmd/api → infrastructure/api → adapter/* → application/usecase → application/gateway(<I>) → domain
                                              ↑
                               adapter/gateway が application/gateway(<I>) を実装
```
