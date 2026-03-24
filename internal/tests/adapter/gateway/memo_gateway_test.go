// internal/tests/adapter/gateway/memo_gateway_test.go
package gateway_test

import (
	"os"
	"testing"
	"time"

	adaptergateway "tech-memo/internal/adapter/gateway"
	"tech-memo/internal/domain"
	sqlserverinfra "tech-memo/internal/infrastructure/persistence/sqlserver"
)

func setupMemoGW(t *testing.T) *adaptergateway.GORMMemoGateway {
	t.Helper()
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "sqlserver://sa:TechMemo123!@localhost:1433?database=tech_memo"
	}
	db, err := sqlserverinfra.Open(dsn)
	if err != nil {
		t.Skipf("SQL Server unavailable: %v", err)
	}
	gw := adaptergateway.NewGORMMemoGateway(db)
	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Exec("DELETE FROM memos")
	})
	return gw
}

func TestMemoGateway_SaveAndFindByID(t *testing.T) {
	gw := setupMemoGW(t)
	memo := &domain.Memo{
		ID: "memo-1", UserID: "user-1", Title: "テストメモ", Content: "内容",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	if err := gw.Save(memo); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := gw.FindByID("memo-1")
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if got == nil || got.Title != memo.Title {
		t.Errorf("Title mismatch: got %v", got)
	}
}

func TestMemoGateway_Delete_SoftDelete(t *testing.T) {
	gw := setupMemoGW(t)
	memo := &domain.Memo{
		ID: "memo-del", UserID: "u1", Title: "削除", Content: "",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	_ = gw.Save(memo)
	if err := gw.Delete("memo-del"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	got, err := gw.FindByID("memo-del")
	if err != nil {
		t.Fatalf("FindByID after delete: %v", err)
	}
	if got != nil {
		t.Error("削除済みメモが取得できてしまう")
	}
}

func TestMemoGateway_Search(t *testing.T) {
	gw := setupMemoGW(t)
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
