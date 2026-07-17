package aws

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

// GetAddressOfRdsInstanceContextE gets the address of the given RDS Instance in the given region.
// The ctx parameter supports cancellation and timeouts.
func GetAddressOfRdsInstanceContextE(t testing.TestingT, ctx context.Context, dbInstanceID string, awsRegion string) (string, error) {
	dbInstance, err := GetRdsInstanceDetailsContextE(t, ctx, dbInstanceID, awsRegion)
	if err != nil {
		return "", err
	}

	if dbInstance.Endpoint == nil {
		return "", fmt.Errorf("RDS instance %s endpoint is not yet available", dbInstanceID)
	}

	return aws.ToString(dbInstance.Endpoint.Address), nil
}

// GetAddressOfRdsInstanceContext gets the address of the given RDS Instance in the given region.
// The ctx parameter supports cancellation and timeouts.
func GetAddressOfRdsInstanceContext(t testing.TestingT, ctx context.Context, dbInstanceID string, awsRegion string) string {
	t.Helper()

	address, err := GetAddressOfRdsInstanceContextE(t, ctx, dbInstanceID, awsRegion)
	require.NoError(t, err)

	return address
}

// GetPortOfRdsInstanceContextE gets the port of the given RDS Instance in the given region.
// The ctx parameter supports cancellation and timeouts.
func GetPortOfRdsInstanceContextE(t testing.TestingT, ctx context.Context, dbInstanceID string, awsRegion string) (int32, error) {
	dbInstance, err := GetRdsInstanceDetailsContextE(t, ctx, dbInstanceID, awsRegion)
	if err != nil {
		return -1, err
	}

	if dbInstance.Endpoint == nil {
		return -1, fmt.Errorf("RDS instance %s endpoint is not yet available", dbInstanceID)
	}

	return *dbInstance.Endpoint.Port, nil
}

// GetPortOfRdsInstanceContext gets the port of the given RDS Instance in the given region.
// The ctx parameter supports cancellation and timeouts.
func GetPortOfRdsInstanceContext(t testing.TestingT, ctx context.Context, dbInstanceID string, awsRegion string) int32 {
	t.Helper()

	port, err := GetPortOfRdsInstanceContextE(t, ctx, dbInstanceID, awsRegion)
	require.NoError(t, err)

	return port
}

// GetWhetherSchemaExistsInRdsMySQLInstanceContextE checks whether the specified schema/table name exists in the RDS MySQL instance.
// The ctx parameter supports cancellation and timeouts.
func GetWhetherSchemaExistsInRdsMySQLInstanceContextE(t testing.TestingT, ctx context.Context, dbURL string, dbPort int32, dbUsername string, dbPassword string, expectedSchemaName string) (bool, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/", dbUsername, dbPassword, dbURL, dbPort)

	db, connErr := sql.Open("mysql", connectionString)
	if connErr != nil {
		return false, connErr
	}

	defer func() { _ = db.Close() }()

	var schemaName string

	sqlStatement := "SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME=?;"
	row := db.QueryRowContext(ctx, sqlStatement, expectedSchemaName)

	scanErr := row.Scan(&schemaName)
	if scanErr != nil {
		return false, scanErr
	}

	return true, nil
}

// GetWhetherSchemaExistsInRdsMySQLInstanceContext checks whether the specified schema/table name exists in the RDS MySQL instance.
// The ctx parameter supports cancellation and timeouts.
func GetWhetherSchemaExistsInRdsMySQLInstanceContext(t testing.TestingT, ctx context.Context, dbURL string, dbPort int32, dbUsername string, dbPassword string, expectedSchemaName string) bool {
	t.Helper()

	output, err := GetWhetherSchemaExistsInRdsMySQLInstanceContextE(t, ctx, dbURL, dbPort, dbUsername, dbPassword, expectedSchemaName)
	require.NoError(t, err)

	return output
}

