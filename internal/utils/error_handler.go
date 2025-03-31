package utils

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func ErrorHandler(log *zap.SugaredLogger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		// Log the error
		log.Errorw("Unhandled error", "path", c.Path(), "error", err)

		// Determine status code
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		// Send consistent JSON error
		return c.Status(code).JSON(ErrorResponse{
			Code:    code,
			Message: "An unexpected error occurred",
			Details: err.Error(),
		})
	}
}
