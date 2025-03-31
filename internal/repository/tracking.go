package repository

import (
	"sweng-task/internal/model"
)

type TrackingRepository interface {
	Store(event *model.TrackingEventEntity) error
	FindAll() ([]*model.TrackingEventEntity, error)
	CountEvents(lineItemID string, placement string) (model.EventCounts, error)
}
