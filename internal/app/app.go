package app

import (
	"sweng-task/internal/scheduler"
	"sweng-task/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"go.uber.org/zap"

	"sweng-task/internal/config"
	"sweng-task/internal/db"
	"sweng-task/internal/handler"
	"sweng-task/internal/repository/postgres"
	"sweng-task/internal/service"
)

func SetupApp(cfg *config.Config, log *zap.SugaredLogger) *fiber.App {
	// Database
	database := db.InitDatabase(cfg.Database, log)

	// Repositories
	lineItemRepo := postgres.NewLineItemPostgresRepository(database, log)
	trackingRepo := postgres.NewTrackingPostgresRepository(database, log)

	// Services
	lineItemService := service.NewLineItemService(lineItemRepo, log)
	trackingService := service.NewTrackingService(trackingRepo, lineItemService, log)
	adService := service.NewAdService(lineItemService, trackingService, log)

	// Handlers
	lineItemHandler := handler.NewLineItemHandler(lineItemService, log)
	adSelectionHandler := handler.NewAdSelectionHandler(adService, log)
	trackingHandler := handler.NewTrackingHandler(trackingService, log)

	// Fiber instance
	app := fiber.New(fiber.Config{
		AppName:      "Ad Bidding Service",
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 10,
		ErrorHandler: utils.ErrorHandler(log),
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	// Routes
	RegisterRoutes(app, lineItemHandler, adSelectionHandler, trackingHandler)

	// Schedulers
	schedule := scheduler.NewScheduler(lineItemService, log)
	schedule.Start()

	return app
}
