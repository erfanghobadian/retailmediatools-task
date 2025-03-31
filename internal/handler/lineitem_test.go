package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"sweng-task/internal/model"
	"sweng-task/internal/repository/mocks"
	"sweng-task/internal/service"
	"sweng-task/internal/testutil"
)

func setupLineItemTest(t *testing.T) (*fiber.App, *mocks.LineItemRepository) {
	app := testutil.SetupTestApp(t)
	mockRepo := mocks.NewInMemoryLineItemRepository()
	logger := testutil.GetTestLogger()

	svc := service.NewLineItemService(
		mockRepo,
		logger,
	)

	handler := NewLineItemHandler(svc, logger)

	app.Post("/api/v1/lineitems", handler.Create)
	app.Get("/api/v1/lineitems/:id", handler.GetByID)
	app.Get("/api/v1/lineitems", handler.GetAll)

	return app, mockRepo
}

func TestLineItemHandler_Create(t *testing.T) {
	app, mockRepo := setupLineItemTest(t)

	input := testutil.CreateTestLineItemCreate()
	expected := testutil.CreateTestLineItemEntity()
	_ = mockRepo.Create(expected)

	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/lineitems", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result model.LineItem
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, input.Name, result.Name)
	assert.Equal(t, input.AdvertiserID, result.AdvertiserID)
	assert.Equal(t, input.Bid, result.Bid)
	assert.Equal(t, input.Budget, result.Budget)
	assert.Equal(t, input.Placement, result.Placement)
	assert.Equal(t, input.Categories, result.Categories)
	assert.Equal(t, input.Keywords, result.Keywords)
}

func TestLineItemHandler_Create_InvalidInput(t *testing.T) {
	app, _ := setupLineItemTest(t)

	input := model.LineItemCreate{
		Name:   "Test Ad",
		Bid:    2.5,
		Budget: 1000.0,
	}

	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/lineitems", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestLineItemHandler_GetByID(t *testing.T) {
	app, mockRepo := setupLineItemTest(t)
	expected := testutil.CreateTestLineItemEntity()
	_ = mockRepo.Create(expected)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lineitems/"+expected.ID, nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result model.LineItem
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)
	assert.Equal(t, expected.Name, result.Name)
}

func TestLineItemHandler_GetByID_NotFound(t *testing.T) {
	app, _ := setupLineItemTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lineitems/nonexistent_id", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, "Line item not found", result["message"])
}

func TestLineItemHandler_GetAll(t *testing.T) {
	app, mockRepo := setupLineItemTest(t)

	li1 := testutil.CreateTestLineItemEntity()
	li2 := testutil.CreateTestLineItemEntity()
	_ = mockRepo.Create(li1)
	_ = mockRepo.Create(li2)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lineitems?advertiser_id="+li1.AdvertiserID+"&placement="+li1.Placement, nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result []model.LineItem
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}