// GetWhetherSchemaExistsInRdsPostgresInstanceContextE checks whether the specified schema/table name exists in the RDS Postgres instance.
// The ctx parameter supports cancellation and timeouts.
func GetWhetherSchemaExistsInRdsPostgresInstanceContextE(t testing.TestingT, ctx context.Context, dbURL string, dbPort int32, dbUsername string, dbPassword string, expectedSchemaName string) (bool, error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", dbURL, dbPort, dbUsername, dbPassword, expectedSchemaName)

	db, connErr := sql.Open("pgx", connectionString)
	if connErr != nil {
		return false, connErr
	}

	defer func() { _ = db.Close() }()

	var schemaName string

	sqlStatement := `SELECT "catalog_name" FROM "information_schema"."schemata" where catalog_name=$1`
	row := db.QueryRowContext(ctx, sqlStatement, expectedSchemaName)

	scanErr := row.Scan(&schemaName)
	if scanErr != nil {
		return false, scanErr
	}

	return true, nil
}

// GetWhetherSchemaExistsInRdsPostgresInstanceContext checks whether the specified schema/table name exists in the RDS Postgres instance.
// The ctx parameter supports cancellation and timeouts.
func GetWhetherSchemaExistsInRdsPostgresInstanceContext(t testing.TestingT, ctx context.Context, dbURL string, dbPort int32, dbUsername string, dbPassword string, expectedSchemaName string) bool {
	t.Helper()

	output, err := GetWhetherSchemaExistsInRdsPostgresInstanceContextE(t, ctx, dbURL, dbPort, dbUsername, dbPassword, expectedSchemaName)
	require.NoError(t, err)

	return output
}

// GetParameterValueForParameterOfRdsInstanceContextE gets the value of the parameter name specified for the RDS instance in the given region.
// The ctx parameter supports cancellation and timeouts.
func GetParameterValueForParameterOfRdsInstanceContextE(t testing.TestingT, ctx context.Context, parameterName string, dbInstanceID string, awsRegion string) (string, error) {
	output, err := GetAllParametersOfRdsInstanceContextE(t, ctx, dbInstanceID, awsRegion)
	if err != nil {
		return "", err
	}

	for _, parameter := range output {
		if aws.ToString(parameter.ParameterName) == parameterName {
			return aws.ToString(parameter.ParameterValue), nil
		}
	}

	return "", ParameterForDBInstanceNotFound{ParameterName: parameterName, DbInstanceID: dbInstanceID, AwsRegion: awsRegion}
}

// GetParameterValueForParameterOfRdsInstanceContext gets the value of the parameter name specified for the RDS instance in the given region.
// The ctx parameter supports cancellation and timeouts.
func GetParameterValueForParameterOfRdsInstanceContext(t testing.TestingT, ctx context.Context, parameterName string, dbInstanceID string, awsRegion string) string {
	t.Helper()

	parameterValue, err := GetParameterValueForParameterOfRdsInstanceContextE(t, ctx, parameterName, dbInstanceID, awsRegion)
	require.NoError(t, err)

	return parameterValue
}

// GetOptionSettingForOfRdsInstanceContextE gets the value of the option name in the option group specified for the RDS instance in the given region.
// The ctx parameter supports cancellation and timeouts.
func GetOptionSettingForOfRdsInstanceContextE(t testing.TestingT, ctx context.Context, optionName string, optionSettingName string, dbInstanceID, awsRegion string) (string, error) {
	optionGroupName, err := GetOptionGroupNameOfRdsInstanceContextE(t, ctx, dbInstanceID, awsRegion)
	if err != nil {
		return "", err
	}

	options, err := GetOptionsOfOptionGroupContextE(t, ctx, optionGroupName, awsRegion)
	if err != nil {
		return "", err
	}

	for i := range options {
		if aws.ToString(options[i].OptionName) == optionName {
			for _, optionSetting := range options[i].OptionSettings {
				if aws.ToString(optionSetting.Name) == optionSettingName {
					return aws.ToString(optionSetting.Value), nil
				}
			}
		}
	}

	return "", OptionGroupOptionSettingForDBInstanceNotFound{OptionName: optionName, OptionSettingName: optionSettingName, DbInstanceID: dbInstanceID, AwsRegion: awsRegion}
}

