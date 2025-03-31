package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sweng-task/internal/repository/mocks"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"sweng-task/internal/model"
	"sweng-task/internal/service"
	"sweng-task/internal/testutil"
)

func setupAdHandlerTest(t *testing.T) (*fiber.App, *service.AdService) {
	app := testutil.SetupTestApp(t)

	mockLineItemRepo := mocks.NewInMemoryLineItemRepository()
	mockTrackingRepo := mocks.NewInMemoryTrackingRepository()
	logger := testutil.GetTestLogger()

	lineItemService := service.NewLineItemService(mockLineItemRepo, logger)
	trackingService := service.NewTrackingService(mockTrackingRepo, lineItemService, logger)
	adService := service.NewAdService(lineItemService, trackingService, logger)

	h := NewAdSelectionHandler(adService, logger)
	app.Get("/api/v1/ads", h.GetWinningAds)

	return app, adService
}

func TestAdSelectionHandler_GetWinningAds_Success(t *testing.T) {
	app, _ := setupAdHandlerTest(t)
	item := testutil.CreateTestLineItemEntity()
	_ = mocks.NewInMemoryLineItemRepository().Create(item)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/ads?placement="+item.Placement+"&limit=1", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var ads []model.Ad
	err = json.NewDecoder(resp.Body).Decode(&ads)
	assert.NoError(t, err)
	assert.LessOrEqual(t, len(ads), 1)
}

func TestAdSelectionHandler_GetWinningAds_MissingPlacement(t *testing.T) {
	app, _ := setupAdHandlerTest(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/ads", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var body map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Equal(t, float64(400), body["code"])
	assert.Contains(t, body["message"], "Invalid request")
}
