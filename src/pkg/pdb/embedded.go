package pdb

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var dbRef *EmbeddedDB

type EmbeddedDB struct {
	*Queries
	DB   *sql.DB
	path string
}

const (
	dbGetErrFmtStr  = "db get err: %w"
	dbSyncErrFmtStr = "db sync err: %w"
)

// Close prepared statements and database conn
func (embedded *EmbeddedDB) Close() error {
	statementsCloseErr := embedded.Queries.Close()
	dbCloseErr := embedded.DB.Close()
	rmDBPathErr := os.RemoveAll(embedded.path)
	return errors.Join(statementsCloseErr, dbCloseErr, rmDBPathErr)
}

const (
	sqlite3String   = "sqlite3"
	dbConnStringFmt = "file:%s?cache=shared&mode=rwc&_synchronous=NORMAL&_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=true"
)

func Get() (*EmbeddedDB, error) {
	if !globalContextAssigned {
		return nil, fmt.Errorf(dbSyncErrFmtStr, errGlobalContextNotAssigned)
	}
	if dbRef != nil {
		return dbRef, nil
	}

	// connector, connectorErr := libsql.NewEmbeddedReplicaConnector(dbPath, dbURL, libsql.WithAuthToken(dbToken), libsql.WithReadYourWrites(true))
	connectionString := fmt.Sprintf(dbConnStringFmt, dbPath)

	db, dbErr := sql.Open(sqlite3String, connectionString)
	if dbErr != nil {
		return nil, dbErr
	}

	// schema verify
	schemaErr := applyMigrations(db)
	if schemaErr != nil {
		return nil, fmt.Errorf(dbSyncErrFmtStr, schemaErr)
	}

	querier, querierErr := Prepare(globalContext, db)
	if querierErr != nil {
		return nil, fmt.Errorf(dbSyncErrFmtStr, querierErr)
	}

	embeddedDB := EmbeddedDB{
		Queries: querier,
		DB:      db,
		path:    dbPath,
	}

	dbRef = &embeddedDB

	return dbRef, nil
}

func Obtain() *EmbeddedDB {
	return dbRef
}
