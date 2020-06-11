// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package azure

import (
	"testing"

	"github.com/stretchr/testify/require"
)

/*
The below tests are currently stubbed out, with the expectation that they will throw errors.
If/when CRUD methods are introduced for Azure Virtual Machines, these tests can be extended
(see AWS S3 tests for reference).
*/

func TestGetAppServicePlan(t *testing.T) {
	t.Parallel()

	subId := ""
	resGroupName := ""
	planName := ""

	_, err := getAppServicePlanE(planName, resGroupName, subId)

	require.Error(t, err)
}

func TestGetAppServicePlanClient(t *testing.T) {
	t.Parallel()

	subId := ""

	_, err := getAppServicePlanClient(subId)

	require.Error(t, err)
}