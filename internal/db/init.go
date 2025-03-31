package db

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sweng-task/internal/config"
)

func InitDatabase(cfg config.DatabaseConfig, log *zap.SugaredLogger) *gorm.DB {
	db, err := ConnectPostgres(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB from GORM: %v", err)

	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}
	log.Info("Connected to PostgreSQL")

	if err := RunMigrations(db, log); err != nil {

		log.Fatalf("Migration error: %v", err)

	}
	log.Info("Migrations applied successfully")

	return db
}