// GetOptionSettingForOfRdsInstanceContext gets the value of the option name in the option group specified for the RDS instance in the given region.
// The ctx parameter supports cancellation and timeouts.
func GetOptionSettingForOfRdsInstanceContext(t testing.TestingT, ctx context.Context, optionName string, optionSettingName string, dbInstanceID, awsRegion string) string {
	t.Helper()

	optionValue, err := GetOptionSettingForOfRdsInstanceContextE(t, ctx, optionName, optionSettingName, dbInstanceID, awsRegion)
	require.NoError(t, err)

	return optionValue
}

// GetOptionGroupNameOfRdsInstanceContextE gets the name of the option group associated with the RDS instance.
// The ctx parameter supports cancellation and timeouts.
func GetOptionGroupNameOfRdsInstanceContextE(t testing.TestingT, ctx context.Context, dbInstanceID string, awsRegion string) (string, error) {
	dbInstance, err := GetRdsInstanceDetailsContextE(t, ctx, dbInstanceID, awsRegion)
	if err != nil {
		return "", err
	}

	if len(dbInstance.OptionGroupMemberships) == 0 {
		return "", fmt.Errorf("RDS instance %s in region %s has no option group memberships", dbInstanceID, awsRegion)
	}

	return aws.ToString(dbInstance.OptionGroupMemberships[0].OptionGroupName), nil
}

// GetOptionGroupNameOfRdsInstanceContext gets the name of the option group associated with the RDS instance.
// The ctx parameter supports cancellation and timeouts.
func GetOptionGroupNameOfRdsInstanceContext(t testing.TestingT, ctx context.Context, dbInstanceID string, awsRegion string) string {
	t.Helper()

	dbInstance, err := GetOptionGroupNameOfRdsInstanceContextE(t, ctx, dbInstanceID, awsRegion)
	require.NoError(t, err)

	return dbInstance
}

// GetOptionsOfOptionGroupContextE gets the options of the option group specified.
// The ctx parameter supports cancellation and timeouts.
func GetOptionsOfOptionGroupContextE(t testing.TestingT, ctx context.Context, optionGroupName string, awsRegion string) ([]types.Option, error) {
	rdsClient, err := NewRdsClientContextE(t, ctx, awsRegion)
	if err != nil {
		return []types.Option{}, err
	}

	paginator := rds.NewDescribeOptionGroupsPaginator(rdsClient, &rds.DescribeOptionGroupsInput{
		OptionGroupName: aws.String(optionGroupName),
	})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return []types.Option{}, err
		}

		if len(page.OptionGroupsList) > 0 {
			return page.OptionGroupsList[0].Options, nil
		}
	}

	return []types.Option{}, fmt.Errorf("no option groups found for name %s in region %s", optionGroupName, awsRegion)
}

// GetOptionsOfOptionGroupContext gets the options of the option group specified.
// The ctx parameter supports cancellation and timeouts.
func GetOptionsOfOptionGroupContext(t testing.TestingT, ctx context.Context, optionGroupName string, awsRegion string) []types.Option {
	t.Helper()

	output, err := GetOptionsOfOptionGroupContextE(t, ctx, optionGroupName, awsRegion)
	require.NoError(t, err)

	return output
}

// GetAllParametersOfRdsInstanceContextE gets all the parameters defined in the parameter group for the RDS instance in the given region.
// The ctx parameter supports cancellation and timeouts.
func GetAllParametersOfRdsInstanceContextE(t testing.TestingT, ctx context.Context, dbInstanceID string, awsRegion string) ([]types.Parameter, error) {
	dbInstance, dbInstanceErr := GetRdsInstanceDetailsContextE(t, ctx, dbInstanceID, awsRegion)
	if dbInstanceErr != nil {
		return []types.Parameter{}, dbInstanceErr
	}

	if len(dbInstance.DBParameterGroups) == 0 {
		return []types.Parameter{}, fmt.Errorf("RDS instance %s in region %s has no parameter groups", dbInstanceID, awsRegion)
	}

	parameterGroupName := aws.ToString(dbInstance.DBParameterGroups[0].DBParameterGroupName)

	rdsClient, err := NewRdsClientContextE(t, ctx, awsRegion)
	if err != nil {
		return []types.Parameter{}, err
	}

	input := rds.DescribeDBParametersInput{DBParameterGroupName: aws.String(parameterGroupName)}

	var allParameters []types.Parameter

	for {
		output, err := rdsClient.DescribeDBParameters(ctx, &input)
		if err != nil {
			return []types.Parameter{}, err
		}

		allParameters = append(allParameters, output.Parameters...)

		if output.Marker == nil {
			break
		}

		input.Marker = output.Marker
	}

	return allParameters, nil
}

