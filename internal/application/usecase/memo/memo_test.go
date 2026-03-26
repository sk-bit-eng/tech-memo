package memo

import (
	"testing"

	memodto "tech-memo/internal/application/dto/memo"
	"tech-memo/internal/domain"
)

// --- モックゲートウェイ ---

type mockRepository struct {
	memos   map[string]*domain.Memo
	saved   []*domain.Memo
	updated []*domain.Memo
	deleted []string
}

func newMockRepository() *mockRepository {
	return &mockRepository{memos: make(map[string]*domain.Memo)}
}

func (m *mockRepository) FindByID(id string) (*domain.Memo, error) { return m.memos[id], nil }
func (m *mockRepository) FindByUserID(userID string) ([]*domain.Memo, error) { return nil, nil }
func (m *mockRepository) FindByCategory(userID, categoryID string) ([]*domain.Memo, error) {
	return nil, nil
}
func (m *mockRepository) Search(userID, query string) ([]*domain.Memo, error) { return nil, nil }
func (m *mockRepository) Save(memo *domain.Memo) error {
	m.memos[memo.ID] = memo
	m.saved = append(m.saved, memo)
	return nil
}
func (m *mockRepository) Update(memo *domain.Memo) error {
	m.memos[memo.ID] = memo
	m.updated = append(m.updated, memo)
	return nil
}
func (m *mockRepository) Delete(id string) error {
	m.deleted = append(m.deleted, id)
	delete(m.memos, id)
	return nil
}

// --- テスト ---

func TestCreate_SetsIDAndTimestamps(t *testing.T) {
	gw := newMockRepository()
	uc := NewInteractor(gw)

	memo, err := uc.Create(memodto.CreateInput{
		UserID:  "user1",
		Title:   "タイトル",
		Content: "内容",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if memo.ID == "" {
		t.Error("ID should be set")
	}
	if memo.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}
	if memo.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should be set")
	}
	if len(gw.saved) != 1 {
		t.Errorf("expected 1 save, got %d", len(gw.saved))
	}
}

func TestGetByID_NotFound(t *testing.T) {
	gw := newMockRepository()
	uc := NewInteractor(gw)

	_, err := uc.GetByID("not-exist")

	if err == nil {
		t.Error("expected error for not found")
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

func TestDelete_CallsRepository(t *testing.T) {
	gw := newMockRepository()
	gw.memos["memo1"] = &domain.Memo{ID: "memo1"}
	uc := NewInteractor(gw)

	err := uc.Delete("memo1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(gw.deleted) != 1 || gw.deleted[0] != "memo1" {
		t.Error("expected repository.Delete to be called with memo1")
	}
}

func TestTogglePin_FlipsIsPinned(t *testing.T) {
	gw := newMockRepository()
	gw.memos["memo1"] = &domain.Memo{ID: "memo1", IsPinned: false}
	uc := NewInteractor(gw)

	memo, err := uc.TogglePin("memo1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !memo.IsPinned {
		t.Error("IsPinned should be true after toggle")
	}
}

func TestUpdate_NotFound(t *testing.T) {
	gw := newMockRepository()
	uc := NewInteractor(gw)

	_, err := uc.Update(memodto.UpdateInput{ID: "not-exist"})

	if err == nil {
		t.Error("expected error for not found")
	}
}
