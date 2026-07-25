package mappings

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// Column is a Fluent column rule (EF Property().IsRequired() / HasMaxLength equivalent).
type Column struct {
	Name    string
	Type    string // e.g. "varchar(255)", "text", "boolean"
	NotNull bool
	Default string // raw SQL default expression, e.g. "false", "''"
}

// ApplyColumns aligns table columns with Fluent rules after AutoMigrate.
func ApplyColumns(tx *gorm.DB, table string, columns []Column) error {
	if table == "" {
		return fmt.Errorf("mappings: empty table name")
	}
	for _, col := range columns {
		if col.Name == "" || col.Type == "" {
			return fmt.Errorf("mappings: column on %q needs Name and Type", table)
		}
		if err := applyColumn(tx, table, col); err != nil {
			return err
		}
	}
	return nil
}

func applyColumn(tx *gorm.DB, table string, col Column) error {
	qTable := quoteIdent(table)
	qCol := quoteIdent(col.Name)

	stmt := fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s USING %s::%s",
		qTable, qCol, col.Type, qCol, baseType(col.Type))
	if err := tx.Exec(stmt).Error; err != nil {
		return fmt.Errorf("mappings: alter type %s.%s: %w", table, col.Name, err)
	}

	if col.NotNull {
		stmt = fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s SET NOT NULL", qTable, qCol)
	} else {
		stmt = fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s DROP NOT NULL", qTable, qCol)
	}
	if err := tx.Exec(stmt).Error; err != nil {
		return fmt.Errorf("mappings: alter nullability %s.%s: %w", table, col.Name, err)
	}

	if col.Default != "" {
		stmt = fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s SET DEFAULT %s", qTable, qCol, col.Default)
		if err := tx.Exec(stmt).Error; err != nil {
			return fmt.Errorf("mappings: set default %s.%s: %w", table, col.Name, err)
		}
	} else {
		stmt = fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s DROP DEFAULT", qTable, qCol)
		if err := tx.Exec(stmt).Error; err != nil {
			return fmt.Errorf("mappings: drop default %s.%s: %w", table, col.Name, err)
		}
	}

	return nil
}

func quoteIdent(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// baseType strips length/precision for USING casts (varchar(255) → varchar).
func baseType(t string) string {
	if i := strings.IndexByte(t, '('); i >= 0 {
		return t[:i]
	}
	return t
}
