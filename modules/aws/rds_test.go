package aws

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRecommendedRdsInstanceTypeHappyPath(t *testing.T) {
	type TestingScenerios struct {
		name               string
		region             string
		databaseEngine     string
		engineMajorVersion string
		instanceTypes      []string
		expected           string
	}

	testingScenerios := []TestingScenerios{
		{
			name:               "US region, mysql, first offering available",
			region:             "us-east-2",
			databaseEngine:     "mysql",
			engineMajorVersion: "8.0",
			instanceTypes:      []string{"db.t4g.micro", "db.t4g.small"},
			expected:           "db.t4g.micro",
		},
		{
			name:               "EU region, postgres, 2nd offering available based on region",
			region:             "eu-north-1",
			databaseEngine:     "postgres",
			engineMajorVersion: "13",
			instanceTypes:      []string{"db.t2.micro", "db.m5.large"},
			expected:           "db.m5.large",
		},
		{
			name:               "US region, oracle-ee, 2nd offering available based on db type",
			region:             "us-west-2",
			databaseEngine:     "oracle-ee",
			engineMajorVersion: "19",
			instanceTypes:      []string{"db.m5d.xlarge", "db.m5.large"},
			expected:           "db.m5d.xlarge",
		},
		{
			name:               "US region, oracle-ee, 2nd offering available based on db engine version",
			region:             "us-west-2",
			databaseEngine:     "oracle-ee",
			engineMajorVersion: "19",
			instanceTypes:      []string{"db.t3.micro", "db.t3.small"},
			expected:           "db.t3.small",
		},
	}

	for _, scenerio := range testingScenerios {
		scenerio := scenerio

		t.Run(scenerio.name, func(t *testing.T) {
			t.Parallel()
			engineVersion := GetValidEngineVersion(t, scenerio.region, scenerio.databaseEngine, scenerio.engineMajorVersion)
			actual, err := GetRecommendedRdsInstanceTypeE(t, scenerio.region, scenerio.databaseEngine, engineVersion, scenerio.instanceTypes)
			assert.NoError(t, err)
			assert.Equal(t, scenerio.expected, actual)
		})
	}
}

func TestGetRecommendedRdsInstanceTypeErrors(t *testing.T) {
	type TestingScenerios struct {
		name                  string
		region                string
		databaseEngine        string
		databaseEngineVersion string
		instanceTypes         []string
	}

	testingScenerios := []TestingScenerios{
		{
			name:                  "All empty",
			region:                "",
			databaseEngine:        "",
			databaseEngineVersion: "",
			instanceTypes:         nil,
		},
		{
			name:                  "No engine, version, or instance type",
			region:                "us-east-2",
			databaseEngine:        "",
			databaseEngineVersion: "",
			instanceTypes:         nil,
		},
		{
			name:                  "No instance types or version",
			region:                "us-east-2",
			databaseEngine:        "mysql",
			databaseEngineVersion: "",
			instanceTypes:         nil,
		},
		{
			name:                  "No engine version",
			region:                "us-east-2",
			databaseEngine:        "mysql",
			databaseEngineVersion: "",
			instanceTypes:         []string{"db.t3.small"},
		},
		{
			name:                  "Invalid instance types",
			region:                "us-east-2",
			databaseEngine:        "mysql",
			databaseEngineVersion: "8.0.32",
			instanceTypes:         []string{"garbage"},
		},
		{
			name:                  "Region has no instance type available",
			region:                "eu-north-1",
			databaseEngine:        "mysql",
			databaseEngineVersion: "8.0.32",
			instanceTypes:         []string{"db.t2.micro"},
		},
		{
			name:                  "No instance type available for engine",
			region:                "us-east-1",
			databaseEngine:        "oracle-ee",
			databaseEngineVersion: "19.0.0.0.ru-2024-04.rur-2024-04.r1",
			instanceTypes:         []string{"db.r5a.large"},
		},
		{
			name:                  "No instance type available for engine version",
			region:                "us-east-1",
			databaseEngine:        "oracle-ee",
			databaseEngineVersion: "19.0.0.0.ru-2021-01.rur-2021-01.r1",
			instanceTypes:         []string{"db.t3.micro"},
		},
	}

	for _, scenerio := range testingScenerios {
		scenerio := scenerio

		t.Run(scenerio.name, func(t *testing.T) {
			t.Parallel()

			_, err := GetRecommendedRdsInstanceTypeE(t, scenerio.region, scenerio.databaseEngine, scenerio.databaseEngineVersion, scenerio.instanceTypes)
			fmt.Println(err)
			assert.EqualError(t, err, NoRdsInstanceTypeError{InstanceTypeOptions: scenerio.instanceTypes, DatabaseEngine: scenerio.databaseEngine, DatabaseEngineVersion: scenerio.databaseEngineVersion}.Error())
		})
	}
}
