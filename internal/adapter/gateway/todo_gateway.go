// internal/adapter/gateway/todo_gateway.go
package gateway

import (
	"errors"

	"gorm.io/gorm"
	"tech-memo/internal/adapter/gateway/model"
	appgateway "tech-memo/internal/application/gateway"
	"tech-memo/internal/domain"
)

var _ appgateway.TodoGateway = (*GORMTodoGateway)(nil)

type GORMTodoGateway struct {
	db *gorm.DB
}

func NewGORMTodoGateway(db *gorm.DB) *GORMTodoGateway {
	return &GORMTodoGateway{db: db}
}

func (g *GORMTodoGateway) FindByID(id string) (*domain.Todo, error) {
	var m model.TodoModel
	if err := g.db.First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (g *GORMTodoGateway) FindByUserID(userID string) ([]*domain.Todo, error) {
	var models []model.TodoModel
	if err := g.db.Where("user_id = ?", userID).Order("updated_at desc").Find(&models).Error; err != nil {
		return nil, err
	}
	return toTodos(models), nil
}

func (g *GORMTodoGateway) FindByCategory(userID, categoryID string) ([]*domain.Todo, error) {
	var models []model.TodoModel
	if err := g.db.Where("user_id = ? AND category_id = ?", userID, categoryID).
		Order("updated_at desc").Find(&models).Error; err != nil {
		return nil, err
	}
	return toTodos(models), nil
}

func (g *GORMTodoGateway) FindPending(userID string) ([]*domain.Todo, error) {
	var models []model.TodoModel
	if err := g.db.Where("user_id = ? AND completed_at IS NULL", userID).
		Order("updated_at desc").Find(&models).Error; err != nil {
		return nil, err
	}
	return toTodos(models), nil
}

func (g *GORMTodoGateway) FindCompleted(userID string) ([]*domain.Todo, error) {
	var models []model.TodoModel
	if err := g.db.Where("user_id = ? AND completed_at IS NOT NULL", userID).
		Order("updated_at desc").Find(&models).Error; err != nil {
		return nil, err
	}
	return toTodos(models), nil
}

func (g *GORMTodoGateway) Search(userID, query string) ([]*domain.Todo, error) {
	var models []model.TodoModel
	like := "%" + query + "%"
	if err := g.db.Where("user_id = ? AND (title LIKE ? OR content LIKE ?)", userID, like, like).
		Order("updated_at desc").Find(&models).Error; err != nil {
		return nil, err
	}
	return toTodos(models), nil
}

func (g *GORMTodoGateway) Save(todo *domain.Todo) error {
	m := model.FromTodo(todo)
	return g.db.Create(&m).Error
}

func (g *GORMTodoGateway) Update(todo *domain.Todo) error {
	m := model.FromTodo(todo)
	return g.db.Save(&m).Error
}

func (g *GORMTodoGateway) Delete(id string) error {
	return g.db.Where("id = ?", id).Delete(&model.TodoModel{}).Error
}

func toTodos(models []model.TodoModel) []*domain.Todo {
	todos := make([]*domain.Todo, len(models))
	for i, m := range models {
		todos[i] = m.ToDomain()
	}
	return todos
}
