// Package extensions wires infrastructure at the composition root.
package extensions

import (
	"gorm.io/gorm"

	"github.com/furkancmn57/go-architecture-template/src/config"
	"github.com/furkancmn57/go-architecture-template/src/data"
)

// AddDatabase creates the DB if missing, opens Postgres, and runs versioned migrations.
func AddDatabase(cfg config.Postgres) (*gorm.DB, error) {
	return data.Setup(cfg)
}
