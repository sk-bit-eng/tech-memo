// internal/tests/usecase/todo_test.go
package usecase_test

import (
	"testing"

	"tech-memo/internal/application/dto"
	"tech-memo/internal/application/usecase"
	"tech-memo/internal/domain"
)

// --- モックゲートウェイ ---

type mockTodoGateway struct {
	todos   map[string]*domain.Todo
	updated []*domain.Todo
	deleted []string
}

func newMockTodoGateway() *mockTodoGateway {
	return &mockTodoGateway{todos: make(map[string]*domain.Todo)}
}

func (m *mockTodoGateway) FindByID(id string) (*domain.Todo, error)           { return m.todos[id], nil }
func (m *mockTodoGateway) FindByUserID(userID string) ([]*domain.Todo, error) { return nil, nil }
func (m *mockTodoGateway) FindByCategory(userID, categoryID string) ([]*domain.Todo, error) {
	return nil, nil
}
func (m *mockTodoGateway) FindPending(userID string) ([]*domain.Todo, error)   { return nil, nil }
func (m *mockTodoGateway) FindCompleted(userID string) ([]*domain.Todo, error) { return nil, nil }
func (m *mockTodoGateway) Search(userID, query string) ([]*domain.Todo, error) { return nil, nil }
func (m *mockTodoGateway) Save(todo *domain.Todo) error {
	m.todos[todo.ID] = todo
	return nil
}
func (m *mockTodoGateway) Update(todo *domain.Todo) error {
	m.todos[todo.ID] = todo
	m.updated = append(m.updated, todo)
	return nil
}
func (m *mockTodoGateway) Delete(id string) error {
	m.deleted = append(m.deleted, id)
	delete(m.todos, id)
	return nil
}

// --- テスト ---

func TestTodoCreate_SetsIDAndTimestamps(t *testing.T) {
	gw := newMockTodoGateway()
	uc := usecase.NewTodoInteractor(gw)

	todo, err := uc.Create(dto.CreateTodoInput{
		UserID:  "user1",
		Title:   "タスク",
		Content: "内容",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if todo.ID == "" {
		t.Error("ID should be set")
	}
	if todo.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}
}

func TestTodoComplete_SetsCompletedAt(t *testing.T) {
	gw := newMockTodoGateway()
	gw.todos["todo1"] = &domain.Todo{ID: "todo1", CompletedAt: nil}
	uc := usecase.NewTodoInteractor(gw)

	err := uc.Complete("todo1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	updated := gw.updated[len(gw.updated)-1]
	if updated.CompletedAt == nil {
		t.Error("CompletedAt should be set after Complete")
	}
}

func TestTodoIncomplete_ClearsCompletedAt(t *testing.T) {
	gw := newMockTodoGateway()
	gw.todos["todo1"] = &domain.Todo{ID: "todo1"}
	uc := usecase.NewTodoInteractor(gw)
	_ = uc.Complete("todo1")

	err := uc.Incomplete("todo1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	updated := gw.updated[len(gw.updated)-1]
	if updated.CompletedAt != nil {
		t.Error("CompletedAt should be nil after Incomplete")
	}
}

func TestTodoTogglePin_FlipsIsPinned(t *testing.T) {
	gw := newMockTodoGateway()
	gw.todos["todo1"] = &domain.Todo{ID: "todo1", IsPinned: false}
	uc := usecase.NewTodoInteractor(gw)

	todo, err := uc.TogglePin("todo1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !todo.IsPinned {
		t.Error("IsPinned should be true after toggle")
	}
}

func TestTodoDelete_NotFound(t *testing.T) {
	gw := newMockTodoGateway()
	uc := usecase.NewTodoInteractor(gw)

	err := uc.Delete("not-exist")

	if err == nil {
		t.Error("expected error for not found")
	}
}
