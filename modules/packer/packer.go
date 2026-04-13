// Package packer allows to interact with Packer.
package packer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/require"
)

// ErrArtifactIDNotFound is returned when the Packer output does not contain an artifact ID.
var ErrArtifactIDNotFound = errors.New("could not find artifact ID pattern in packer output")

// BuildNameNotFoundError is returned when the specified build name is not found in the manifest file.
type BuildNameNotFoundError struct {
	BuildName    string
	ManifestPath string
}

// Error implements the error interface for BuildNameNotFoundError.
func (e *BuildNameNotFoundError) Error() string {
	return fmt.Sprintf("build name %s not found in manifest file %s", e.BuildName, e.ManifestPath)
}

// Options are the options for Packer.
type Options struct {
	Vars                       map[string]string // The custom vars to pass when running the build command
	Env                        map[string]string // Custom environment variables to set when running Packer
	RetryableErrors            map[string]string // If packer build fails with one of these (transient) errors, retry. The keys are a regexp to match against the error and the message is what to display to a user if that error is matched.
	Logger                     *logger.Logger    // If set, use a non-default logger
	Template                   string            // The path to the Packer template
	Only                       string            // If specified, only run the build of this name
	Except                     string            // Runs the build excluding the specified builds and post-processors
	WorkingDir                 string            // The directory to run packer in
	VarFiles                   []string          // Var file paths to pass Packer using -var-file option
	TimeBetweenRetries         time.Duration     // The amount of time to wait between retries
	MaxRetries                 int               // Maximum number of times to retry errors matching RetryableErrors
	DisableTemporaryPluginPath bool              // If set, do not use a temporary directory for Packer plugins.
}

// BuildArtifactsContextE can take a map of identifierName <-> Options and then parallelize
// the packer builds. Once all the packer builds have completed a map of identifierName <-> generated identifier
// is returned. If any artifact fails to build, the errors are accumulated and returned
// as a MultiError. The identifierName can be anything you want, it is only used so that you can
// know which generated artifact is which.
// The ctx parameter supports cancellation and timeouts.
func BuildArtifactsContextE(t testing.TestingT, ctx context.Context, artifactNameToOptions map[string]*Options) (map[string]string, error) {
	var waitForArtifacts sync.WaitGroup

	waitForArtifacts.Add(len(artifactNameToOptions))

	artifactNameToArtifactID := map[string]string{}
	errorsOccurred := new(multierror.Error)

	var mu sync.Mutex

	for artifactName, curOptions := range artifactNameToOptions {
		go func() {
			defer waitForArtifacts.Done()

			artifactID, err := BuildArtifactContextE(t, ctx, curOptions)

			mu.Lock()
			if err != nil {
				errorsOccurred = multierror.Append(errorsOccurred, err)
			} else {
				artifactNameToArtifactID[artifactName] = artifactID
			}
			mu.Unlock()
		}()
	}

	waitForArtifacts.Wait()

	return artifactNameToArtifactID, errorsOccurred.ErrorOrNil()
}

// BuildArtifactsContext can take a map of identifierName <-> Options and then parallelize
// the packer builds. Once all the packer builds have completed a map of identifierName <-> generated identifier
// is returned. The identifierName can be anything you want, it is only used so that you can
// know which generated artifact is which.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func BuildArtifactsContext(t testing.TestingT, ctx context.Context, artifactNameToOptions map[string]*Options) map[string]string {
	t.Helper()

	result, err := BuildArtifactsContextE(t, ctx, artifactNameToOptions)
	require.NoError(t, err)

	return result
}

// BuildArtifacts can take a map of identifierName <-> Options and then parallelize
// the packer builds. Once all the packer builds have completed a map of identifierName <-> generated identifier
// is returned. The identifierName can be anything you want, it is only used so that you can
// know which generated artifact is which.
//
// Deprecated: Use [BuildArtifactsContext] instead.
func BuildArtifacts(t testing.TestingT, artifactNameToOptions map[string]*Options) map[string]string {
	t.Helper()

	return BuildArtifactsContext(t, context.Background(), artifactNameToOptions)
}

// BuildArtifactsE can take a map of identifierName <-> Options and then parallelize
// the packer builds. Once all the packer builds have completed a map of identifierName <-> generated identifier
// is returned. If any artifact fails to build, the errors are accumulated and returned
// as a MultiError. The identifierName can be anything you want, it is only used so that you can
// know which generated artifact is which.
//
// Deprecated: Use [BuildArtifactsContextE] instead.
func BuildArtifactsE(t testing.TestingT, artifactNameToOptions map[string]*Options) (map[string]string, error) {
	return BuildArtifactsContextE(t, context.Background(), artifactNameToOptions)
}

