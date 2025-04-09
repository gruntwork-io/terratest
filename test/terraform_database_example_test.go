package test

import (
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
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Setting database configuration, including host, port, username, password and database name
	var dbConfig database.DBConfig
	dbConfig.Host = terraform.Output(t, terraformOptions, "host")
	dbConfig.Port = terraform.Output(t, terraformOptions, "port")
	dbConfig.User = terraform.Output(t, terraformOptions, "username")
	dbConfig.Password = terraform.Output(t, terraformOptions, "password")
	dbConfig.Database = terraform.Output(t, terraformOptions, "database_name")

	// It can take a minute or so for the database to boot up, so retry a few times
	maxRetries := 15
	timeBetweenRetries := 15 * time.Second
	description := fmt.Sprintf("Executing commands on database %s", dbConfig.Host)

	// Verify that we can connect to the database and run SQL commands
	retry.DoWithRetry(t, description, maxRetries, timeBetweenRetries, func() (string, error) {
		// Connect to specific database, i.e. postgres
		db, err := database.DBConnectionE(t, "postgres", dbConfig)
		if err != nil {
			return "", err
		}

		// Create a table
		creation := "create table person (id integer, name varchar(30), primary key (id))"
		database.DBExecution(t, db, creation)

		// Insert a row
		expectedID := 12345
		expectedName := "azure"
		insertion := fmt.Sprintf("insert into person values (%d, '%s')", expectedID, expectedName)
		database.DBExecution(t, db, insertion)

		// Query the table and check the output
		query := "select name from person"
		database.DBQueryWithValidation(t, db, query, "azure")

		// Drop the table
		drop := "drop table person"
		database.DBExecution(t, db, drop)
		fmt.Println("Executed SQL commands correctly")

		defer db.Close()

		return "", nil
	})
}
