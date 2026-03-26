package todo

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
		sqlDB.Exec("DELETE FROM todos")
	})
	return r
}

func TestSaveAndFindByID(t *testing.T) {
	r := setup(t)
	todo := &domain.Todo{
		ID: "todo-1", UserID: "user-1", Title: "テストTodo", Content: "内容",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	if err := r.Save(todo); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := r.FindByID("todo-1")
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if got == nil || got.Title != todo.Title {
		t.Errorf("Title mismatch: got %v", got)
	}
}

func TestFindPendingAndCompleted(t *testing.T) {
	r := setup(t)
	now := time.Now()
	todos := []*domain.Todo{
		{ID: "t1", UserID: "u1", Title: "未完了", Content: "", CreatedAt: now, UpdatedAt: now},
		{ID: "t2", UserID: "u1", Title: "完了済み", Content: "", CompletedAt: &now, CreatedAt: now, UpdatedAt: now},
	}
	for _, td := range todos {
		_ = r.Save(td)
	}
	pending, err := r.FindPending("u1")
	if err != nil {
		t.Fatalf("FindPending: %v", err)
	}
	if len(pending) != 1 || pending[0].ID != "t1" {
		t.Errorf("FindPending: got %d items", len(pending))
	}
	completed, err := r.FindCompleted("u1")
	if err != nil {
		t.Fatalf("FindCompleted: %v", err)
	}
	if len(completed) != 1 || completed[0].ID != "t2" {
		t.Errorf("FindCompleted: got %d items", len(completed))
	}
}

func TestDelete_SoftDelete(t *testing.T) {
	r := setup(t)
	todo := &domain.Todo{
		ID: "todo-del", UserID: "u1", Title: "削除", Content: "",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	_ = r.Save(todo)
	if err := r.Delete("todo-del"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	got, err := r.FindByID("todo-del")
	if err != nil {
		t.Fatalf("FindByID after delete: %v", err)
	}
	if got != nil {
		t.Error("削除済みTodoが取得できてしまう")
	}
}
