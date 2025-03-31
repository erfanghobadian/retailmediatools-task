package postgres

import (
	"go.uber.org/zap"
	"gorm.io/gorm"

	"sweng-task/internal/model"
)

type TrackingPostgresRepository struct {
	db  *gorm.DB
	log *zap.SugaredLogger
}

func NewTrackingPostgresRepository(db *gorm.DB, log *zap.SugaredLogger) *TrackingPostgresRepository {
	return &TrackingPostgresRepository{db: db, log: log}
}

func (r *TrackingPostgresRepository) Store(event *model.TrackingEventEntity) error {

	if err := r.db.Create(&event).Error; err != nil {
		r.log.Errorw("Failed to insert tracking event", "error", err)
		return err
	}

	r.log.Infow("Tracking event stored",
		"event_type", event.EventType,
		"line_item_id", event.LineItemID,
		"placement", event.Placement,
		"user_id", event.UserID,
	)
	return nil
}

func (r *TrackingPostgresRepository) FindAll() ([]*model.TrackingEventEntity, error) {
	var events []*model.TrackingEventEntity

	if err := r.db.Find(&events).Error; err != nil {
		r.log.Errorw("Failed to fetch tracking events", "error", err)
		return nil, err
	}

	r.log.Infow("Fetched tracking events", "count", len(events))
	return events, nil
}

func (r *TrackingPostgresRepository) CountEvents(lineItemID string, placement string) (model.EventCounts, error) {
	var groupedCounts []struct {
		EventType string
		Count     int
	}

	query := r.db.Model(&model.TrackingEventEntity{}).
		Select("event_type, COUNT(*) as count")

	if lineItemID != "" {
		query = query.Where("line_item_id = ?", lineItemID)
	}
	if placement != "" {
		query = query.Where("placement = ?", placement)
	}

	if err := query.Group("event_type").Scan(&groupedCounts).Error; err != nil {
		return model.EventCounts{}, err
	}

	// Map results
	var counts model.EventCounts
	for _, row := range groupedCounts {
		switch row.EventType {
		case string(model.TrackingEventTypeImpression):
			counts.Impressions = row.Count
		case string(model.TrackingEventTypeClick):
			counts.Clicks = row.Count
		case string(model.TrackingEventTypeConversion):
			counts.Conversions = row.Count
		}
	}

	return counts, nil
}
