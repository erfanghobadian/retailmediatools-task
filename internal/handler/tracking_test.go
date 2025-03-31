package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"sweng-task/internal/model"
	"sweng-task/internal/repository"
	"sweng-task/internal/repository/mocks"
	"sweng-task/internal/service"
	"sweng-task/internal/testutil"
)

func setupTrackingTest(t *testing.T) (*fiber.App, repository.LineItemRepository, repository.TrackingRepository) {
	app := testutil.SetupTestApp(t)

	trackingRepo := mocks.NewInMemoryTrackingRepository()
	lineItemRepo := mocks.NewInMemoryLineItemRepository()
	logger := testutil.GetTestLogger()

	lineItemService := service.NewLineItemService(lineItemRepo, logger)
	trackingService := service.NewTrackingService(trackingRepo, lineItemService, logger)
	handler := NewTrackingHandler(trackingService, logger)

	app.Post("/api/v1/tracking", handler.TrackEvent)

	return app, lineItemRepo, trackingRepo
}

func TestTrackingHandler_TrackEvent(t *testing.T) {
	app, lineItemRepo, _ := setupTrackingTest(t)

	lineItem := testutil.CreateTestLineItemEntity()
	err := lineItemRepo.Create(lineItem)
	assert.NoError(t, err)

	event := testutil.CreateTestTrackingEvent(lineItem.ID)

	body, _ := json.Marshal(event)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tracking", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusAccepted, resp.StatusCode)
}

func TestTrackingHandler_TrackEvent_InvalidInput(t *testing.T) {
	app, lineItemRepo, _ := setupTrackingTest(t)
	lineItem := testutil.CreateTestLineItemEntity()
	err := lineItemRepo.Create(lineItem)
	event := map[string]interface{}{
		"event_type":   "invalid_type",
		"line_item_id": lineItem.ID,
		"timestamp":    time.Now(),
		"placement":    "homepage_top",
		"user_id":      "user_123",
		"metadata":     map[string]string{"device": "mobile"},
	}

	body, _ := json.Marshal(event)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tracking", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestTrackingHandler_TrackEvent_LineItemNotFound(t *testing.T) {
	app, _, _ := setupTrackingTest(t)

	event := testutil.CreateTestTrackingEvent("li_missing")

	body, _ := json.Marshal(event)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tracking", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, "Line item not found", result["message"])
}

func TestTrackingHandler_TrackEvent_DifferentEventTypes(t *testing.T) {
	app, lineItemRepo, _ := setupTrackingTest(t)

	lineItem := testutil.CreateTestLineItemEntity()
	err := lineItemRepo.Create(lineItem)
	assert.NoError(t, err)

	eventTypes := []model.TrackingEventType{
		model.TrackingEventTypeImpression,
		model.TrackingEventTypeClick,
		model.TrackingEventTypeConversion,
	}

	for _, eventType := range eventTypes {
		t.Run(string(eventType), func(t *testing.T) {
			event := testutil.CreateTestTrackingEvent(lineItem.ID)
			event.EventType = eventType

			body, _ := json.Marshal(event)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/tracking", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusAccepted, resp.StatusCode)
		})
	}
}
