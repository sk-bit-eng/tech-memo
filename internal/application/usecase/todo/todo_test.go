package todo

import (
	"testing"

	"tech-memo/internal/application/dto"
	"tech-memo/internal/domain"
)

// --- モックゲートウェイ ---

type mockRepository struct {
	todos   map[string]*domain.Todo
	updated []*domain.Todo
	deleted []string
}

func newMockRepository() *mockRepository {
	return &mockRepository{todos: make(map[string]*domain.Todo)}
}

func (m *mockRepository) FindByID(id string) (*domain.Todo, error)           { return m.todos[id], nil }
func (m *mockRepository) FindByUserID(userID string) ([]*domain.Todo, error) { return nil, nil }
func (m *mockRepository) FindByCategory(userID, categoryID string) ([]*domain.Todo, error) {
	return nil, nil
}
func (m *mockRepository) FindPending(userID string) ([]*domain.Todo, error)   { return nil, nil }
func (m *mockRepository) FindCompleted(userID string) ([]*domain.Todo, error) { return nil, nil }
func (m *mockRepository) Search(userID, query string) ([]*domain.Todo, error) { return nil, nil }
func (m *mockRepository) Save(todo *domain.Todo) error {
	m.todos[todo.ID] = todo
	return nil
}
func (m *mockRepository) Update(todo *domain.Todo) error {
	m.todos[todo.ID] = todo
	m.updated = append(m.updated, todo)
	return nil
}
func (m *mockRepository) Delete(id string) error {
	m.deleted = append(m.deleted, id)
	delete(m.todos, id)
	return nil
}

// --- テスト ---

func TestCreate_SetsIDAndTimestamps(t *testing.T) {
	gw := newMockRepository()
	uc := NewInteractor(gw)

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

func TestComplete_SetsCompletedAt(t *testing.T) {
	gw := newMockRepository()
	gw.todos["todo1"] = &domain.Todo{ID: "todo1", CompletedAt: nil}
	uc := NewInteractor(gw)

	err := uc.Complete("todo1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	updated := gw.updated[len(gw.updated)-1]
	if updated.CompletedAt == nil {
		t.Error("CompletedAt should be set after Complete")
	}
}

func TestIncomplete_ClearsCompletedAt(t *testing.T) {
	gw := newMockRepository()
	gw.todos["todo1"] = &domain.Todo{ID: "todo1"}
	uc := NewInteractor(gw)
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

func TestTogglePin_FlipsIsPinned(t *testing.T) {
	gw := newMockRepository()
	gw.todos["todo1"] = &domain.Todo{ID: "todo1", IsPinned: false}
	uc := NewInteractor(gw)

	todo, err := uc.TogglePin("todo1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !todo.IsPinned {
		t.Error("IsPinned should be true after toggle")
	}
}

func TestDelete_NotFound(t *testing.T) {
	gw := newMockRepository()
	uc := NewInteractor(gw)

	err := uc.Delete("not-exist")

	if err == nil {
		t.Error("expected error for not found")
	}
}
