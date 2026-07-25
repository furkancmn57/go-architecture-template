package data

import (
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"

	"github.com/furkancmn57/go-base-template/src/config"
	"github.com/furkancmn57/go-base-template/src/data/migrations"
)

// schemaMigration tracks which code-first migrations have been applied.
type schemaMigration struct {
	Version   int64     `gorm:"primaryKey"`
	Name      string    `gorm:"size:255;not null"`
	AppliedAt time.Time `gorm:"not null"`
}

func (schemaMigration) TableName() string { return "schema_migrations" }

// Setup ensures the database exists, connects, and runs pending migrations.
func Setup(cfg config.Postgres) (*gorm.DB, error) {
	if err := EnsureDatabase(cfg); err != nil {
		return nil, err
	}
	db, err := New(cfg)
	if err != nil {
		return nil, err
	}
	if err := Migrate(db); err != nil {
		return nil, err
	}
	return db, nil
}

// Migrate applies pending code-first migrations and records each version in schema_migrations.
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&schemaMigration{}); err != nil {
		return fmt.Errorf("data: failed to ensure schema_migrations: %w", err)
	}

	pending, err := migrations.All()
	if err != nil {
		return fmt.Errorf("data: %w", err)
	}

	var applied []schemaMigration
	if err := db.Order("version asc").Find(&applied).Error; err != nil {
		return fmt.Errorf("data: failed to read schema_migrations: %w", err)
	}
	done := make(map[int64]struct{}, len(applied))
	for _, row := range applied {
		done[row.Version] = struct{}{}
	}

	for _, m := range pending {
		if _, ok := done[m.Version]; ok {
			continue
		}

		slog.Info("applying migration", "version", m.Version, "name", m.Name)
		err := db.Transaction(func(tx *gorm.DB) error {
			if err := m.Up(tx); err != nil {
				return err
			}
			return tx.Create(&schemaMigration{
				Version:   m.Version,
				Name:      m.Name,
				AppliedAt: time.Now().UTC(),
			}).Error
		})
		if err != nil {
			return fmt.Errorf("data: migration %d_%s failed: %w", m.Version, m.Name, err)
		}
		slog.Info("migration applied", "version", m.Version, "name", m.Name)
	}

	return nil
}
