package terraform

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestNewPlanInfoFlattensResources(t *testing.T) {
	raw, err := ioutil.ReadFile("testdata/plan_output/root_and_child_modules.json")
	info, err := NewPlanInfo(string(raw))
	assert.Nil(t, err)
	assert.Equal(t, len(info.AllResources), 2)
	assert.Equal(t, "some-name", info.AllResources[0].Name)
	assert.Equal(t, "some-other-name", info.AllResources[1].Name)
}
