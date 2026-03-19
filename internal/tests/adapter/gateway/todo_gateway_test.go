package gateway_test

import (
	"testing"
	"time"

	adaptergateway "tech-memo/internal/adapter/gateway"
	"tech-memo/internal/domain"
	"tech-memo/internal/infrastructure/persistence/sqlite"
)

func setupTodoDB(t *testing.T) *adaptergateway.SQLiteTodoGateway {
	t.Helper()
	db, err := sqlite.Open(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return adaptergateway.NewSQLiteTodoGateway(db)
}

func TestTodoGateway_SaveAndFindByID(t *testing.T) {
	gw := setupTodoDB(t)
	todo := &domain.Todo{
		ID:        "todo-1",
		UserID:    "user-1",
		Title:     "テストTodo",
		Content:   "内容",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := gw.Save(todo); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := gw.FindByID("todo-1")
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if got == nil {
		t.Fatal("todo not found")
	}
	if got.Title != todo.Title {
		t.Errorf("Title: got %q, want %q", got.Title, todo.Title)
	}
}

func TestTodoGateway_FindPendingAndCompleted(t *testing.T) {
	gw := setupTodoDB(t)
	now := time.Now()
	todos := []*domain.Todo{
		{ID: "t1", UserID: "u1", Title: "未完了", Content: "", CreatedAt: now, UpdatedAt: now},
		{ID: "t2", UserID: "u1", Title: "完了済み", Content: "", CompletedAt: &now, CreatedAt: now, UpdatedAt: now},
	}
	for _, td := range todos {
		_ = gw.Save(td)
	}

	pending, err := gw.FindPending("u1")
	if err != nil {
		t.Fatalf("FindPending: %v", err)
	}
	if len(pending) != 1 || pending[0].ID != "t1" {
		t.Errorf("FindPending: got %d items", len(pending))
	}

	completed, err := gw.FindCompleted("u1")
	if err != nil {
		t.Fatalf("FindCompleted: %v", err)
	}
	if len(completed) != 1 || completed[0].ID != "t2" {
		t.Errorf("FindCompleted: got %d items", len(completed))
	}
}

func TestTodoGateway_Delete_SoftDelete(t *testing.T) {
	gw := setupTodoDB(t)
	todo := &domain.Todo{
		ID: "todo-del", UserID: "u1", Title: "削除", Content: "",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	_ = gw.Save(todo)
	if err := gw.Delete("todo-del"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	got, err := gw.FindByID("todo-del")
	if err != nil {
		t.Fatalf("FindByID after delete: %v", err)
	}
	if got != nil {
		t.Error("削除済みTodoが取得できてしまう")
	}
}
