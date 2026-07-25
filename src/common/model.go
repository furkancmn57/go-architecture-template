package common

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Model is the base struct every GORM entity embeds.
// Domain column types live in data/mappings Fluent maps; only shared PK / soft-delete tags stay here.
type Model struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate assigns a UUID before the row is inserted, unless one has
// already been set by the caller.
func (m *Model) BeforeCreate(_ *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
