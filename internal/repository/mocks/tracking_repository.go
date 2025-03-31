package mocks

import (
	"sync"

	"sweng-task/internal/model"
)

type TrackingRepository struct {
	store []model.TrackingEventEntity
	mu    sync.RWMutex
}

func NewInMemoryTrackingRepository() *TrackingRepository {
	return &TrackingRepository{
		store: make([]model.TrackingEventEntity, 0),
	}
}

func (m *TrackingRepository) Store(event *model.TrackingEventEntity) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.store = append(m.store, *event)
	return nil
}

func (m *TrackingRepository) CountEvents(lineItemID string, placement string) (model.EventCounts, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var impressions, clicks, conversions int

	for _, e := range m.store {
		if (lineItemID == "" || e.LineItemID == lineItemID) &&
			(placement == "" || e.Placement == placement) {

			switch e.EventType {
			case model.TrackingEventTypeImpression:
				impressions++
			case model.TrackingEventTypeClick:
				clicks++
			case model.TrackingEventTypeConversion:
				conversions++
			}
		}
	}

	return model.EventCounts{
		Impressions: impressions,
		Clicks:      clicks,
		Conversions: conversions,
	}, nil
}

func (m *TrackingRepository) FindAll() ([]*model.TrackingEventEntity, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []*model.TrackingEventEntity
	for _, e := range m.store {
		results = append(results, &model.TrackingEventEntity{
			EventType:  e.EventType,
			LineItemID: e.LineItemID,
			Timestamp:  e.Timestamp,
			Placement:  e.Placement,
			UserID:     e.UserID,
			Metadata:   e.Metadata,
		})
	}

	return results, nil
}
