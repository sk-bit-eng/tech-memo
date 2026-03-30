package memo

import (
	"os"
	"testing"
	"time"

	"tech-memo/internal/domain"
	sqlserverinfra "tech-memo/internal/infrastructure/persistence/sqlserver"
)

func setup(t *testing.T) *Repository {
	t.Helper()
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "sqlserver://sa:Test@1234@localhost:1433?database=tech_memo"
	}
	db, err := sqlserverinfra.Open(dsn)
	if err != nil {
		t.Skipf("SQL Server unavailable: %v", err)
	}
	r := NewRepository(db)
	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Exec("DELETE FROM memos")
	})
	return r
}

func TestSaveAndFindByID(t *testing.T) {
	r := setup(t)
	memo := &domain.Memo{
		ID: "memo-1", UserID: "user-1", Title: "テストメモ", Content: "内容",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	if err := r.Save(memo); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := r.FindByID("memo-1")
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if got == nil || got.Title != memo.Title {
		t.Errorf("Title mismatch: got %v", got)
	}
}

func TestDelete_SoftDelete(t *testing.T) {
	r := setup(t)
	memo := &domain.Memo{
		ID: "memo-del", UserID: "u1", Title: "削除", Content: "",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	_ = r.Save(memo)
	if err := r.Delete("memo-del"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	got, err := r.FindByID("memo-del")
	if err != nil {
		t.Fatalf("FindByID after delete: %v", err)
	}
	if got != nil {
		t.Error("削除済みメモが取得できてしまう")
	}
}

func TestSearch(t *testing.T) {
	r := setup(t)
	memos := []*domain.Memo{
		{ID: "m1", UserID: "u1", Title: "Goチュートリアル", Content: "基礎", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "m2", UserID: "u1", Title: "Python入門", Content: "基礎", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	for _, m := range memos {
		_ = r.Save(m)
	}
	results, err := r.Search("u1", "Go")
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Search結果: got %d, want 1", len(results))
	}
}

func TestFindByUserIDAndCategory(t *testing.T) {
	r := setup(t)
	memos := []*domain.Memo{
		{ID: "u1-1", UserID: "u1", Title: "A", Content: "1", CategoryID: "c1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "u1-2", UserID: "u1", Title: "B", Content: "2", CategoryID: "c2", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "u2-1", UserID: "u2", Title: "C", Content: "3", CategoryID: "c1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	for _, m := range memos {
		_ = r.Save(m)
	}
	gotByUser, err := r.FindByUserID("u1")
	if err != nil {
		t.Fatalf("FindByUserID: %v", err)
	}
	if len(gotByUser) != 2 {
		t.Fatalf("FindByUserID: got %d, want 2", len(gotByUser))
	}

	gotByCategory, err := r.FindByCategory("u1", "c2")
	if err != nil {
		t.Fatalf("FindByCategory: %v", err)
	}
	if len(gotByCategory) != 1 {
		t.Fatalf("FindByCategory: got %d, want 1", len(gotByCategory))
	}
	if gotByCategory[0].ID != "u1-2" {
		t.Errorf("FindByCategory: expected u1-2, got %s", gotByCategory[0].ID)
	}
}

func TestUpdateAndSoftDeleteAndFindByID(t *testing.T) {
	r := setup(t)
	memo := &domain.Memo{ID: "memo-upd", UserID: "u9", Title: "Old", Content: "C", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := r.Save(memo); err != nil {
		t.Fatalf("Save: %v", err)
	}

	memo.Title = "New"
	if err := r.Update(memo); err != nil {
		t.Fatalf("Update: %v", err)
	}
	got, err := r.FindByID("memo-upd")
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if got == nil || got.Title != "New" {
		t.Fatalf("Update反映: got %v", got)
	}

	if err := r.Delete("memo-upd"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	gotDeleted, err := r.FindByID("memo-upd")
	if err != nil {
		t.Fatalf("FindByID after delete: %v", err)
	}
	if gotDeleted != nil {
		t.Fatalf("soft delete後にnil期待: got %+v", gotDeleted)
	}
}

func TestFindByIDNotFound(t *testing.T) {
	r := setup(t)
	got, err := r.FindByID("unknown")
	if err != nil {
		t.Fatalf("FindByID unknown: %v", err)
	}
	if got != nil {
		t.Fatal("FindByID unknown が nil であること")
	}
}
