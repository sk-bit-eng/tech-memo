package memo

import (
	"errors"

	model "tech-memo/internal/adapter/gateway/model/memo"
	memogtw "tech-memo/internal/application/gateway/memo"
	"tech-memo/internal/domain"

	"gorm.io/gorm"
)

var _ memogtw.Repository = (*Repository)(nil)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByID(id string) (*domain.Memo, error) {
	var m model.MemoModel
	if err := r.db.First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *Repository) FindByUserID(userID string) ([]*domain.Memo, error) {
	var models []model.MemoModel
	if err := r.db.Where("user_id = ?", userID).Order("updated_at desc").Find(&models).Error; err != nil {
		return nil, err
	}
	return toMemos(models), nil
}

func (r *Repository) FindByCategory(userID, categoryID string) ([]*domain.Memo, error) {
	var models []model.MemoModel
	if err := r.db.Where("user_id = ? AND category_id = ?", userID, categoryID).
		Order("updated_at desc").Find(&models).Error; err != nil {
		return nil, err
	}
	return toMemos(models), nil
}

func (r *Repository) Search(userID, query string) ([]*domain.Memo, error) {
	var models []model.MemoModel
	like := "%" + query + "%"
	if err := r.db.Where("user_id = ? AND (title LIKE ? OR content LIKE ?)", userID, like, like).
		Order("updated_at desc").Find(&models).Error; err != nil {
		return nil, err
	}
	return toMemos(models), nil
}

func (r *Repository) Save(memo *domain.Memo) error {
	m := model.FromMemo(memo)
	return r.db.Create(&m).Error
}

func (r *Repository) Update(memo *domain.Memo) error {
	m := model.FromMemo(memo)
	return r.db.Save(&m).Error
}

func (r *Repository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.MemoModel{}).Error
}

func toMemos(models []model.MemoModel) []*domain.Memo {
	memos := make([]*domain.Memo, len(models))
	for i, m := range models {
		memos[i] = m.ToDomain()
	}
	return memos
}
