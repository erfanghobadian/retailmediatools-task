package postgres

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sweng-task/internal/model"
)

type LineItemPostgresRepository struct {
	db  *gorm.DB
	log *zap.SugaredLogger
}

func NewLineItemPostgresRepository(db *gorm.DB, log *zap.SugaredLogger) *LineItemPostgresRepository {
	return &LineItemPostgresRepository{db: db, log: log}
}

func (r *LineItemPostgresRepository) Create(item *model.LineItemEntity) error {
	result := r.db.Create(item)
	return result.Error
}

func (r *LineItemPostgresRepository) GetByID(id string) (*model.LineItemEntity, error) {
	var item model.LineItemEntity
	result := r.db.First(&item, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (r *LineItemPostgresRepository) GetAll(advertiserID, placement string) ([]*model.LineItemEntity, error) {
	var items []*model.LineItemEntity
	query := r.db.Model(&model.LineItemEntity{})

	if advertiserID != "" {
		query = query.Where("advertiser_id = ?", advertiserID)
	}
	if placement != "" {
		query = query.Where("placement = ?", placement)
	}

	err := query.Find(&items).Error
	return items, err
}

func (r *LineItemPostgresRepository) FindMatchingLineItems(placement, category, keyword string) ([]*model.LineItemEntity, error) {
	var items []*model.LineItemEntity

	query := r.db.Where("placement = ? AND status = ? AND daily_budget_remaining >= bid", placement, "active")

	if category != "" {
		query = query.Where("? = ANY(categories)", category)
	}
	if keyword != "" {
		query = query.Where("? = ANY(keywords)", keyword)
	}

	err := query.Find(&items).Error
	return items, err
}

func (r *LineItemPostgresRepository) ResetDailySpending() error {
	result := r.db.Model(&model.LineItemEntity{}).
		Where("daily_spending > 0").
		Update("daily_spending", 0)

	if result.Error != nil {
		r.log.Errorw("Failed to reset daily budgets", "error", result.Error)
		return result.Error
	}

	r.log.Infow("Daily budgets reset", "affected_rows", result.RowsAffected)
	return nil
}

func (r *LineItemPostgresRepository) IncreaseDailySpending(lineItemID string, amount float64) error {

	result := r.db.Model(&model.LineItemEntity{}).
		Where("id = ?", lineItemID).
		Update("daily_spending", gorm.Expr("daily_spending + ?", amount))

	if result.Error != nil {
		return result.Error
	}
	return nil
}
