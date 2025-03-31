package handler

import (
	"sweng-task/internal/model"
	"sweng-task/internal/service"
	"sweng-task/internal/utils"
	"sweng-task/internal/validator"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type TrackingHandler struct {
	service *service.TrackingService
	logger  *zap.SugaredLogger
}

func NewTrackingHandler(service *service.TrackingService, logger *zap.SugaredLogger) *TrackingHandler {
	return &TrackingHandler{service: service, logger: logger}
}

func (h *TrackingHandler) TrackEvent(c *fiber.Ctx) error {
	var event model.TrackingEvent

	if err := c.BodyParser(&event); err != nil {
		h.logger.Warnw("Invalid tracking payload", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	if fieldErr, err := validator.ValidateStruct(&event); err != nil {
		h.logger.Warnw("Validation failed", "error", err)
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

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	if err := h.service.Track(event); err != nil {
		h.logger.Errorw("Failed to store tracking event", "error", err)

		if err == service.ErrLineItemNotFound {
			return c.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse{
				Code:    fiber.StatusNotFound,
				Message: "Line item not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to track event",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"success": true})
}
