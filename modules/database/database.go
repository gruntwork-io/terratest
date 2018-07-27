package database

import (
	"database/sql"
	"fmt"
	"testing"

	// Microsoft SQL Database Driver
	_ "github.com/denisenkom/go-mssqldb"

	// PostgreSQL Database Driver
	_ "github.com/lib/pq"

	// MySQL Database Driver
	_ "github.com/go-sql-driver/mysql"
)

// DBConfig using server name, user name, password and database name
type DBConfig struct {
	host     string
	port     int
	user     string
	password string
	database string
}

// DBConnection connects to the database using database configuration and database type, i.e. mssql, and then return the database. If there's any error, fail the test.
func DBConnection(t *testing.T, dbType string, dbConfig DBConfig) *sql.DB {
	db, err := DBConnectionE(t, dbType, dbConfig)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

// DBConnectionE connects to the database using database configuration and database type, i.e. mssql. Return the database or an error.
func DBConnectionE(t *testing.T, dbType string, dbConfig DBConfig) (*sql.DB, error) {
	config := ""
	switch dbType {
	case "mssql":
		config = fmt.Sprintf("server = %s; port = %d; user id = %s; password = %s; database = %s", dbConfig.host, dbConfig.port, dbConfig.user, dbConfig.password, dbConfig.database)
	case "postgres":
		config = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", dbConfig.host, dbConfig.port, dbConfig.user, dbConfig.password, dbConfig.database)
	case "mysql":
		config = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?allowNativePasswords=true", dbConfig.user, dbConfig.password, dbConfig.host, dbConfig.port, dbConfig.database)
	default:
		return nil, DBUnknown{dbType: dbType}
	}
	db, err := sql.Open(dbType, config)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

// DBExecution executes specific SQL commands, i.e. insertion. If there's any error, fail the test.
func DBExecution(t *testing.T, db *sql.DB, command string) {
	_, err := DBExecutionE(t, db, command)
	if err != nil {
		t.Fatal(err)
	}
}

// DBExecutionE executes specific SQL commands, i.e. insertion. Return the result or an error.
func DBExecutionE(t *testing.T, db *sql.DB, command string) (sql.Result, error) {
	result, err := db.Exec(command)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DBQuery queries from database, i.e. selection, and then return the result. If there's any error, fail the test.
func DBQuery(t *testing.T, db *sql.DB, command string) *sql.Rows {
	rows, err := DBQueryE(t, db, command)
	if err != nil {
		t.Fatal(err)
	}
	return rows
}

// DBQueryE queries from database, i.e. selection. Return the result or an error.
func DBQueryE(t *testing.T, db *sql.DB, command string) (*sql.Rows, error) {
	rows, err := db.Query(command)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// DBQueryWithValidation queries from database and validate whether the result is the same as expected text. If there's any error, fail the test.
func DBQueryWithValidation(t *testing.T, db *sql.DB, command string, expected string) {
	err := DBQueryWithValidationE(t, db, command, expected)
	if err != nil {
		t.Fatal(err)
	}
}

// DBQueryWithValidationE queries from database and validate whether the result is the same as expected text. If not, return an error.
func DBQueryWithValidationE(t *testing.T, db *sql.DB, command string, expected string) error {
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

// DBQueryWithCustomValidation queries from database and validate whether the result meets the requirement. If there's any error, fail the test.
func DBQueryWithCustomValidation(t *testing.T, db *sql.DB, command string, validateResponse func(*sql.Rows) bool) {
	err := DBQueryWithCustomValidationE(t, db, command, validateResponse)
	if err != nil {
		t.Fatal(err)
	}
}

// DBQueryWithCustomValidationE queries from database and validate whether the result meets the requirement. If not, return an error.
func DBQueryWithCustomValidationE(t *testing.T, db *sql.DB, command string, validateResponse func(*sql.Rows) bool) error {
	rows, err := DBQueryE(t, db, command)
	defer rows.Close()
	if err != nil {
		return err
	}
	if !validateResponse(rows) {
		return ValidationFunctionFailed{command: command}
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