// GetAllParametersOfRdsInstanceContext gets all the parameters defined in the parameter group for the RDS instance in the given region.
// The ctx parameter supports cancellation and timeouts.
func GetAllParametersOfRdsInstanceContext(t testing.TestingT, ctx context.Context, dbInstanceID string, awsRegion string) []types.Parameter {
	t.Helper()

	parameters, err := GetAllParametersOfRdsInstanceContextE(t, ctx, dbInstanceID, awsRegion)
	require.NoError(t, err)

	return parameters
}

// GetRdsInstanceDetailsContextE gets the details of a single DB instance whose identifier is passed.
// The ctx parameter supports cancellation and timeouts.
func GetRdsInstanceDetailsContextE(t testing.TestingT, ctx context.Context, dbInstanceID string, awsRegion string) (*types.DBInstance, error) {
	rdsClient, err := NewRdsClientContextE(t, ctx, awsRegion)
	if err != nil {
		return nil, err
	}

	input := rds.DescribeDBInstancesInput{DBInstanceIdentifier: aws.String(dbInstanceID)}

	output, err := rdsClient.DescribeDBInstances(ctx, &input)
	if err != nil {
		return nil, err
	}

	if len(output.DBInstances) == 0 {
		return nil, fmt.Errorf("RDS instance %s not found in region %s", dbInstanceID, awsRegion)
	}

	return &output.DBInstances[0], nil
}

// NewRdsClientContextE creates an RDS client.
// The ctx parameter supports cancellation and timeouts.
func NewRdsClientContextE(t testing.TestingT, ctx context.Context, region string) (*rds.Client, error) {
	sess, err := NewAuthenticatedSessionContext(ctx, region)
	if err != nil {
		return nil, err
	}

	return rds.NewFromConfig(*sess), nil
}

// NewRdsClientContext creates an RDS client.
// The ctx parameter supports cancellation and timeouts.
func NewRdsClientContext(t testing.TestingT, ctx context.Context, region string) *rds.Client {
	t.Helper()

	client, err := NewRdsClientContextE(t, ctx, region)
	require.NoError(t, err)

	return client
}

// GetRecommendedRdsInstanceTypeContextE takes in a list of RDS instance types (e.g., "db.t2.micro", "db.t3.micro") and returns the
// first instance type in the list that is available in the given region and for the given database engine type.
// If none of the instances provided are available for your combination of region and database engine, this function will return an error.
// The ctx parameter supports cancellation and timeouts.
func GetRecommendedRdsInstanceTypeContextE(t testing.TestingT, ctx context.Context, region string, engine string, engineVersion string, instanceTypeOptions []string) (string, error) {
	client, err := NewRdsClientContextE(t, ctx, region)
	if err != nil {
		return "", err
	}

	return GetRecommendedRdsInstanceTypeWithClientContextE(t, ctx, client, engine, engineVersion, instanceTypeOptions)
}

// GetRecommendedRdsInstanceTypeContext takes in a list of RDS instance types (e.g., "db.t2.micro", "db.t3.micro") and returns the
// first instance type in the list that is available in the given region and for the given database engine type.
// If none of the instances provided are available for your combination of region and database engine, this function will exit with an error.
// The ctx parameter supports cancellation and timeouts.
func GetRecommendedRdsInstanceTypeContext(t testing.TestingT, ctx context.Context, region string, engine string, engineVersion string, instanceTypeOptions []string) string {
	t.Helper()

	out, err := GetRecommendedRdsInstanceTypeContextE(t, ctx, region, engine, engineVersion, instanceTypeOptions)
	require.NoError(t, err)

	return out
}

