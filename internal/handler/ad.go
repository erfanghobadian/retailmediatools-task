package handler

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"sweng-task/internal/service"
	"sweng-task/internal/utils"
	"sweng-task/internal/validator"
)

type AdSelectionHandler struct {
	adService *service.AdService
	log       *zap.SugaredLogger
}

func NewAdSelectionHandler(adService *service.AdService, logger *zap.SugaredLogger) *AdSelectionHandler {
	return &AdSelectionHandler{
		adService: adService,
		log:       logger,
	}
}

// GetWinningAds handles GET /ads requests
func (h *AdSelectionHandler) GetWinningAds(c *fiber.Ctx) error {
	const defaultLimit = 1

	var q validator.AdQueryParams
	if err := c.QueryParser(&q); err != nil {
		h.log.Warnw("Failed to parse query parameters", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid query parameters",
			Details: err.Error(),
		})
	}

	// Default fallback for limit
	if q.Limit == 0 {
		q.Limit = defaultLimit
	}

	if fieldErr, err := validator.ValidateStruct(&q); err != nil {
		h.log.Warnw("Validation failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Validation failed",
		})
	} else if fieldErr != nil {
		h.log.Warnw("Query parameter validation error", "field", fieldErr.Field, "reason", fieldErr.Reason)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request parameters",
			Details: fieldErr,
		})
	}

	h.log.Infow("Received ad request",
		"placement", q.Placement,
		"category", q.Category,
		"keyword", q.Keyword,
		"limit", q.Limit,
	)

	ads, err := h.adService.GetWinningAds(q.Placement, q.Category, q.Keyword, q.Limit)
	if err != nil {
		h.log.Errorw("Failed to get winning ads", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to retrieve ads",
			Details: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(ads)
}
