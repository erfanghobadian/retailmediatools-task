package handler

import (
	"sweng-task/internal/model"
	"sweng-task/internal/utils"
	"sweng-task/internal/validator"

	"sweng-task/internal/service"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// LineItemHandler handles HTTP requests related to line items
type LineItemHandler struct {
	service *service.LineItemService
	log     *zap.SugaredLogger
}

// NewLineItemHandler creates a new LineItemHandler
func NewLineItemHandler(service *service.LineItemService, log *zap.SugaredLogger) *LineItemHandler {
	return &LineItemHandler{
		service: service,
		log:     log,
	}
}

// Create handles the creation of a new line item
func (h *LineItemHandler) Create(c *fiber.Ctx) error {
	var input model.LineItemCreate

	if err := c.BodyParser(&input); err != nil {
		h.log.Warnw("Invalid line item payload", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Details: err.Error(),
		})
	}

	if fieldErr, err := validator.ValidateStruct(&input); err != nil {
		h.log.Warnw("Line item validation failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Validation failed",
		})
	} else if fieldErr != nil {
		h.log.Warnw("Line item field validation error", "field", fieldErr.Field, "reason", fieldErr.Reason)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request",
			Details: fieldErr,
		})
	}

	lineItem, err := h.service.Create(input)
	if err != nil {
		h.log.Errorw("Failed to create line item", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to create line item",
			Details: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(lineItem)
}

// GetByID handles retrieving a line item by ID
func (h *LineItemHandler) GetByID(c *fiber.Ctx) error {
	var param validator.IDParam
	if err := c.ParamsParser(&param); err != nil {
		h.log.Warnw("Failed to parse path parameters", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid path parameters",
			Details: err.Error(),
		})
	}

	if fieldErr, err := validator.ValidateStruct(&param); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Validation failed",
		})
	} else if fieldErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request parameters",
			Details: fieldErr,
		})
	}

	lineItem, err := h.service.GetByID(param.ID)
	if err != nil {
		if err == service.ErrLineItemNotFound {
			return c.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Line item not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to retrieve line item",
			Details: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(lineItem)
}

// GetAll handles retrieving all line items with optional filtering
func (h *LineItemHandler) GetAll(c *fiber.Ctx) error {
	var query validator.LineItemQueryParams
	if err := c.QueryParser(&query); err != nil {
		h.log.Warnw("Failed to parse query", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid query parameters",
			Details: err.Error(),
		})
	}

	if fieldErr, err := validator.ValidateStruct(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Validation failed",
		})
	} else if fieldErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid query parameters",
			Details: fieldErr,
		})
	}

	lineItems, err := h.service.GetAll(query.AdvertiserID, query.Placement)
	if err != nil {
		h.log.Errorw("Failed to retrieve line items", "query", query, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to retrieve line items",
			Details: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(lineItems)
}
