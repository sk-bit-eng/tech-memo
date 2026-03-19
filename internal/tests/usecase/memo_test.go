// internal/tests/usecase/memo_test.go
package usecase_test

import (
	"testing"

	"tech-memo/internal/application/dto"
	"tech-memo/internal/application/usecase"
	"tech-memo/internal/domain"
)

// --- モックゲートウェイ ---

type mockMemoGateway struct {
	memos   map[string]*domain.Memo
	saved   []*domain.Memo
	updated []*domain.Memo
	deleted []string
}

func newMockMemoGateway() *mockMemoGateway {
	return &mockMemoGateway{memos: make(map[string]*domain.Memo)}
}

func (m *mockMemoGateway) FindByID(id string) (*domain.Memo, error) {
	return m.memos[id], nil
}
func (m *mockMemoGateway) FindByUserID(userID string) ([]*domain.Memo, error) { return nil, nil }
func (m *mockMemoGateway) FindByCategory(userID, categoryID string) ([]*domain.Memo, error) {
	return nil, nil
}
func (m *mockMemoGateway) Search(userID, query string) ([]*domain.Memo, error) { return nil, nil }
func (m *mockMemoGateway) Save(memo *domain.Memo) error {
	m.memos[memo.ID] = memo
	m.saved = append(m.saved, memo)
	return nil
}
func (m *mockMemoGateway) Update(memo *domain.Memo) error {
	m.memos[memo.ID] = memo
	m.updated = append(m.updated, memo)
	return nil
}
func (m *mockMemoGateway) Delete(id string) error {
	m.deleted = append(m.deleted, id)
	delete(m.memos, id)
	return nil
}

// --- テスト ---

func TestMemoCreate_SetsIDAndTimestamps(t *testing.T) {
	gw := newMockMemoGateway()
	uc := usecase.NewMemoInteractor(gw)

	memo, err := uc.Create(dto.CreateMemoInput{
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

func TestMemoGetByID_NotFound(t *testing.T) {
	gw := newMockMemoGateway()
	uc := usecase.NewMemoInteractor(gw)

	_, err := uc.GetByID("not-exist")

	if err == nil {
		t.Error("expected error for not found")
	}
}

func TestMemoDelete_NotFound(t *testing.T) {
	gw := newMockMemoGateway()
	uc := usecase.NewMemoInteractor(gw)

	err := uc.Delete("not-exist")

	if err == nil {
		t.Error("expected error for not found")
	}
}

func TestMemoDelete_CallsGateway(t *testing.T) {
	gw := newMockMemoGateway()
	gw.memos["memo1"] = &domain.Memo{ID: "memo1"}
	uc := usecase.NewMemoInteractor(gw)

	err := uc.Delete("memo1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(gw.deleted) != 1 || gw.deleted[0] != "memo1" {
		t.Error("expected gateway.Delete to be called with memo1")
	}
}

func TestMemoTogglePin_FlipsIsPinned(t *testing.T) {
	gw := newMockMemoGateway()
	gw.memos["memo1"] = &domain.Memo{ID: "memo1", IsPinned: false}
	uc := usecase.NewMemoInteractor(gw)

	memo, err := uc.TogglePin("memo1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !memo.IsPinned {
		t.Error("IsPinned should be true after toggle")
	}
}

func TestMemoUpdate_NotFound(t *testing.T) {
	gw := newMockMemoGateway()
	uc := usecase.NewMemoInteractor(gw)

	_, err := uc.Update(dto.UpdateMemoInput{ID: "not-exist"})

	if err == nil {
		t.Error("expected error for not found")
	}
}
