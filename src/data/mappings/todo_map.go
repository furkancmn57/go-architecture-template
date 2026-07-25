package mappings

import (
	"github.com/furkancmn57/go-base-template/src/data/entities"
	"gorm.io/gorm"
)

// TodoMap maps entities.Todo to the todos table (NotificationApi Fluent-style).
type TodoMap struct{}

// TableName is the physical table name.
func (TodoMap) TableName() string { return "todos" }

// Entity is the code-first model this map configures.
func (TodoMap) Entity() any { return &entities.Todo{} }

// Columns are Fluent column rules for domain fields (not common.Model).
func (TodoMap) Columns() []Column {
	return []Column{
		{Name: "title", Type: "varchar(255)", NotNull: true},
		{Name: "description", Type: "text"},
		{Name: "completed", Type: "boolean", NotNull: true, Default: "false"},
	}
}

// Migrate applies this map's schema (used from versioned migrations).
func (TodoMap) Migrate(tx *gorm.DB) error {
	if err := tx.AutoMigrate(TodoMap{}.Entity()); err != nil {
		return err
	}
	return ApplyColumns(tx, TodoMap{}.TableName(), TodoMap{}.Columns())
}
