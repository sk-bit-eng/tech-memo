package todo

import (
	"fmt"
	"time"

	"tech-memo/internal/application/dto"
	todogtw "tech-memo/internal/application/gateway/todo"
	"tech-memo/internal/domain"

	"github.com/google/uuid"
)

type interactor struct {
	gw todogtw.Repository
}

func NewInteractor(gw todogtw.Repository) UseCase {
	return &interactor{gw: gw}
}

func (uc *interactor) GetByID(id string) (*domain.Todo, error) {
	todo, err := uc.gw.FindByID(id)
	if err != nil {
		return nil, err
	}
	if todo == nil {
		return nil, fmt.Errorf("todo not found: %s", id)
	}
	return todo, nil
}

func (uc *interactor) ListByUser(userID string) ([]*domain.Todo, error) {
	return uc.gw.FindByUserID(userID)
}

func (uc *interactor) ListByCategory(userID, categoryID string) ([]*domain.Todo, error) {
	return uc.gw.FindByCategory(userID, categoryID)
}

func (uc *interactor) ListPending(userID string) ([]*domain.Todo, error) {
	return uc.gw.FindPending(userID)
}

func (uc *interactor) ListCompleted(userID string) ([]*domain.Todo, error) {
	return uc.gw.FindCompleted(userID)
}

func (uc *interactor) Search(userID, query string) ([]*domain.Todo, error) {
	return uc.gw.Search(userID, query)
}

func (uc *interactor) Create(input dto.CreateTodoInput) (*domain.Todo, error) {
	now := time.Now()
	todo := &domain.Todo{
		ID:         uuid.New().String(),
		UserID:     input.UserID,
		Title:      input.Title,
		Content:    input.Content,
		CategoryID: input.CategoryID,
		Parameters: input.Parameters,
		IsPinned:   input.IsPinned,
		DueAt:      input.DueAt,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := uc.gw.Save(todo); err != nil {
		return nil, err
	}
	return todo, nil
}

func (uc *interactor) Update(input dto.UpdateTodoInput) (*domain.Todo, error) {
	todo, err := uc.gw.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if todo == nil {
		return nil, fmt.Errorf("todo not found: %s", input.ID)
	}
	todo.Title = input.Title
	todo.Content = input.Content
	todo.CategoryID = input.CategoryID
	todo.Parameters = input.Parameters
	todo.IsPinned = input.IsPinned
	todo.DueAt = input.DueAt
	todo.UpdatedAt = time.Now()
	if err := uc.gw.Update(todo); err != nil {
		return nil, err
	}
	return todo, nil
}

func (uc *interactor) Complete(id string) error {
	todo, err := uc.gw.FindByID(id)
	if err != nil {
		return err
	}
	if todo == nil {
		return fmt.Errorf("todo not found: %s", id)
	}
	now := time.Now()
	todo.CompletedAt = &now
	todo.UpdatedAt = now
	return uc.gw.Update(todo)
}

func (uc *interactor) Incomplete(id string) error {
	todo, err := uc.gw.FindByID(id)
	if err != nil {
		return err
	}
	if todo == nil {
		return fmt.Errorf("todo not found: %s", id)
	}
	todo.CompletedAt = nil
	todo.UpdatedAt = time.Now()
	return uc.gw.Update(todo)
}

func (uc *interactor) Delete(id string) error {
	todo, err := uc.gw.FindByID(id)
	if err != nil {
		return err
	}
	if todo == nil {
		return fmt.Errorf("todo not found: %s", id)
	}
	return uc.gw.Delete(id)
}

func (uc *interactor) TogglePin(id string) (*domain.Todo, error) {
	todo, err := uc.gw.FindByID(id)
	if err != nil {
		return nil, err
	}
	if todo == nil {
		return nil, fmt.Errorf("todo not found: %s", id)
	}
	todo.IsPinned = !todo.IsPinned
	todo.UpdatedAt = time.Now()
	if err := uc.gw.Update(todo); err != nil {
		return nil, err
	}
	return todo, nil
}
