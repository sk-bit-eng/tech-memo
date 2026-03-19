// internal/application/usecase/todo.go
package usecase

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	appgateway "tech-memo/internal/application/gateway"
	"tech-memo/internal/application/dto"
	"tech-memo/internal/domain"
)

// ---- インターフェース ----

type TodoUseCase interface {
	GetByID(id string) (*domain.Todo, error)
	ListByUser(userID string) ([]*domain.Todo, error)
	ListByCategory(userID, categoryID string) ([]*domain.Todo, error)
	ListPending(userID string) ([]*domain.Todo, error)
	ListCompleted(userID string) ([]*domain.Todo, error)
	Search(userID, query string) ([]*domain.Todo, error)
	Create(input dto.CreateTodoInput) (*domain.Todo, error)
	Update(input dto.UpdateTodoInput) (*domain.Todo, error)
	Complete(id string) error
	Incomplete(id string) error
	Delete(id string) error
	TogglePin(id string) (*domain.Todo, error)
}

// ---- インタラクター（実装）----

type todoInteracter struct {
	gw appgateway.TodoGateway
}

func NewTodoInteracter(gw appgateway.TodoGateway) TodoUseCase {
	return &todoInteracter{gw: gw}
}

func (uc *todoInteracter) GetByID(id string) (*domain.Todo, error) {
	todo, err := uc.gw.FindByID(id)
	if err != nil {
		return nil, err
	}
	if todo == nil {
		return nil, fmt.Errorf("todo not found: %s", id)
	}
	return todo, nil
}

func (uc *todoInteracter) ListByUser(userID string) ([]*domain.Todo, error) {
	return uc.gw.FindByUserID(userID)
}

func (uc *todoInteracter) ListByCategory(userID, categoryID string) ([]*domain.Todo, error) {
	return uc.gw.FindByCategory(userID, categoryID)
}

func (uc *todoInteracter) ListPending(userID string) ([]*domain.Todo, error) {
	return uc.gw.FindPending(userID)
}

func (uc *todoInteracter) ListCompleted(userID string) ([]*domain.Todo, error) {
	return uc.gw.FindCompleted(userID)
}

func (uc *todoInteracter) Search(userID, query string) ([]*domain.Todo, error) {
	return uc.gw.Search(userID, query)
}

func (uc *todoInteracter) Create(input dto.CreateTodoInput) (*domain.Todo, error) {
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

func (uc *todoInteracter) Update(input dto.UpdateTodoInput) (*domain.Todo, error) {
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

func (uc *todoInteracter) Complete(id string) error {
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

func (uc *todoInteracter) Incomplete(id string) error {
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

func (uc *todoInteracter) Delete(id string) error {
	todo, err := uc.gw.FindByID(id)
	if err != nil {
		return err
	}
	if todo == nil {
		return fmt.Errorf("todo not found: %s", id)
	}
	return uc.gw.Delete(id)
}

func (uc *todoInteracter) TogglePin(id string) (*domain.Todo, error) {
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
