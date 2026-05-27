package external_test

import (
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/core/v2/collections"
	"github.com/gruntwork-io/terratest/modules/core/v2/random"
	tttesting "github.com/gruntwork-io/terratest/modules/core/v2/testing"
	shell "github.com/gruntwork-io/terratest/modules/shell/v2"
	"github.com/stretchr/testify/assert"
)

func TestExternalConsumer_CorePackages(t *testing.T) {
	t.Parallel()

	id := random.UniqueId()
	assert.Len(t, id, 6)
	assert.NotEmpty(t, strings.TrimSpace(id))
	assert.True(t, collections.ListContains([]string{"a", "b", "c"}, "b"))
}

func TestExternalConsumer_CrossTier(t *testing.T) {
	t.Parallel()

	var tt tttesting.TestingT = t
	out := shell.RunCommandAndGetOutput(tt, shell.Command{
		Command: "echo",
		Args:    []string{"external-consumer-ok"},
	})
	assert.Equal(t, "external-consumer-ok", out)
}
