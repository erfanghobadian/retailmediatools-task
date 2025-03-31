package testutil

import (
	"testing"
	"time"

	"sweng-task/internal/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func SetupTestApp(t *testing.T) *fiber.App {
	t.Helper()
	app := fiber.New()
	return app
}

func CreateTestLineItem() *model.LineItem {
	return &model.LineItem{
		ID:           "li_" + uuid.New().String(),
		Name:         "Test Ad",
		AdvertiserID: "adv_123",
		Bid:          2.5,
		Budget:       1000.0,
		Placement:    "homepage_top",
		Categories:   []string{"electronics"},
		Keywords:     []string{"test"},
		Status:       model.LineItemStatusActive,
	}
}

func CreateTestLineItemEntity() *model.LineItemEntity {
	return &model.LineItemEntity{
		ID:           "li_" + uuid.New().String(),
		Name:         "Test Ad",
		AdvertiserID: "adv_123",
		Bid:          2.5,
		Budget:       1000.0,
		Placement:    "homepage_top",
		Categories:   []string{"electronics"},
		Keywords:     []string{"test"},
		Status:       model.LineItemStatusActive,
	}
}

func CreateTestLineItemCreate() model.LineItemCreate {
	return model.LineItemCreate{
		Name:         "Test Ad",
		AdvertiserID: "adv_123",
		Bid:          2.5,
		Budget:       1000.0,
		Placement:    "homepage_top",
		Categories:   []string{"electronics"},
		Keywords:     []string{"test"},
	}
}

func CreateTestTrackingEvent(lineItemID string) model.TrackingEvent {
	return model.TrackingEvent{
		EventType:  model.TrackingEventTypeImpression,
		LineItemID: lineItemID,
		Timestamp:  time.Now(),
		Placement:  "homepage_top",
		UserID:     "user_123",
		Metadata:   map[string]string{"device": "mobile"},
	}
}

func CreateTestTrackingEventEntity(lineItemID string) *model.TrackingEventEntity {
	return &model.TrackingEventEntity{
		EventType:  model.TrackingEventTypeImpression,
		LineItemID: lineItemID,
		Timestamp:  time.Now(),
		Placement:  "homepage_top",
		UserID:     "user_123",
		Metadata:   map[string]string{"device": "mobile"},
	}
}

// GetTestLogger returns a no-op logger for testing
func GetTestLogger() *zap.SugaredLogger {
	return zap.NewNop().Sugar()
}
