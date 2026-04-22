// Package database provides helper functions to connect to and interact with databases during automated tests.
package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"

	// Microsoft SQL Database Driver (Microsoft-maintained fork of denisenkom/go-mssqldb)
	_ "github.com/microsoft/go-mssqldb"

	// PostgreSQL Database Driver
	_ "github.com/lib/pq"

	// MySQL Database Driver
	_ "github.com/go-sql-driver/mysql"
)

const (
	_databaseTypeMSSQL    = "mssql"
	_databaseTypePostgres = "postgres"
	_databaseTypeMySQL    = "mysql"
	_postgresConnStr      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable"
	_mssqlConnStr         = "server = %s; port = %s; user id = %s; password = %s; database = %s"
	_mysqlConnStr         = "%s:%s@tcp(%s:%s)/%s?allowNativePasswords=true"
)

// DBConfig using server name, user name, password and database name.
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// DBConnection connects to the database using database configuration and database type, i.e. mssql, and then returns
// the database. If there's any error, fail the test.
//
// Deprecated: Use DBConnectionWithContextE instead.
func DBConnection(t testing.TestingT, dbType string, dbConfig DBConfig) *sql.DB { //nolint:gocritic // Preserving original signature for backward compatibility.
	t.Helper()

	db, err := DBConnectionE(t, dbType, dbConfig)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

// DBConnectionE connects to the database using database configuration and database type, i.e. mssql. Returns the
// database or an error.
//
// Deprecated: Use DBConnectionWithContextE instead.
func DBConnectionE(t testing.TestingT, dbType string, dbConfig DBConfig) (*sql.DB, error) { //nolint:gocritic // Preserving original signature for backward compatibility.
	t.Helper()

	return DBConnectionWithContextE(t, context.Background(), dbType, &dbConfig)
}

// DBConnectionWithContextE connects to the database using database configuration and database type, i.e. mssql.
// Returns the database or an error.
func DBConnectionWithContextE(t testing.TestingT, ctx context.Context, dbType string, dbConfig *DBConfig) (*sql.DB, error) {
	t.Helper()

	var config string

	switch dbType {
	case _databaseTypeMSSQL:
		config = fmt.Sprintf(_mssqlConnStr, dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Database)
	case _databaseTypePostgres:
		config = fmt.Sprintf(_postgresConnStr, dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Database)
	case _databaseTypeMySQL:
		config = fmt.Sprintf(_mysqlConnStr, dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)
	default:
		return nil, DBUnknown{dbType: dbType}
	}

	db, err := sql.Open(dbType, config)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// DBExecution executes specific SQL commands, i.e. insertion. If there's any error, fail the test.
//
// Deprecated: Use DBExecutionWithContextE instead.
func DBExecution(t testing.TestingT, db *sql.DB, command string) {
	t.Helper()

	_, err := DBExecutionE(t, db, command)
	if err != nil {
		t.Fatal(err)
	}
}

// DBExecutionE executes specific SQL commands, i.e. insertion. Returns the result or an error.
//
// Deprecated: Use DBExecutionWithContextE instead.
func DBExecutionE(t testing.TestingT, db *sql.DB, command string) (sql.Result, error) {
	t.Helper()

	return DBExecutionWithContextE(t, context.Background(), db, command)
}

// DBExecutionWithContextE executes specific SQL commands, i.e. insertion. Returns the result or an error.
func DBExecutionWithContextE(t testing.TestingT, ctx context.Context, db *sql.DB, command string) (sql.Result, error) {
	t.Helper()

	result, err := db.ExecContext(ctx, command)
	if err != nil {
		return nil, fmt.Errorf("failed to execute database command: %w", err)
	}

	return result, nil
}

// DBQuery queries from database, i.e. selection, and then returns the result. If there's any error, fail the test.
//
// Deprecated: Use DBQueryWithContextE instead.
func DBQuery(t testing.TestingT, db *sql.DB, command string) *sql.Rows {
	t.Helper()

	rows, err := DBQueryE(t, db, command)
	if err != nil {
		t.Fatal(err)
	}

	return rows
}

// DBQueryE queries from database, i.e. selection. Returns the result or an error.
//
// Deprecated: Use DBQueryWithContextE instead.
func DBQueryE(t testing.TestingT, db *sql.DB, command string) (*sql.Rows, error) {
	t.Helper()

	return DBQueryWithContextE(t, context.Background(), db, command)
}

// DBQueryWithContextE queries from database, i.e. selection. Returns the result or an error.
func DBQueryWithContextE(t testing.TestingT, ctx context.Context, db *sql.DB, command string) (*sql.Rows, error) {
	t.Helper()

	rows, err := db.QueryContext(ctx, command)
	if err != nil {
		return nil, fmt.Errorf("failed to execute database query: %w", err)
	}

	return rows, nil
}

// DBQueryWithValidation queries from database and validates whether the result is the same as expected text. If
// there's any error, fail the test.
//
// Deprecated: Use DBQueryWithCustomValidationWithContextE instead.
func DBQueryWithValidation(t testing.TestingT, db *sql.DB, command string, expected string) {
	t.Helper()

	err := DBQueryWithValidationE(t, db, command, expected)
	if err != nil {
		t.Fatal(err)
	}
}

// DBQueryWithValidationE queries from database and validates whether the result is the same as expected text. If not,
// returns an error.
//
// Deprecated: Use DBQueryWithCustomValidationWithContextE instead.
func DBQueryWithValidationE(t testing.TestingT, db *sql.DB, command string, expected string) error {
	t.Helper()

	return DBQueryWithCustomValidationE(t, db, command, func(rows *sql.Rows) bool {
		var name string
		for rows.Next() {
			err := rows.Scan(&name)
			if err != nil {
				t.Fatal(err)
			}

			if name != expected {
				return false
			}
		}

		return true
	})
}

// DBQueryWithCustomValidation queries from database and validates whether the result meets the requirement. If
// there's any error, fail the test.
//
// Deprecated: Use DBQueryWithCustomValidationWithContextE instead.
func DBQueryWithCustomValidation(t testing.TestingT, db *sql.DB, command string, validateResponse func(*sql.Rows) bool) {
	t.Helper()

	err := DBQueryWithCustomValidationE(t, db, command, validateResponse)
	if err != nil {
		t.Fatal(err)
	}
}

// DBQueryWithCustomValidationE queries from database and validates whether the result meets the requirement. If not,
// returns an error.
//
// Deprecated: Use DBQueryWithCustomValidationWithContextE instead.
func DBQueryWithCustomValidationE(t testing.TestingT, db *sql.DB, command string, validateResponse func(*sql.Rows) bool) error {
	t.Helper()

	return DBQueryWithCustomValidationWithContextE(t, context.Background(), db, command, validateResponse)
}

// DBQueryWithCustomValidationWithContextE queries from database and validates whether the result meets the requirement.
// If not, returns an error.
func DBQueryWithCustomValidationWithContextE(t testing.TestingT, ctx context.Context, db *sql.DB, command string, validateResponse func(*sql.Rows) bool) error {
	t.Helper()

	rows, err := DBQueryWithContextE(t, ctx, db, command)
	if err != nil {
		return err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			logger.Default.Logf(t, "failed to close database rows: %v", closeErr)
		}
	}()

	if !validateResponse(rows) {
		return ValidationFunctionFailed{command: command}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating database rows: %w", err)
	}

	return nil
}

// ValidationFunctionFailed is an error that occurs if the validation function fails.
type ValidationFunctionFailed struct {
	command string
}

func (err ValidationFunctionFailed) Error() string {
	return fmt.Sprintf("Validation failed for command: %s.", err.command)
}

// DBUnknown is an error that occurs if the given database type is unknown or not supported.
type DBUnknown struct {
	dbType string
}

func (err DBUnknown) Error() string {
	return fmt.Sprintf("Database unknown or not supported: %s. We only support mssql, postgres and mysql.", err.dbType)
}
