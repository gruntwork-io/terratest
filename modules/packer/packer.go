// Package packer allows to interact with Packer.
package packer

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/hashicorp/go-version"
)

// Options are the options for Packer.
type Options struct {
	Template                   string            // The path to the Packer template
	Vars                       map[string]string // The custom vars to pass when running the build command
	VarFiles                   []string          // Var file paths to pass Packer using -var-file option
	Only                       string            // If specified, only run the build of this name
	Except                     string            // Runs the build excluding the specified builds and post-processors
	Env                        map[string]string // Custom environment variables to set when running Packer
	RetryableErrors            map[string]string // If packer build fails with one of these (transient) errors, retry. The keys are a regexp to match against the error and the message is what to display to a user if that error is matched.
	MaxRetries                 int               // Maximum number of times to retry errors matching RetryableErrors
	TimeBetweenRetries         time.Duration     // The amount of time to wait between retries
	WorkingDir                 string            // The directory to run packer in
	Logger                     *logger.Logger    // If set, use a non-default logger
	DisableTemporaryPluginPath bool              // If set, do not use a temporary directory for Packer plugins.
}

// BuildArtifacts can take a map of identifierName <-> Options and then parallelize
// the packer builds. Once all the packer builds have completed a map of identifierName <-> generated identifier
// is returned. The identifierName can be anything you want, it is only used so that you can
// know which generated artifact is which.
func BuildArtifacts(t testing.TestingT, artifactNameToOptions map[string]*Options) map[string]string {
	result, err := BuildArtifactsE(t, artifactNameToOptions)

	if err != nil {
		t.Fatalf("Error building artifacts: %s", err.Error())
	}

	return result
}

// BuildArtifactsE can take a map of identifierName <-> Options and then parallelize
// the packer builds. Once all the packer builds have completed a map of identifierName <-> generated identifier
// is returned. If any artifact fails to build, the errors are accumulated and returned
// as a MultiError. The identifierName can be anything you want, it is only used so that you can
// know which generated artifact is which.
func BuildArtifactsE(t testing.TestingT, artifactNameToOptions map[string]*Options) (map[string]string, error) {
	var waitForArtifacts sync.WaitGroup
	waitForArtifacts.Add(len(artifactNameToOptions))

	var artifactNameToArtifactId = map[string]string{}
	var errorsOccurred = new(multierror.Error)

	for artifactName, curOptions := range artifactNameToOptions {
		// The following is necessary to make sure artifactName and curOptions don't
		// get updated due to concurrency within the scope of t.Run(..) below
		artifactName := artifactName
		curOptions := curOptions
		go func() {
			defer waitForArtifacts.Done()
			artifactId, err := BuildArtifactE(t, curOptions)

			if err != nil {
				errorsOccurred = multierror.Append(errorsOccurred, err)
			} else {
				artifactNameToArtifactId[artifactName] = artifactId
			}
		}()
	}

	waitForArtifacts.Wait()

	return artifactNameToArtifactId, errorsOccurred.ErrorOrNil()
}

// BuildArtifact builds the given Packer template and return the generated Artifact ID.
func BuildArtifact(t testing.TestingT, options *Options) string {
	artifactID, err := BuildArtifactE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return artifactID
}

// BuildArtifactE builds the given Packer template and return the generated Artifact ID.
func BuildArtifactE(t testing.TestingT, options *Options) (string, error) {
	options.Logger.Logf(t, "Running Packer to generate a custom artifact for template %s", options.Template)

	// By default, we download packer plugins to a temporary directory rather than use the global plugin path.
	// This prevents race conditions when multiple tests are running in parallel and each of them attempt
	// to download the same plugin at the same time to the global path.
	// Set DisableTemporaryPluginPath to disable this behavior.
	if !options.DisableTemporaryPluginPath {
		// The built-in env variable defining where plugins are downloaded
		const packerPluginPathEnvVar = "PACKER_PLUGIN_PATH"
		options.Logger.Logf(t, "Creating a temporary directory for Packer plugins")
		pluginDir, err := os.MkdirTemp("", "terratest-packer-")
		require.NoError(t, err)
		if len(options.Env) == 0 {
			options.Env = make(map[string]string)
		}
		options.Env[packerPluginPathEnvVar] = pluginDir
		defer os.RemoveAll(pluginDir)
	}

	err := packerInit(t, options)
	if err != nil {
		return "", err
	}

	cmd := shell.Command{
		Command:    "packer",
		Args:       formatPackerArgs(options),
		Env:        options.Env,
		WorkingDir: options.WorkingDir,
	}

	description := fmt.Sprintf("%s %v", cmd.Command, cmd.Args)
	output, err := retry.DoWithRetryableErrorsE(t, description, options.RetryableErrors, options.MaxRetries, options.TimeBetweenRetries, func() (string, error) {
		return shell.RunCommandAndGetOutputE(t, cmd)
	})

	if err != nil {
		return "", err
	}

	return extractArtifactID(output)
}