// GetRecommendedRdsInstanceTypeWithClientContextE takes in a list of RDS instance types (e.g., "db.t2.micro", "db.t3.micro") and returns the
// first instance type in the list that is available in the given region and for the given database engine type.
// If none of the instances provided are available for your combination of region and database engine, this function will return an error.
// This function expects an authenticated RDS client from the AWS SDK Go library.
// The ctx parameter supports cancellation and timeouts.
func GetRecommendedRdsInstanceTypeWithClientContextE(t testing.TestingT, ctx context.Context, rdsClient *rds.Client, engine string, engineVersion string, instanceTypeOptions []string) (string, error) {
	for _, instanceTypeOption := range instanceTypeOptions {
		instanceTypeExists, err := instanceTypeExistsForEngineAndRegionContextE(ctx, rdsClient, engine, engineVersion, instanceTypeOption)
		if err != nil {
			return "", err
		}

		if instanceTypeExists {
			return instanceTypeOption, nil
		}
	}

	return "", NoRdsInstanceTypeError{InstanceTypeOptions: instanceTypeOptions, DatabaseEngine: engine, DatabaseEngineVersion: engineVersion}
}

// instanceTypeExistsForEngineAndRegionContextE returns a boolean that represents whether the provided instance type (e.g. db.t2.micro) exists for the given region and db engine type.
// This function will return an error if the RDS AWS SDK call fails.
func instanceTypeExistsForEngineAndRegionContextE(ctx context.Context, client *rds.Client, engine string, engineVersion string, instanceType string) (bool, error) {
	paginator := rds.NewDescribeOrderableDBInstanceOptionsPaginator(client, &rds.DescribeOrderableDBInstanceOptionsInput{
		Engine:          aws.String(engine),
		EngineVersion:   aws.String(engineVersion),
		DBInstanceClass: aws.String(instanceType),
	})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return false, err
		}

		if len(page.OrderableDBInstanceOptions) > 0 {
			return true, nil
		}
	}

	return false, nil
}

// GetValidEngineVersionContextE returns a string containing a valid RDS engine version or an error if no valid version is found.
// The ctx parameter supports cancellation and timeouts.
func GetValidEngineVersionContextE(t testing.TestingT, ctx context.Context, region string, engine string, majorVersion string) (string, error) {
	client, err := NewRdsClientContextE(t, ctx, region)
	if err != nil {
		return "", err
	}

	input := rds.DescribeDBEngineVersionsInput{
		Engine:        aws.String(engine),
		EngineVersion: aws.String(majorVersion),
	}

	out, err := client.DescribeDBEngineVersions(ctx, &input)
	if err != nil {
		return "", err
	}

	if len(out.DBEngineVersions) == 0 {
		return "", fmt.Errorf("no engine versions found for engine %s version %s in region %s", engine, majorVersion, region)
	}

	return *out.DBEngineVersions[0].EngineVersion, nil
}

// GetValidEngineVersionContext returns a string containing a valid RDS engine version for the provided region and engine type.
// This function will fail the test if no valid engine is found.
// The ctx parameter supports cancellation and timeouts.
func GetValidEngineVersionContext(t testing.TestingT, ctx context.Context, region string, engine string, majorVersion string) string {
	t.Helper()

	out, err := GetValidEngineVersionContextE(t, ctx, region, engine, majorVersion)
	require.NoError(t, err)

	return out
}

// ParameterForDBInstanceNotFound is an error that occurs when the parameter group specified is not found for the DB instance.
type ParameterForDBInstanceNotFound struct {
	ParameterName string
	DbInstanceID  string //nolint:staticcheck,revive // preserving existing field name
	AwsRegion     string
}

func (err ParameterForDBInstanceNotFound) Error() string {
	return fmt.Sprintf("Could not find a parameter %s in parameter group of database %s in %s", err.ParameterName, err.DbInstanceID, err.AwsRegion)
}

// OptionGroupOptionSettingForDBInstanceNotFound is an error that occurs when the option setting specified is not found in the option group of the DB instance.
type OptionGroupOptionSettingForDBInstanceNotFound struct {
	OptionName        string
	OptionSettingName string
	DbInstanceID      string //nolint:staticcheck,revive // preserving existing field name
	AwsRegion         string
}

func (err OptionGroupOptionSettingForDBInstanceNotFound) Error() string {
	return fmt.Sprintf("Could not find a option setting %s in option name %s of database %s in %s", err.OptionName, err.OptionSettingName, err.DbInstanceID, err.AwsRegion)
}
