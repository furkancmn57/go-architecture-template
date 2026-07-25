package data

import (
	"database/sql"
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	_ "github.com/lib/pq"

	"github.com/furkancmn57/go-architecture-template/src/config"
)

// maintenanceDB is the catalog DB used only to create the application database.
const maintenanceDB = "postgres"

var safeDBName = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// EnsureDatabase creates POSTGRES_DB when it does not exist.
func EnsureDatabase(cfg config.Postgres) error {
	if !safeDBName.MatchString(cfg.DBName) {
		return fmt.Errorf("data: invalid database name %q", cfg.DBName)
	}

	admin := cfg
	admin.DBName = maintenanceDB

	db, err := sql.Open("postgres", admin.DSN())
	if err != nil {
		return fmt.Errorf(
			"data: failed to open maintenance db (host=%s port=%d user=%s db=%s): %w",
			admin.Host, admin.Port, admin.User, admin.DBName, err,
		)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf(
			"data: failed to connect to maintenance db (host=%s port=%d user=%s db=%s): %w",
			admin.Host, admin.Port, admin.User, admin.DBName, err,
		)
	}

	var exists bool
	if err := db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)",
		cfg.DBName,
	).Scan(&exists); err != nil {
		return fmt.Errorf("data: failed to check database existence: %w", err)
	}
	if exists {
		return nil
	}

	slog.Info("creating database", "name", cfg.DBName)
	quoted := `"` + strings.ReplaceAll(cfg.DBName, `"`, `""`) + `"`
	if _, err := db.Exec("CREATE DATABASE " + quoted); err != nil {
		return fmt.Errorf("data: failed to create database %q: %w", cfg.DBName, err)
	}
	return nil
}