// BuildAmi builds the given Packer template and return the generated AMI ID.
//
// Deprecated: Use BuildArtifact instead.
func BuildAmi(t testing.TestingT, options *Options) string {
	return BuildArtifact(t, options)
}

// BuildAmiE builds the given Packer template and return the generated AMI ID.
//
// Deprecated: Use BuildArtifactE instead.
func BuildAmiE(t testing.TestingT, options *Options) (string, error) {
	return BuildArtifactE(t, options)
}

// The Packer machine-readable log output should contain an entry of this format:
//
// AWS: <timestamp>,<builder>,artifact,<index>,id,<region>:<image_id>
// GCP: <timestamp>,<builder>,artifact,<index>,id,<image_id>
//
// For example:
//
// 1456332887,amazon-ebs,artifact,0,id,us-east-1:ami-b481b3de
// 1533742764,googlecompute,artifact,0,id,terratest-packer-example-2018-08-08t15-35-19z
func extractArtifactID(packerLogOutput string) (string, error) {
	re := regexp.MustCompile(`.+artifact,\d+?,id,(?:.+?:|)(.+)`)
	matches := re.FindStringSubmatch(packerLogOutput)

	if len(matches) == 2 {
		return matches[1], nil
	}
	return "", errors.New("Could not find Artifact ID pattern in Packer output")
}

// Check if the local version of Packer has init
func hasPackerInit(t testing.TestingT, options *Options) (bool, error) {
	// The init command was introduced in Packer 1.7.0
	const packerInitVersion = "1.7.0"
	minInitVersion, err := version.NewVersion(packerInitVersion)
	if err != nil {
		return false, err
	}

	cmd := shell.Command{
		Command:    "packer",
		Args:       []string{"-version"},
		Env:        options.Env,
		WorkingDir: options.WorkingDir,
	}
	versionCmdOutput, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		return false, err
	}
	localVersion := trimPackerVersion(versionCmdOutput)
	thisVersion, err := version.NewVersion(localVersion)
	if err != nil {
		return false, err
	}

	if thisVersion.LessThan(minInitVersion) {
		return false, nil
	}

	return true, nil
}

// packerInit runs 'packer init' if it is supported by the local packer
func packerInit(t testing.TestingT, options *Options) error {
	hasInit, err := hasPackerInit(t, options)
	if err != nil {
		return err
	}
	if !hasInit {
		options.Logger.Logf(t, "Skipping 'packer init' because it is not present in this version")
		return nil
	}

	extension := filepath.Ext(options.Template)
	if extension != ".hcl" {
		options.Logger.Logf(t, "Skipping 'packer init' because it is only supported for HCL2 templates")
		return nil
	}

	cmd := shell.Command{
		Command:    "packer",
		Args:       []string{"init", options.Template},
		Env:        options.Env,
		WorkingDir: options.WorkingDir,
	}

	description := "Running Packer init"
	_, err = retry.DoWithRetryableErrorsE(t, description, options.RetryableErrors, options.MaxRetries, options.TimeBetweenRetries, func() (string, error) {
		return shell.RunCommandAndGetOutputE(t, cmd)
	})

	if err != nil {
		return err
	}

	return nil
}

// Convert the inputs to a format palatable to packer. The build command should have the format:
//
// packer build [OPTIONS] template
func formatPackerArgs(options *Options) []string {
	args := []string{"build", "-machine-readable"}

	for key, value := range options.Vars {
		args = append(args, "-var", fmt.Sprintf("%s=%s", key, value))
	}

	for _, filePath := range options.VarFiles {
		args = append(args, "-var-file", filePath)
	}

	if options.Only != "" {
		args = append(args, fmt.Sprintf("-only=%s", options.Only))
	}

	if options.Except != "" {
		args = append(args, fmt.Sprintf("-except=%s", options.Except))
	}

	return append(args, options.Template)
}

// From packer 1.10 the -version command output is prefixed with Packer v
func trimPackerVersion(versionCmdOutput string) string {
	re := regexp.MustCompile(`(?:Packer v?|)(\d+\.\d+\.\d+)`)
	matches := re.FindStringSubmatch(versionCmdOutput)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
