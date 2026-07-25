package entities

import "github.com/furkancmn57/go-base-template/src/common"

// Todo is the app-facing persistence model (POCO). Schema/columns live in mappings.TodoMap.
// TableName must match mappings.TodoMap.TableName().
type Todo struct {
	common.Model
	Title       string
	Description string
	Completed   bool
}

// TableName is the physical table (kept here so entities does not import mappings).
func (Todo) TableName() string { return "todos" }
