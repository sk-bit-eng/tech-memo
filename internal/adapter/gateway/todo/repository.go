package todo

import (
	"errors"

	"gorm.io/gorm"
	"tech-memo/internal/adapter/gateway/model"
	todogtw "tech-memo/internal/application/gateway/todo"
	"tech-memo/internal/domain"
)

var _ todogtw.Repository = (*Repository)(nil)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByID(id string) (*domain.Todo, error) {
	var m model.TodoModel
	if err := r.db.First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *Repository) FindByUserID(userID string) ([]*domain.Todo, error) {
	var models []model.TodoModel
	if err := r.db.Where("user_id = ?", userID).Order("updated_at desc").Find(&models).Error; err != nil {
		return nil, err
	}
	return toTodos(models), nil
}

func (r *Repository) FindByCategory(userID, categoryID string) ([]*domain.Todo, error) {
	var models []model.TodoModel
	if err := r.db.Where("user_id = ? AND category_id = ?", userID, categoryID).
		Order("updated_at desc").Find(&models).Error; err != nil {
		return nil, err
	}
	return toTodos(models), nil
}

func (r *Repository) FindPending(userID string) ([]*domain.Todo, error) {
	var models []model.TodoModel
	if err := r.db.Where("user_id = ? AND completed_at IS NULL", userID).
		Order("updated_at desc").Find(&models).Error; err != nil {
		return nil, err
	}
	return toTodos(models), nil
}

func (r *Repository) FindCompleted(userID string) ([]*domain.Todo, error) {
	var models []model.TodoModel
	if err := r.db.Where("user_id = ? AND completed_at IS NOT NULL", userID).
		Order("updated_at desc").Find(&models).Error; err != nil {
		return nil, err
	}
	return toTodos(models), nil
}

func (r *Repository) Search(userID, query string) ([]*domain.Todo, error) {
	var models []model.TodoModel
	like := "%" + query + "%"
	if err := r.db.Where("user_id = ? AND (title LIKE ? OR content LIKE ?)", userID, like, like).
		Order("updated_at desc").Find(&models).Error; err != nil {
		return nil, err
	}
	return toTodos(models), nil
}

func (r *Repository) Save(todo *domain.Todo) error {
	m := model.FromTodo(todo)
	return r.db.Create(&m).Error
}

func (r *Repository) Update(todo *domain.Todo) error {
	m := model.FromTodo(todo)
	return r.db.Save(&m).Error
}

func (r *Repository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.TodoModel{}).Error
}

func toTodos(models []model.TodoModel) []*domain.Todo {
	todos := make([]*domain.Todo, len(models))
	for i, m := range models {
		todos[i] = m.ToDomain()
	}
	return todos
}
