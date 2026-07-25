package migrations

import "gorm.io/gorm"

// Migration is a versioned schema change tracked in schema_migrations.
type Migration struct {
	Version int64
	Name    string
	Up      func(tx *gorm.DB) error
}
