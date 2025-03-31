package app

import (
	"github.com/gofiber/fiber/v2"

	"sweng-task/internal/handler"
)

func RegisterRoutes(app *fiber.App,
	lineItemHandler *handler.LineItemHandler,
	adSelectionHandler *handler.AdSelectionHandler,
	trackingHandler *handler.TrackingHandler,
) {
	app.Get("/health", handler.HealthCheck)

	api := app.Group("/api/v1")

	// Line items
	api.Post("/lineitems", lineItemHandler.Create)
	api.Get("/lineitems", lineItemHandler.GetAll)
	api.Get("/lineitems/:id", lineItemHandler.GetByID)

	// Ad selection
	api.Get("/ads", adSelectionHandler.GetWinningAds)

	// Tracking
	api.Post("/tracking", trackingHandler.TrackEvent)
}