// BuildArtifactContextE builds the given Packer template and return the generated Artifact ID.
// The ctx parameter supports cancellation and timeouts.
func BuildArtifactContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
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
		if err != nil {
			return "", fmt.Errorf("creating temporary plugin directory: %w", err)
		}

		if len(options.Env) == 0 {
			options.Env = make(map[string]string)
		}

		options.Env[packerPluginPathEnvVar] = pluginDir

		defer func() { _ = os.RemoveAll(pluginDir) }()
	}

	err := packerInit(t, ctx, options)
	if err != nil {
		return "", err
	}

	cmd := shell.Command{
		Command:    "packer",
		Args:       FormatPackerArgs(options),
		Env:        options.Env,
		WorkingDir: options.WorkingDir,
	}

	description := fmt.Sprintf("%s %v", cmd.Command, cmd.Args)

	output, err := retry.DoWithRetryableErrorsContextE(t, ctx, description, options.RetryableErrors, options.MaxRetries, options.TimeBetweenRetries, func() (string, error) {
		return shell.RunCommandContextAndGetOutputE(t, ctx, &cmd)
	})
	if err != nil {
		return "", err
	}

	return ExtractArtifactID(output)
}

// BuildArtifactContext builds the given Packer template and return the generated Artifact ID.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func BuildArtifactContext(t testing.TestingT, ctx context.Context, options *Options) string {
	t.Helper()

	result, err := BuildArtifactContextE(t, ctx, options)
	require.NoError(t, err)

	return result
}

// BuildArtifact builds the given Packer template and return the generated Artifact ID.
//
// Deprecated: Use [BuildArtifactContext] instead.
func BuildArtifact(t testing.TestingT, options *Options) string {
	t.Helper()

	return BuildArtifactContext(t, context.Background(), options)
}

// BuildArtifactE builds the given Packer template and return the generated Artifact ID.
//
// Deprecated: Use [BuildArtifactContextE] instead.
func BuildArtifactE(t testing.TestingT, options *Options) (string, error) {
	return BuildArtifactContextE(t, context.Background(), options)
}

// BuildAmi builds the given Packer template and return the generated AMI ID.
//
// Deprecated: Use [BuildArtifactContext] instead.
func BuildAmi(t testing.TestingT, options *Options) string {
	t.Helper()

	return BuildArtifactContext(t, context.Background(), options)
}

// BuildAmiE builds the given Packer template and return the generated AMI ID.
//
// Deprecated: Use [BuildArtifactContextE] instead.
func BuildAmiE(t testing.TestingT, options *Options) (string, error) {
	return BuildArtifactContextE(t, context.Background(), options)
}

// artifactIDMatchLen is the expected number of submatches (full match + capture group)
// from the artifact ID regex.
const artifactIDMatchLen = 2

// ExtractArtifactID extracts the artifact ID from Packer machine-readable log output.
//
// The Packer machine-readable log output should contain an entry of this format:
//
// AWS: <timestamp>,<builder>,artifact,<index>,id,<region>:<image_id>
// GCP: <timestamp>,<builder>,artifact,<index>,id,<image_id>
//
// For example:
//
// 1456332887,amazon-ebs,artifact,0,id,us-east-1:ami-b481b3de
// 1533742764,googlecompute,artifact,0,id,terratest-packer-example-2018-08-08t15-35-19z
func ExtractArtifactID(packerLogOutput string) (string, error) {
	re := regexp.MustCompile(`.+artifact,\d+?,id,(?:.+?:|)(.+)`)
	matches := re.FindStringSubmatch(packerLogOutput)

	if len(matches) == artifactIDMatchLen {
		return matches[1], nil
	}

	return "", ErrArtifactIDNotFound
}

// hasPackerInit checks if the local version of Packer supports the init command.
func hasPackerInit(t testing.TestingT, ctx context.Context, options *Options) (bool, error) {
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

	versionCmdOutput, err := shell.RunCommandContextAndGetOutputE(t, ctx, &cmd)
	if err != nil {
		return false, err
	}

	localVersion := TrimPackerVersion(versionCmdOutput)

	thisVersion, err := version.NewVersion(localVersion)
	if err != nil {
		return false, err
	}

	if thisVersion.LessThan(minInitVersion) {
		return false, nil
	}

	return true, nil
}

