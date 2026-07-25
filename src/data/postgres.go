package data

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/furkancmn57/go-base-template/src/config"
)

// New opens a *gorm.DB connection to Postgres.
func New(cfg config.Postgres) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf(
			"data: failed to connect to postgres (host=%s port=%d user=%s db=%s): %w",
			cfg.Host, cfg.Port, cfg.User, cfg.DBName, err,
		)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("data: failed to get underlying sql.DB: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("data: failed to ping postgres: %w", err)
	}

	return db, nil
}
