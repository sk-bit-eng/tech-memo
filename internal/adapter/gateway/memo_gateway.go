// internal/adapter/gateway/memo_gateway.go
package gateway

import (
	"errors"

	"gorm.io/gorm"
	"tech-memo/internal/adapter/gateway/model"
	memogtw "tech-memo/internal/application/gateway/memo"
	"tech-memo/internal/domain"
)

var _ memogtw.Repository = (*GORMMemoGateway)(nil)

type GORMMemoGateway struct {
	db *gorm.DB
}

func NewGORMMemoGateway(db *gorm.DB) *GORMMemoGateway {
	return &GORMMemoGateway{db: db}
}

func (g *GORMMemoGateway) FindByID(id string) (*domain.Memo, error) {
	var m model.MemoModel
	if err := g.db.First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (g *GORMMemoGateway) FindByUserID(userID string) ([]*domain.Memo, error) {
	var models []model.MemoModel
	if err := g.db.Where("user_id = ?", userID).Order("updated_at desc").Find(&models).Error; err != nil {
		return nil, err
	}
	return toMemos(models), nil
}

func (g *GORMMemoGateway) FindByCategory(userID, categoryID string) ([]*domain.Memo, error) {
	var models []model.MemoModel
	if err := g.db.Where("user_id = ? AND category_id = ?", userID, categoryID).
		Order("updated_at desc").Find(&models).Error; err != nil {
		return nil, err
	}
	return toMemos(models), nil
}

func (g *GORMMemoGateway) Search(userID, query string) ([]*domain.Memo, error) {
	var models []model.MemoModel
	like := "%" + query + "%"
	if err := g.db.Where("user_id = ? AND (title LIKE ? OR content LIKE ?)", userID, like, like).
		Order("updated_at desc").Find(&models).Error; err != nil {
		return nil, err
	}
	return toMemos(models), nil
}

func (g *GORMMemoGateway) Save(memo *domain.Memo) error {
	m := model.FromMemo(memo)
	return g.db.Create(&m).Error
}

func (g *GORMMemoGateway) Update(memo *domain.Memo) error {
	m := model.FromMemo(memo)
	return g.db.Save(&m).Error
}

func (g *GORMMemoGateway) Delete(id string) error {
	return g.db.Where("id = ?", id).Delete(&model.MemoModel{}).Error
}

func toMemos(models []model.MemoModel) []*domain.Memo {
	memos := make([]*domain.Memo, len(models))
	for i, m := range models {
		memos[i] = m.ToDomain()
	}
	return memos
}
