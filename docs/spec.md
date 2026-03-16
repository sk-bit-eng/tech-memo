# Tech Memo API 仕様書

> 作成日: 2026-03-04

---

## 目次

1. [アプリ概要](#1-アプリ概要)
2. [技術スタック](#2-技術スタック)
3. [アーキテクチャ](#3-アーキテクチャ)
4. [データ構造](#4-データ構造)
5. [API仕様](#5-api仕様)
6. [エラーレスポンス](#6-エラーレスポンス)
7. [環境変数](#7-環境変数)
8. [起動方法](#8-起動方法)
9. [実装上の制約・注意事項](#9-実装上の制約注意事項)

---

## 1. アプリ概要

技術メモ・コードスニペットをREST APIで管理するGoアプリケーション。
タイトル・本文・タグ・プログラミング言語を記録し、全文検索・タグ検索に対応する。

---

## 2. 技術スタック

| 項目 | 内容 |
|------|------|
| 言語 | Go 1.22 |
| DB | SQLite（`tech_memo.db` ファイルで永続化） |
| DBドライバ | `github.com/mattn/go-sqlite3` v1.14.22（CGO使用） |
| ID生成 | `github.com/google/uuid` v1.6.0 |
| HTTPサーバー | Go標準ライブラリ `net/http` |
| アーキテクチャ | クリーンアーキテクチャ |

---

## 3. アーキテクチャ

### 層構成

```
tech-memo/
├── main.go                          # DIコンテナ・エントリーポイント
├── domain/
│   ├── memo.go                      # Memoエンティティ・バリデーション
│   └── memo_repository.go           # MemoRepository インターフェース定義
├── usecase/
│   └── memo_usecase.go              # ビジネスロジック（CRUD・検索）
├── infrastructure/
│   └── sqlite_memo_repository.go    # SQLiteによるリポジトリ実装
└── interface/
    ├── handler/memo_handler.go      # HTTPリクエスト/レスポンス処理
    └── router/router.go             # URLルーティング定義
```

### 依存関係図

```
┌─────────────┐     ┌─────────────┐     ┌──────────────────┐
│  Interface  │────▶│   UseCase   │────▶│     Domain       │
│  (handler,  │     │             │     │ (Memo, Repository│
│   router)   │     │             │     │   Interface)     │
└─────────────┘     └─────────────┘     └──────────────────┘
                                                 ▲
                                                 │ implements
                                        ┌────────────────┐
                                        │ Infrastructure │
                                        │  (SQLite実装)  │
                                        └────────────────┘
```

依存の方向は常に **外側 → 内側（Domain）** のみ。Domain層は外部パッケージに一切依存しない。

### 各層の責務

| 層 | パッケージ | 責務 |
|----|-----------|------|
| Domain | `domain` | エンティティ定義・バリデーション・リポジトリIF |
| UseCase | `usecase` | ビジネスロジック（ID生成・存在確認・エラー変換） |
| Infrastructure | `infrastructure` | SQLite接続・テーブル作成・SQL実行 |
| Interface | `interface/handler`, `interface/router` | HTTPパース・ルーティング・JSONレスポンス |
| main | `main` | 全層のDI組み立て・サーバー起動 |

---

## 4. データ構造

### Memo エンティティ

| フィールド | 型 | JSON キー | 説明 |
|-----------|-----|-----------|------|
| ID | `string` | `id` | UUID v4（自動採番） |
| Title | `string` | `title` | タイトル（必須・200文字以内） |
| Content | `string` | `content` | 本文・コードスニペット（必須） |
| Tags | `[]string` | `tags` | タグ一覧（省略時は空配列） |
| Language | `string` | `language` | プログラミング言語（例: `go`, `sql`） |
| CreatedAt | `time.Time` | `created_at` | 作成日時（RFC3339Nano形式） |
| UpdatedAt | `time.Time` | `updated_at` | 更新日時（RFC3339Nano形式） |

### SQLite テーブル定義

```sql
CREATE TABLE IF NOT EXISTS memos (
    id         TEXT PRIMARY KEY,
    title      TEXT NOT NULL,
    content    TEXT NOT NULL,
    tags       TEXT NOT NULL DEFAULT '',   -- カンマ区切り文字列で保存
    language   TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL,          -- RFC3339Nano文字列
    updated_at DATETIME NOT NULL           -- RFC3339Nano文字列
);
```

> **タグの保存形式**: `["go", "api", "rest"]` → `"go,api,rest"` として TEXT 列に保存。取得時に再分割して `[]string` に戻す。

### バリデーションルール

| フィールド | ルール |
|-----------|--------|
| title | 空文字・スペースのみ不可。200文字以内 |
| content | 空文字・スペースのみ不可 |
| tags | 省略可（`null` の場合は空配列として扱う） |
| language | 省略可 |

---

## 5. API仕様

**ベースURL**: `http://localhost:8080`
**Content-Type**: `application/json`
**ソート順**: 全リスト系エンドポイントは `updated_at DESC` で返却

---

### 5.1 ヘルスチェック

#### `GET /health`

サーバーの死活確認。

**レスポンス** `200 OK`
```json
{
  "status": "ok"
}
```

---

### 5.2 メモ一覧取得

#### `GET /api/memos`

全件取得。クエリパラメータで絞り込み可能。

**クエリパラメータ**

| パラメータ | 型 | 説明 |
|-----------|-----|------|
| `q` | string | 全文検索（`title` または `content` の部分一致） |
| `tag` | string | タグ完全一致検索 |

> `q` と `tag` を同時指定した場合は `q` が優先される。

**レスポンス** `200 OK`
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Goのgoroutine基礎",
    "content": "go func() { ... }()",
    "tags": ["go", "concurrency"],
    "language": "go",
    "created_at": "2026-03-04T10:00:00.000000000Z",
    "updated_at": "2026-03-04T10:00:00.000000000Z"
  }
]
```

> 件数ゼロの場合は空配列 `[]` を返す。

**例: 全文検索**
```
GET /api/memos?q=goroutine
```

**例: タグ検索**
```
GET /api/memos?tag=go
```

---

### 5.3 メモ1件取得

#### `GET /api/memos/{id}`

**パスパラメータ**

| パラメータ | 型 | 説明 |
|-----------|-----|------|
| `id` | string | メモのUUID |

**レスポンス** `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Goのgoroutine基礎",
  "content": "go func() { ... }()",
  "tags": ["go", "concurrency"],
  "language": "go",
  "created_at": "2026-03-04T10:00:00.000000000Z",
  "updated_at": "2026-03-04T10:00:00.000000000Z"
}
```

**エラー** `404 Not Found` — 指定IDが存在しない場合

---

### 5.4 メモ新規作成

#### `POST /api/memos`

**リクエストボディ**

```json
{
  "title": "Goのgoroutine基礎",
  "content": "go func() { ... }()",
  "tags": ["go", "concurrency"],
  "language": "go"
}
```

| フィールド | 必須 | 説明 |
|-----------|------|------|
| `title` | ✅ | 200文字以内 |
| `content` | ✅ | 本文 |
| `tags` | — | 省略時は空配列 |
| `language` | — | 省略時は空文字 |

**レスポンス** `201 Created`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Goのgoroutine基礎",
  "content": "go func() { ... }()",
  "tags": ["go", "concurrency"],
  "language": "go",
  "created_at": "2026-03-04T10:00:00.000000000Z",
  "updated_at": "2026-03-04T10:00:00.000000000Z"
}
```

**エラー** `400 Bad Request` — バリデーション失敗またはJSONパース失敗

---

### 5.5 メモ更新

#### `PUT /api/memos/{id}`

全フィールドを上書き更新する（部分更新不可）。

**パスパラメータ**

| パラメータ | 型 | 説明 |
|-----------|-----|------|
| `id` | string | メモのUUID |

**リクエストボディ** — POST と同一スキーマ

```json
{
  "title": "Goのgoroutine基礎（改訂版）",
  "content": "go func() { fmt.Println(\"hello\") }()",
  "tags": ["go", "concurrency", "goroutine"],
  "language": "go"
}
```

**レスポンス** `200 OK` — 更新後のメモオブジェクト（`updated_at` が更新される）

**エラー**
- `400 Bad Request` — バリデーション失敗
- `404 Not Found` — 指定IDが存在しない

---

### 5.6 メモ削除

#### `DELETE /api/memos/{id}`

**パスパラメータ**

| パラメータ | 型 | 説明 |
|-----------|-----|------|
| `id` | string | メモのUUID |

**レスポンス** `204 No Content` — ボディなし

**エラー** `404 Not Found` — 指定IDが存在しない場合

---

## 6. エラーレスポンス

全エラーは以下の統一フォーマットで返す。

```json
{
  "error": "エラーメッセージ"
}
```

### HTTPステータスコード一覧

| コード | 意味 | 発生ケース |
|--------|------|-----------|
| `200` | OK | GET・PUT 成功 |
| `201` | Created | POST 成功 |
| `204` | No Content | DELETE 成功 |
| `400` | Bad Request | JSONパース失敗・バリデーション失敗 |
| `404` | Not Found | 指定IDのメモが存在しない |
| `405` | Method Not Allowed | 未対応HTTPメソッド |
| `500` | Internal Server Error | DB操作エラー |

---

## 7. 環境変数

| 変数名 | デフォルト値 | 説明 |
|--------|------------|------|
| `PORT` | `8080` | リッスンポート |
| `DB_PATH` | `tech_memo.db` | SQLiteファイルのパス |

---

## 8. 起動方法

### 前提条件

- Go 1.22 以上
- GCC（`go-sqlite3` のCGOビルドに必要）
  - Windows: [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) または WSL を推奨
  - macOS: Xcode Command Line Tools（`xcode-select --install`）
  - Linux: `gcc` パッケージ

### コマンド

```bash
cd ~/tech-memo

# 初回のみ: 依存パッケージ取得
go mod tidy

# 開発サーバー起動
make run
# → http://localhost:8080 で起動

# バイナリビルド
make build
# → ./tech-memo バイナリが生成される

# 環境変数でポート・DBパスを変更する場合
PORT=9090 DB_PATH=/data/memo.db make run
```

### 動作確認 (curl)

```bash
# ヘルスチェック
curl http://localhost:8080/health

# メモ作成
curl -X POST http://localhost:8080/api/memos \
  -H "Content-Type: application/json" \
  -d '{
    "title": "SELECT文の基本",
    "content": "SELECT * FROM users WHERE id = 1;",
    "tags": ["sql", "database"],
    "language": "sql"
  }'

# 全件取得
curl http://localhost:8080/api/memos

# キーワード検索
curl "http://localhost:8080/api/memos?q=SELECT"

# タグ検索
curl "http://localhost:8080/api/memos?tag=sql"

# ID指定取得（IDは作成時のレスポンスから取得）
curl http://localhost:8080/api/memos/{id}

# 更新
curl -X PUT http://localhost:8080/api/memos/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "title": "SELECT文の基本（改訂）",
    "content": "SELECT id, name FROM users WHERE id = 1;",
    "tags": ["sql", "database"],
    "language": "sql"
  }'

# 削除
curl -X DELETE http://localhost:8080/api/memos/{id}
```

---

## 9. 実装上の制約・注意事項

| 項目 | 内容 |
|------|------|
| タグ検索 | タグはDBにカンマ区切りで保存するため、カンマを含むタグ名は使用不可 |
| PUT の部分更新 | PUT は全フィールド置換。`tags`・`language` を省略すると空になる |
| 全文検索のケース | `LIKE` による部分一致のため大文字小文字を区別しない（SQLiteのデフォルト） |
| タイムゾーン | `created_at`・`updated_at` はサーバーのローカル時刻で記録される |
| 並行アクセス | SQLiteの書き込みロックにより、高並行書き込みには向かない |
| CGO依存 | `go-sqlite3` はCGOを使用するため、クロスコンパイル時は注意が必要 |
