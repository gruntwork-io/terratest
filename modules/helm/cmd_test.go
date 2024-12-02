package helm

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/stretchr/testify/assert"
)

func TestPrepareHelmCommand(t *testing.T) {
	t.Parallel()

	options := &Options{
		KubectlOptions: &k8s.KubectlOptions{
			Namespace: "test-namespace",
		},
		EnvVars: map[string]string{"SampleEnv": "test_value"},
		Logger:  logger.Default,
	}
	t.Run("command without additional args", func(t *testing.T) {
		cmd := prepareHelmCommand(t, options, "install")
		assert.Equal(t, "helm", cmd.Command)
		assert.Contains(t, cmd.Args, "install")
		assert.Contains(t, cmd.Args, "--namespace")
		assert.Contains(t, cmd.Args, "test-namespace")
		assert.Equal(t, ".", cmd.WorkingDir)
		assert.Equal(t, options.EnvVars, cmd.Env)
		assert.Equal(t, options.Logger, cmd.Logger)
	})
	t.Run("Command with additional args", func(t *testing.T) {
		cmd := prepareHelmCommand(t, options, "upgrade", "--install", "my-release", "my-chart")
		assert.Equal(t, "helm", cmd.Command)
		assert.Contains(t, cmd.Args, "upgrade")
		assert.Contains(t, cmd.Args, "--install")
		assert.Contains(t, cmd.Args, "my-release")
		assert.Contains(t, cmd.Args, "my-chart")
		assert.Contains(t, cmd.Args, "--namespace")
		assert.Contains(t, cmd.Args, "test-namespace")
	})

	t.Run("Command with namespace in additional args", func(t *testing.T) {
		cmd := prepareHelmCommand(t, options, "install", "--namespace", "custom-namespace")
		assert.Equal(t, "helm", cmd.Command)
		assert.Contains(t, cmd.Args, "install")
		assert.Contains(t, cmd.Args, "--namespace")
		assert.Contains(t, cmd.Args, "custom-namespace")
		assert.NotContains(t, cmd.Args, "test-namespace")
	})
}
