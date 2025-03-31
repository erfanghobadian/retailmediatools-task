package service

import (
	//"github.com/alicebob/miniredis/v2"
	"go.uber.org/zap"
	"sweng-task/internal/model"
	"sweng-task/internal/repository"
)

type TrackingService struct {
	repo            repository.TrackingRepository
	lineItemService *LineItemService
	logger          *zap.SugaredLogger
}

func NewTrackingService(repo repository.TrackingRepository, lineItemService *LineItemService, logger *zap.SugaredLogger) *TrackingService {
	return &TrackingService{repo: repo, lineItemService: lineItemService, logger: logger}
}

func (s *TrackingService) Track(event model.TrackingEvent) error {
	s.logger.Infow("Tracking event", "event_type", event.EventType, "line_item_id", event.LineItemID)

	// 1. Check if LineItem exists
	lineItem, err := s.lineItemService.GetByID(event.LineItemID)
	if err != nil {
		return ErrLineItemNotFound
	}

	// 2. Increase daily spending if relevant event
	switch event.EventType {
	case model.TrackingEventTypeImpression,
		model.TrackingEventTypeClick,
		model.TrackingEventTypeConversion:
		costPerEvent := lineItem.Bid / 1000

		if err := s.lineItemService.IncreaseDailySpending(lineItem.ID, costPerEvent); err != nil {
			s.logger.Errorw("Failed to increase daily spending", "line_item_id", lineItem.ID, "error", err)
			return err
		}
	}

	// 3. Store the tracking event
	eventEntity := model.ToEntityTrackingEvent(event)
	if err := s.repo.Store(&eventEntity); err != nil {
		s.logger.Errorw("Failed to store tracking event", "error", err)
		return err
	}

	return nil
}

func (s *TrackingService) GetEventCounts(lineItemID string, placement string) (model.EventCounts, error) {
	return s.repo.CountEvents(lineItemID, placement)
}
