package db

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sweng-task/internal/model"
)

func RunMigrations(db *gorm.DB, log *zap.SugaredLogger) error {
	log.Info("Running GORM migrations")

	return db.AutoMigrate(
		&model.LineItemEntity{},
		&model.TrackingEventEntity{},
	)

}
