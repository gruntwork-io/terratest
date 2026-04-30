package aws

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gruntwork-io/terratest/modules/testing"
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

// GetAddressOfRdsInstance gets the address of the given RDS Instance in the given region.
//
// Deprecated: Use [GetAddressOfRdsInstanceContext] instead.
func GetAddressOfRdsInstance(t testing.TestingT, dbInstanceID string, awsRegion string) string {
	t.Helper()

	return GetAddressOfRdsInstanceContext(t, context.Background(), dbInstanceID, awsRegion)
}

// GetAddressOfRdsInstanceE gets the address of the given RDS Instance in the given region.
//
// Deprecated: Use [GetAddressOfRdsInstanceContextE] instead.
func GetAddressOfRdsInstanceE(t testing.TestingT, dbInstanceID string, awsRegion string) (string, error) {
	return GetAddressOfRdsInstanceContextE(t, context.Background(), dbInstanceID, awsRegion)
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

// GetPortOfRdsInstance gets the port of the given RDS Instance in the given region.
//
// Deprecated: Use [GetPortOfRdsInstanceContext] instead.
func GetPortOfRdsInstance(t testing.TestingT, dbInstanceID string, awsRegion string) int32 {
	t.Helper()

	return GetPortOfRdsInstanceContext(t, context.Background(), dbInstanceID, awsRegion)
}

// GetPortOfRdsInstanceE gets the port of the given RDS Instance in the given region.
//
// Deprecated: Use [GetPortOfRdsInstanceContextE] instead.
func GetPortOfRdsInstanceE(t testing.TestingT, dbInstanceID string, awsRegion string) (int32, error) {
	return GetPortOfRdsInstanceContextE(t, context.Background(), dbInstanceID, awsRegion)
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

// GetWhetherSchemaExistsInRdsMySQLInstance checks whether the specified schema/table name exists in the RDS MySQL instance.
//
// Deprecated: Use [GetWhetherSchemaExistsInRdsMySQLInstanceContext] instead.
func GetWhetherSchemaExistsInRdsMySQLInstance(t testing.TestingT, dbURL string, dbPort int32, dbUsername string, dbPassword string, expectedSchemaName string) bool {
	t.Helper()

	return GetWhetherSchemaExistsInRdsMySQLInstanceContext(t, context.Background(), dbURL, dbPort, dbUsername, dbPassword, expectedSchemaName)
}

// GetWhetherSchemaExistsInRdsMySQLInstanceE checks whether the specified schema/table name exists in the RDS MySQL instance.
//
// Deprecated: Use [GetWhetherSchemaExistsInRdsMySQLInstanceContextE] instead.
func GetWhetherSchemaExistsInRdsMySQLInstanceE(t testing.TestingT, dbURL string, dbPort int32, dbUsername string, dbPassword string, expectedSchemaName string) (bool, error) {
	return GetWhetherSchemaExistsInRdsMySQLInstanceContextE(t, context.Background(), dbURL, dbPort, dbUsername, dbPassword, expectedSchemaName)
}

// GetWhetherSchemaExistsInRdsMySqlInstance checks whether the specified schema/table name exists in the RDS MySQL instance.
//
// Deprecated: Use [GetWhetherSchemaExistsInRdsMySQLInstanceContext] instead.
//
//nolint:staticcheck,revive // preserving deprecated function name
func GetWhetherSchemaExistsInRdsMySqlInstance(t testing.TestingT, dbURL string, dbPort int32, dbUsername string, dbPassword string, expectedSchemaName string) bool {
	t.Helper()

	return GetWhetherSchemaExistsInRdsMySQLInstanceContext(t, context.Background(), dbURL, dbPort, dbUsername, dbPassword, expectedSchemaName)
}

// GetWhetherSchemaExistsInRdsMySqlInstanceE checks whether the specified schema/table name exists in the RDS MySQL instance.
//
// Deprecated: Use [GetWhetherSchemaExistsInRdsMySQLInstanceContextE] instead.
//
//nolint:staticcheck,revive // preserving deprecated function name
func GetWhetherSchemaExistsInRdsMySqlInstanceE(t testing.TestingT, dbURL string, dbPort int32, dbUsername string, dbPassword string, expectedSchemaName string) (bool, error) {
	return GetWhetherSchemaExistsInRdsMySQLInstanceContextE(t, context.Background(), dbURL, dbPort, dbUsername, dbPassword, expectedSchemaName)
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

// GetWhetherSchemaExistsInRdsPostgresInstance checks whether the specified schema/table name exists in the RDS Postgres instance.
//
// Deprecated: Use [GetWhetherSchemaExistsInRdsPostgresInstanceContext] instead.
func GetWhetherSchemaExistsInRdsPostgresInstance(t testing.TestingT, dbURL string, dbPort int32, dbUsername string, dbPassword string, expectedSchemaName string) bool {
	t.Helper()

	return GetWhetherSchemaExistsInRdsPostgresInstanceContext(t, context.Background(), dbURL, dbPort, dbUsername, dbPassword, expectedSchemaName)
}

// GetWhetherSchemaExistsInRdsPostgresInstanceE checks whether the specified schema/table name exists in the RDS Postgres instance.
//
// Deprecated: Use [GetWhetherSchemaExistsInRdsPostgresInstanceContextE] instead.
func GetWhetherSchemaExistsInRdsPostgresInstanceE(t testing.TestingT, dbURL string, dbPort int32, dbUsername string, dbPassword string, expectedSchemaName string) (bool, error) {
	return GetWhetherSchemaExistsInRdsPostgresInstanceContextE(t, context.Background(), dbURL, dbPort, dbUsername, dbPassword, expectedSchemaName)
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

	return "", ParameterForDbInstanceNotFound{ParameterName: parameterName, DbInstanceID: dbInstanceID, AwsRegion: awsRegion} //nolint:staticcheck,revive // preserving existing field name
}

// GetParameterValueForParameterOfRdsInstanceContext gets the value of the parameter name specified for the RDS instance in the given region.
// The ctx parameter supports cancellation and timeouts.
func GetParameterValueForParameterOfRdsInstanceContext(t testing.TestingT, ctx context.Context, parameterName string, dbInstanceID string, awsRegion string) string {
	t.Helper()

	parameterValue, err := GetParameterValueForParameterOfRdsInstanceContextE(t, ctx, parameterName, dbInstanceID, awsRegion)
	require.NoError(t, err)

	return parameterValue
}

// GetParameterValueForParameterOfRdsInstance gets the value of the parameter name specified for the RDS instance in the given region.
//
// Deprecated: Use [GetParameterValueForParameterOfRdsInstanceContext] instead.
func GetParameterValueForParameterOfRdsInstance(t testing.TestingT, parameterName string, dbInstanceID string, awsRegion string) string {
	t.Helper()

	return GetParameterValueForParameterOfRdsInstanceContext(t, context.Background(), parameterName, dbInstanceID, awsRegion)
}

// GetParameterValueForParameterOfRdsInstanceE gets the value of the parameter name specified for the RDS instance in the given region.
//
// Deprecated: Use [GetParameterValueForParameterOfRdsInstanceContextE] instead.
func GetParameterValueForParameterOfRdsInstanceE(t testing.TestingT, parameterName string, dbInstanceID string, awsRegion string) (string, error) {
	return GetParameterValueForParameterOfRdsInstanceContextE(t, context.Background(), parameterName, dbInstanceID, awsRegion)
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

	return "", OptionGroupOptionSettingForDbInstanceNotFound{OptionName: optionName, OptionSettingName: optionSettingName, DbInstanceID: dbInstanceID, AwsRegion: awsRegion} //nolint:staticcheck,revive // preserving existing field name
}

// GetOptionSettingForOfRdsInstanceContext gets the value of the option name in the option group specified for the RDS instance in the given region.
// The ctx parameter supports cancellation and timeouts.
func GetOptionSettingForOfRdsInstanceContext(t testing.TestingT, ctx context.Context, optionName string, optionSettingName string, dbInstanceID, awsRegion string) string {
	t.Helper()

	optionValue, err := GetOptionSettingForOfRdsInstanceContextE(t, ctx, optionName, optionSettingName, dbInstanceID, awsRegion)
	require.NoError(t, err)

	return optionValue
}

// GetOptionSettingForOfRdsInstance gets the value of the option name in the option group specified for the RDS instance in the given region.
//
// Deprecated: Use [GetOptionSettingForOfRdsInstanceContext] instead.
func GetOptionSettingForOfRdsInstance(t testing.TestingT, optionName string, optionSettingName string, dbInstanceID, awsRegion string) string {
	t.Helper()

	return GetOptionSettingForOfRdsInstanceContext(t, context.Background(), optionName, optionSettingName, dbInstanceID, awsRegion)
}

// GetOptionSettingForOfRdsInstanceE gets the value of the option name in the option group specified for the RDS instance in the given region.
//
// Deprecated: Use [GetOptionSettingForOfRdsInstanceContextE] instead.
func GetOptionSettingForOfRdsInstanceE(t testing.TestingT, optionName string, optionSettingName string, dbInstanceID, awsRegion string) (string, error) {
	return GetOptionSettingForOfRdsInstanceContextE(t, context.Background(), optionName, optionSettingName, dbInstanceID, awsRegion)
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

// GetOptionGroupNameOfRdsInstance gets the name of the option group associated with the RDS instance
//
// Deprecated: Use [GetOptionGroupNameOfRdsInstanceContext] instead.
func GetOptionGroupNameOfRdsInstance(t testing.TestingT, dbInstanceID string, awsRegion string) string {
	t.Helper()

	return GetOptionGroupNameOfRdsInstanceContext(t, context.Background(), dbInstanceID, awsRegion)
}

// GetOptionGroupNameOfRdsInstanceE gets the name of the option group associated with the RDS instance
//
// Deprecated: Use [GetOptionGroupNameOfRdsInstanceContextE] instead.
func GetOptionGroupNameOfRdsInstanceE(t testing.TestingT, dbInstanceID string, awsRegion string) (string, error) {
	return GetOptionGroupNameOfRdsInstanceContextE(t, context.Background(), dbInstanceID, awsRegion)
}

// GetOptionsOfOptionGroupContextE gets the options of the option group specified.
// The ctx parameter supports cancellation and timeouts.
func GetOptionsOfOptionGroupContextE(t testing.TestingT, ctx context.Context, optionGroupName string, awsRegion string) ([]types.Option, error) {
	rdsClient, err := NewRdsClientContextE(t, ctx, awsRegion)
	if err != nil {
		return []types.Option{}, err
	}

	input := rds.DescribeOptionGroupsInput{OptionGroupName: aws.String(optionGroupName)}

	output, err := rdsClient.DescribeOptionGroups(ctx, &input)
	if err != nil {
		return []types.Option{}, err
	}

	if len(output.OptionGroupsList) == 0 {
		return []types.Option{}, fmt.Errorf("no option groups found for name %s in region %s", optionGroupName, awsRegion)
	}

	return output.OptionGroupsList[0].Options, nil
}

// GetOptionsOfOptionGroupContext gets the options of the option group specified.
// The ctx parameter supports cancellation and timeouts.
func GetOptionsOfOptionGroupContext(t testing.TestingT, ctx context.Context, optionGroupName string, awsRegion string) []types.Option {
	t.Helper()

	output, err := GetOptionsOfOptionGroupContextE(t, ctx, optionGroupName, awsRegion)
	require.NoError(t, err)

	return output
}

// GetOptionsOfOptionGroup gets the options of the option group specified
//
// Deprecated: Use [GetOptionsOfOptionGroupContext] instead.
func GetOptionsOfOptionGroup(t testing.TestingT, optionGroupName string, awsRegion string) []types.Option {
	t.Helper()

	return GetOptionsOfOptionGroupContext(t, context.Background(), optionGroupName, awsRegion)
}

// GetOptionsOfOptionGroupE gets the options of the option group specified
//
// Deprecated: Use [GetOptionsOfOptionGroupContextE] instead.
func GetOptionsOfOptionGroupE(t testing.TestingT, optionGroupName string, awsRegion string) ([]types.Option, error) {
	return GetOptionsOfOptionGroupContextE(t, context.Background(), optionGroupName, awsRegion)
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

// GetAllParametersOfRdsInstance gets all the parameters defined in the parameter group for the RDS instance in the given region.
//
// Deprecated: Use [GetAllParametersOfRdsInstanceContext] instead.
func GetAllParametersOfRdsInstance(t testing.TestingT, dbInstanceID string, awsRegion string) []types.Parameter {
	t.Helper()

	return GetAllParametersOfRdsInstanceContext(t, context.Background(), dbInstanceID, awsRegion)
}

// GetAllParametersOfRdsInstanceE gets all the parameters defined in the parameter group for the RDS instance in the given region.
//
// Deprecated: Use [GetAllParametersOfRdsInstanceContextE] instead.
func GetAllParametersOfRdsInstanceE(t testing.TestingT, dbInstanceID string, awsRegion string) ([]types.Parameter, error) {
	return GetAllParametersOfRdsInstanceContextE(t, context.Background(), dbInstanceID, awsRegion)
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

// GetRdsInstanceDetailsE gets the details of a single DB instance whose identifier is passed.
//
// Deprecated: Use [GetRdsInstanceDetailsContextE] instead.
func GetRdsInstanceDetailsE(t testing.TestingT, dbInstanceID string, awsRegion string) (*types.DBInstance, error) {
	return GetRdsInstanceDetailsContextE(t, context.Background(), dbInstanceID, awsRegion)
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

// NewRdsClient creates an RDS client.
//
// Deprecated: Use [NewRdsClientContext] instead.
func NewRdsClient(t testing.TestingT, region string) *rds.Client {
	t.Helper()

	return NewRdsClientContext(t, context.Background(), region)
}

// NewRdsClientE creates an RDS client.
//
// Deprecated: Use [NewRdsClientContextE] instead.
func NewRdsClientE(t testing.TestingT, region string) (*rds.Client, error) {
	return NewRdsClientContextE(t, context.Background(), region)
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

// GetRecommendedRdsInstanceType takes in a list of RDS instance types (e.g., "db.t2.micro", "db.t3.micro") and returns the
// first instance type in the list that is available in the given region and for the given database engine type.
// If none of the instances provided are available for your combination of region and database engine, this function will exit with an error.
//
// Deprecated: Use [GetRecommendedRdsInstanceTypeContext] instead.
func GetRecommendedRdsInstanceType(t testing.TestingT, region string, engine string, engineVersion string, instanceTypeOptions []string) string {
	t.Helper()

	return GetRecommendedRdsInstanceTypeContext(t, context.Background(), region, engine, engineVersion, instanceTypeOptions)
}

// GetRecommendedRdsInstanceTypeE takes in a list of RDS instance types (e.g., "db.t2.micro", "db.t3.micro") and returns the
// first instance type in the list that is available in the given region and for the given database engine type.
// If none of the instances provided are available for your combination of region and database engine, this function will return an error.
//
// Deprecated: Use [GetRecommendedRdsInstanceTypeContextE] instead.
func GetRecommendedRdsInstanceTypeE(t testing.TestingT, region string, engine string, engineVersion string, instanceTypeOptions []string) (string, error) {
	return GetRecommendedRdsInstanceTypeContextE(t, context.Background(), region, engine, engineVersion, instanceTypeOptions)
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

// GetRecommendedRdsInstanceTypeWithClientE takes in a list of RDS instance types (e.g., "db.t2.micro", "db.t3.micro") and returns the
// first instance type in the list that is available in the given region and for the given database engine type.
// If none of the instances provided are available for your combination of region and database engine, this function will return an error.
// This function expects an authenticated RDS client from the AWS SDK Go library.
//
// Deprecated: Use [GetRecommendedRdsInstanceTypeWithClientContextE] instead.
func GetRecommendedRdsInstanceTypeWithClientE(t testing.TestingT, rdsClient *rds.Client, engine string, engineVersion string, instanceTypeOptions []string) (string, error) {
	return GetRecommendedRdsInstanceTypeWithClientContextE(t, context.Background(), rdsClient, engine, engineVersion, instanceTypeOptions)
}

// instanceTypeExistsForEngineAndRegionContextE returns a boolean that represents whether the provided instance type (e.g. db.t2.micro) exists for the given region and db engine type.
// This function will return an error if the RDS AWS SDK call fails.
func instanceTypeExistsForEngineAndRegionContextE(ctx context.Context, client *rds.Client, engine string, engineVersion string, instanceType string) (bool, error) {
	input := rds.DescribeOrderableDBInstanceOptionsInput{
		Engine:          aws.String(engine),
		EngineVersion:   aws.String(engineVersion),
		DBInstanceClass: aws.String(instanceType),
	}

	out, err := client.DescribeOrderableDBInstanceOptions(ctx, &input)
	if err != nil {
		return false, err
	}

	if len(out.OrderableDBInstanceOptions) > 0 {
		return true, nil
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

// GetValidEngineVersion returns a string containing a valid RDS engine version for the provided region and engine type.
// This function will fail the test if no valid engine is found.
//
// Deprecated: Use [GetValidEngineVersionContext] instead.
func GetValidEngineVersion(t testing.TestingT, region string, engine string, majorVersion string) string {
	t.Helper()

	return GetValidEngineVersionContext(t, context.Background(), region, engine, majorVersion)
}

// GetValidEngineVersionE returns a string containing a valid RDS engine version or an error if no valid version is found.
//
// Deprecated: Use [GetValidEngineVersionContextE] instead.
func GetValidEngineVersionE(t testing.TestingT, region string, engine string, majorVersion string) (string, error) {
	return GetValidEngineVersionContextE(t, context.Background(), region, engine, majorVersion)
}

// ParameterForDbInstanceNotFound is an error that occurs when the parameter group specified is not found for the DB instance.
//
// Deprecated: Use [ParameterForDBInstanceNotFound] instead.
//
//nolint:staticcheck,revive // preserving existing type name
type ParameterForDbInstanceNotFound struct {
	ParameterName string
	DbInstanceID  string //nolint:staticcheck,revive // preserving existing field name
	AwsRegion     string
}

// ParameterForDBInstanceNotFound is an alias for ParameterForDbInstanceNotFound with the correct Go naming convention.
type ParameterForDBInstanceNotFound = ParameterForDbInstanceNotFound

func (err ParameterForDbInstanceNotFound) Error() string {
	return fmt.Sprintf("Could not find a parameter %s in parameter group of database %s in %s", err.ParameterName, err.DbInstanceID, err.AwsRegion)
}

// OptionGroupOptionSettingForDbInstanceNotFound is an error that occurs when the option setting specified is not found in the option group of the DB instance.
//
// Deprecated: Use [OptionGroupOptionSettingForDBInstanceNotFound] instead.
//
//nolint:staticcheck,revive // preserving existing type name
type OptionGroupOptionSettingForDbInstanceNotFound struct {
	OptionName        string
	OptionSettingName string
	DbInstanceID      string //nolint:staticcheck,revive // preserving existing field name
	AwsRegion         string
}

// OptionGroupOptionSettingForDBInstanceNotFound is an alias for OptionGroupOptionSettingForDbInstanceNotFound with the correct Go naming convention.
type OptionGroupOptionSettingForDBInstanceNotFound = OptionGroupOptionSettingForDbInstanceNotFound

func (err OptionGroupOptionSettingForDbInstanceNotFound) Error() string {
	return fmt.Sprintf("Could not find a option setting %s in option name %s of database %s in %s", err.OptionName, err.OptionSettingName, err.DbInstanceID, err.AwsRegion)
}
