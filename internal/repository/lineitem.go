package repository

import (
	"sweng-task/internal/model"
)

type LineItemRepository interface {
	Create(item *model.LineItemEntity) error
	GetByID(id string) (*model.LineItemEntity, error)
	GetAll(advertiserID, placement string) ([]*model.LineItemEntity, error)
	FindMatchingLineItems(placement, category, keyword string) ([]*model.LineItemEntity, error)
	ResetDailySpending() (err error)
	IncreaseDailySpending(lineItemID string, amount float64) error
}
