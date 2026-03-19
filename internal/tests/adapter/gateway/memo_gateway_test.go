package gateway_test

import (
	"testing"
	"time"

	adaptergateway "tech-memo/internal/adapter/gateway"
	"tech-memo/internal/domain"
	"tech-memo/internal/infrastructure/persistence/sqlite"
)

func setupMemoDB(t *testing.T) *adaptergateway.SQLiteMemoGateway {
	t.Helper()
	db, err := sqlite.Open(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return adaptergateway.NewSQLiteMemoGateway(db)
}

func TestMemoGateway_SaveAndFindByID(t *testing.T) {
	gw := setupMemoDB(t)
	memo := &domain.Memo{
		ID:        "memo-1",
		UserID:    "user-1",
		Title:     "テストメモ",
		Content:   "内容",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := gw.Save(memo); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := gw.FindByID("memo-1")
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if got == nil {
		t.Fatal("memo not found")
	}
	if got.Title != memo.Title {
		t.Errorf("Title: got %q, want %q", got.Title, memo.Title)
	}
}

func TestMemoGateway_Delete_SoftDelete(t *testing.T) {
	gw := setupMemoDB(t)
	memo := &domain.Memo{
		ID:        "memo-2",
		UserID:    "user-1",
		Title:     "削除テスト",
		Content:   "内容",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_ = gw.Save(memo)
	if err := gw.Delete("memo-2"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	got, err := gw.FindByID("memo-2")
	if err != nil {
		t.Fatalf("FindByID after delete: %v", err)
	}
	if got != nil {
		t.Error("削除済みメモが取得できてしまう")
	}
}

func TestMemoGateway_Search(t *testing.T) {
	gw := setupMemoDB(t)
	memos := []*domain.Memo{
		{ID: "m1", UserID: "u1", Title: "Goチュートリアル", Content: "基礎", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "m2", UserID: "u1", Title: "Python入門", Content: "基礎", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	for _, m := range memos {
		_ = gw.Save(m)
	}
	results, err := gw.Search("u1", "Go")
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Search結果: got %d, want 1", len(results))
	}
}
