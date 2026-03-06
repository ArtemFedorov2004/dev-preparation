package database

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"github.com/devprep/backend/migrations"
)

func RunMigrations(dsn string) error {
	src, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("create migration source: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", src, dsn)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			slog.Warn("migration source close error", "error", srcErr)
		}
		if dbErr != nil {
			slog.Warn("migration db close error", "error", dbErr)
		}
	}()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("migrations: no changes")
			return nil
		}
		return fmt.Errorf("run migrations: %w", err)
	}

	version, dirty, _ := m.Version()
	slog.Info("migrations applied", "version", version, "dirty", dirty)
	return nil
}
