// Package pdb wraps the libSQL database for postc.art
// It has support for embedded (local copy for fast reads).
package pdb

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"os"

	"github.com/pressly/goose/v3"
)

var (
	dbURL   string
	dbToken string

	globalContext         context.Context
	globalContextAssigned bool

	temporaryDir string
	dbPath       string
)

var (
	errBadDatabaseURL           = errors.New("db err: bad turso database url provided to db configure")
	errBadDatabaseToken         = errors.New("db err: bad turso database token provided to db configure")
	errGlobalContextNotAssigned = errors.New("db err: global context not yet assigned before using db")
)

const (
	dbConfigureErrFmtStr       = "db configuration err: %w"
	dbApplyMigrationsErrFmtStr = "db migrations err: %w"
)

func Close() error {
	var errs []error
	if dbRef != nil {
		dbErr := dbRef.Close()
		if dbErr != nil {
			errs = append(errs, dbErr)
		}
	}
	if temporaryDir != "" {
		rmErr := os.RemoveAll(temporaryDir)
		if rmErr != nil {
			errs = append(errs, rmErr)
		}
	}

	return errors.Join(errs...)
}

func Configure(ctx context.Context, tursoURL, tursoToken string) error {
	if tursoURL == "" {
		return fmt.Errorf(dbConfigureErrFmtStr, errBadDatabaseURL)
	}
	if tursoToken == "" {
		return fmt.Errorf(dbConfigureErrFmtStr, errBadDatabaseToken)
	}
	dbURL = tursoURL
	dbToken = tursoToken

	globalContext = ctx
	globalContextAssigned = true

	var tempDirErr error
	temporaryDir, tempDirErr = os.MkdirTemp("", "libsql-*")
	if tempDirErr != nil {
		return fmt.Errorf(dbConfigureErrFmtStr, tempDirErr)
	}

	// dbPath = path.Join(temporaryDir, "postcart.db")
	dbPath = "./postcart.db"

	_, dbGetErr := Get()
	if dbGetErr != nil {
		return fmt.Errorf(dbConfigureErrFmtStr, dbGetErr)
	}

	return nil
}

//go:embed sql/migrations/*.sql
var embedMigrations embed.FS

func applyMigrations(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	dialectErr := goose.SetDialect("turso")
	if dialectErr != nil {
		return fmt.Errorf(dbApplyMigrationsErrFmtStr, dialectErr)
	}

	migrations, collectErr := goose.CollectMigrations("sql/migrations", 0, goose.MaxVersion)
	if collectErr != nil {
		return fmt.Errorf(dbApplyMigrationsErrFmtStr, collectErr)
	}

	version, err := goose.EnsureDBVersion(db)
	if err != nil {
		return fmt.Errorf(dbApplyMigrationsErrFmtStr, err)
	}

	if int64(len(migrations)) > version {
		for _, migration := range migrations[version:] {
			if err := migration.Up(db); err != nil {
				return fmt.Errorf(dbApplyMigrationsErrFmtStr, fmt.Errorf("failed to apply migration %d: %w", migration.Version, err))
			}
		}
	}

	return nil
}
