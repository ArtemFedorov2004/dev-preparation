package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/joho/godotenv"

	"github.com/devprep/backend/internal/config"
	"github.com/devprep/backend/migrations"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	m := mustNewMigrator(cfg.Database.DSN())
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			slog.Warn("migration source close", "error", srcErr)
		}
		if dbErr != nil {
			slog.Warn("migration db close", "error", dbErr)
		}
	}()

	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			fatal("migrate up", err)
		}
		printVersion(m)

	case "down":
		steps := 1
		if len(os.Args) >= 3 {
			n, err := strconv.Atoi(os.Args[2])
			if err != nil || n < 1 {
				fatal("invalid steps argument", err)
			}
			steps = n
		}
		if err := m.Steps(-steps); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			fatal("migrate down", err)
		}
		printVersion(m)

	case "version":
		printVersion(m)

	case "force":
		if len(os.Args) < 3 {
			fatal("force requires a version argument", nil)
		}
		v, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fatal("invalid version", err)
		}
		if err := m.Force(v); err != nil {
			fatal("migrate force", err)
		}
		printVersion(m)

	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func mustNewMigrator(dsn string) *migrate.Migrate {
	migrateDSN := "pgx5://" + strings.TrimPrefix(dsn, "postgres://")

	src, err := iofs.New(migrations.FS, ".")
	if err != nil {
		fatal("create migration source", err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", src, migrateDSN)
	if err != nil {
		fatal("create migrator", err)
	}
	return m
}

func printVersion(m *migrate.Migrate) {
	v, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		slog.Warn("could not read version", "error", err)
		return
	}
	slog.Info("current migration version", "version", v, "dirty", dirty)
}

func fatal(msg string, err error) {
	slog.Error(msg, "error", err)
	os.Exit(1)
}

func printUsage() {
	fmt.Fprintln(os.Stderr, `Usage:
  go run ./cmd/migrate <command> [args]

Commands:
  up           Apply all pending migrations
  down [N]     Roll back N steps (default: 1)
  version      Print current schema version
  force <N>    Force-set version (fixes dirty state)`)
}
