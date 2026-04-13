package test_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/database"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestTerraformDatabaseExample(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/terraform-database-example",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{},
	}

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.DestroyContext(t, t.Context(), terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApplyContext(t, t.Context(), terraformOptions)

	// Setting database configuration, including host, port, username, password and database name
	var dbConfig database.DBConfig

	dbConfig.Host = terraform.OutputContext(t, t.Context(), terraformOptions, "host")
	dbConfig.Port = terraform.OutputContext(t, t.Context(), terraformOptions, "port")
	dbConfig.User = terraform.OutputContext(t, t.Context(), terraformOptions, "username")
	dbConfig.Password = terraform.OutputContext(t, t.Context(), terraformOptions, "password")
	dbConfig.Database = terraform.OutputContext(t, t.Context(), terraformOptions, "database_name")

	// It can take a minute or so for the database to boot up, so retry a few times
	maxRetries := 15
	timeBetweenRetries := 15 * time.Second
	description := "Executing commands on database " + dbConfig.Host

	// Verify that we can connect to the database and run SQL commands
	retry.DoWithRetryContext(t, t.Context(), description, maxRetries, timeBetweenRetries, func() (string, error) {
		// Connect to specific database, i.e. postgres
		db, err := database.DBConnectionWithContextE(t, t.Context(), "postgres", &dbConfig)
		if err != nil {
			return "", err
		}

		defer db.Close()

		// Create a table
		creation := "create table person (id integer, name varchar(30), primary key (id))"

		if _, err := database.DBExecutionWithContextE(t, t.Context(), db, creation); err != nil {
			return "", err
		}

		// Insert a row
		expectedID := 12345
		expectedName := "azure"
		insertion := fmt.Sprintf("insert into person values (%d, '%s')", expectedID, expectedName)

		if _, err := database.DBExecutionWithContextE(t, t.Context(), db, insertion); err != nil {
			return "", err
		}

		// Query the table and check the output
		query := "select name from person"

		if err := database.DBQueryWithCustomValidationWithContextE(t, t.Context(), db, query, func(rows *sql.Rows) bool {
			var name string
			for rows.Next() {
				if scanErr := rows.Scan(&name); scanErr != nil {
					t.Fatal(scanErr)
				}

				if name != "azure" {
					return false
				}
			}

			return true
		}); err != nil {
			return "", err
		}

		// Drop the table
		drop := "drop table person"

		if _, err := database.DBExecutionWithContextE(t, t.Context(), db, drop); err != nil {
			return "", err
		}

		fmt.Println("Executed SQL commands correctly")

		return "", nil
	})
}
