package migrations

import (
	"gorm.io/gorm"

	"github.com/furkancmn57/go-base-template/src/data/mappings"
)

func init() {
	register(Migration{
		Version: 1,
		Name:    "create_todos",
		Up: func(tx *gorm.DB) error {
			return mappings.TodoMap{}.Migrate(tx)
		},
	})
}
