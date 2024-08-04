package postgres

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunDBMigration(migrationURL string, dsn string) error {
	migration, err := migrate.New(migrationURL, dsn)
	if err != nil {
		return fmt.Errorf("failed create migrate: %w", err)
	}

	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed up migrate: %w", err)
	}

	return nil
}