// packerInit runs 'packer init' if it is supported by the local packer.
func packerInit(t testing.TestingT, ctx context.Context, options *Options) error {
	hasInit, err := hasPackerInit(t, ctx, options)
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

	_, err = retry.DoWithRetryableErrorsContextE(t, ctx, description, options.RetryableErrors, options.MaxRetries, options.TimeBetweenRetries, func() (string, error) {
		return shell.RunCommandContextAndGetOutputE(t, ctx, &cmd)
	})

	return err
}

// FormatPackerArgs converts the inputs to a format palatable to packer. The build command should have the format:
//
// packer build [OPTIONS] template
func FormatPackerArgs(options *Options) []string {
	args := []string{"build", "-machine-readable"}

	for key, value := range options.Vars {
		args = append(args, "-var", key+"="+value)
	}

	for _, filePath := range options.VarFiles {
		args = append(args, "-var-file", filePath)
	}

	if options.Only != "" {
		args = append(args, "-only="+options.Only)
	}

	if options.Except != "" {
		args = append(args, "-except="+options.Except)
	}

	return append(args, options.Template)
}

// TrimPackerVersion extracts the version number from packer version output.
// From packer 1.10 the -version command output is prefixed with "Packer v".
func TrimPackerVersion(versionCmdOutput string) string {
	re := regexp.MustCompile(`(?:Packer v?|)(\d+\.\d+\.\d+)`)
	matches := re.FindStringSubmatch(versionCmdOutput)

	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}

type packerManifest struct {
	LastRunUUID string                `json:"last_run_uuid"`
	Builds      []packerManifestBuild `json:"builds"`
}

type packerManifestBuild struct {
	Name          string                    `json:"name"`
	BuilderType   string                    `json:"builder_type"`
	ArtifactID    string                    `json:"artifact_id"`
	PackerRunUUID string                    `json:"packer_run_uuid"`
	CustomData    map[string]interface{}    `json:"custom_data"`
	Files         []packerManifestBuildFile `json:"files"`
	BuildTime     int64                     `json:"build_time"`
}

type packerManifestBuildFile struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

// GetArtifactIDFromManifestBuildNameContextE returns the artifact id from a build name contained in the manifest file.
// See https://developer.hashicorp.com/packer/docs/post-processors/manifest for more info.
// The ctx parameter is accepted for API consistency with other Context functions.
func GetArtifactIDFromManifestBuildNameContextE(t testing.TestingT, ctx context.Context, manifestPath string, buildName string) (string, error) {
	b, err := os.ReadFile(manifestPath)
	if err != nil {
		return "", fmt.Errorf("reading manifest file: %w", err)
	}

	var manifest packerManifest

	if err = json.Unmarshal(b, &manifest); err != nil {
		return "", fmt.Errorf("unmarshalling manifest file: %w", err)
	}

	for _, build := range manifest.Builds {
		if build.Name == buildName {
			return build.ArtifactID, nil
		}
	}

	return "", &BuildNameNotFoundError{BuildName: buildName, ManifestPath: manifestPath}
}

// GetArtifactIDFromManifestBuildNameContext returns the artifact id from a build name contained in the manifest file.
// See https://developer.hashicorp.com/packer/docs/post-processors/manifest for more info.
// This function will fail the test if there is an error.
// The ctx parameter is accepted for API consistency with other Context functions.
func GetArtifactIDFromManifestBuildNameContext(t testing.TestingT, ctx context.Context, manifestPath string, buildName string) string {
	t.Helper()

	result, err := GetArtifactIDFromManifestBuildNameContextE(t, ctx, manifestPath, buildName)
	require.NoError(t, err)

	return result
}

// GetArtifactIDFromManifestBuildName returns the artifact id from a build name contained in the manifest file.
// See https://developer.hashicorp.com/packer/docs/post-processors/manifest for more info.
// If the build name is not found, it will fail the test.
//
// Deprecated: Use [GetArtifactIDFromManifestBuildNameContext] instead.
func GetArtifactIDFromManifestBuildName(t testing.TestingT, manifestPath string, buildName string) string {
	t.Helper()

	return GetArtifactIDFromManifestBuildNameContext(t, context.Background(), manifestPath, buildName)
}

// GetArtifactIDFromManifestBuildNameE returns the artifact id from a build name contained in the manifest file.
// See https://developer.hashicorp.com/packer/docs/post-processors/manifest for more info.
//
// Deprecated: Use [GetArtifactIDFromManifestBuildNameContextE] instead.
func GetArtifactIDFromManifestBuildNameE(t testing.TestingT, manifestPath string, buildName string) (string, error) {
	return GetArtifactIDFromManifestBuildNameContextE(t, context.Background(), manifestPath, buildName)
}
