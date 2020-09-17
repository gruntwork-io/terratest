// Package packer allows to interact with Packer.
package packer

import (
	"errors"
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/gruntwork-io/terratest/modules/retry"

	"github.com/gruntwork-io/terratest/modules/customerrors"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// Options are the options for Packer.
type Options struct {
	Template           string            // The path to the Packer template
	Vars               map[string]string // The custom vars to pass when running the build command
	VarFiles           []string          // Var file paths to pass Packer using -var-file option
	Only               string            // If specified, only run the build of this name
	Env                map[string]string // Custom environment variables to set when running Packer
	RetryableErrors    map[string]string // If packer build fails with one of these (transient) errors, retry. The keys are a regexp to match against the error and the message is what to display to a user if that error is matched.
	MaxRetries         int               // Maximum number of times to retry errors matching RetryableErrors
	TimeBetweenRetries time.Duration     // The amount of time to wait between retries
	WorkingDir         string            // The directory to run packer in
	Logger             *logger.Logger    // If set, use a non-default logger
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
	errorsOccurred := []error{}

	for artifactName, curOptions := range artifactNameToOptions {
		// The following is necessary to make sure artifactName and curOptions don't
		// get updated due to concurrency within the scope of t.Run(..) below
		artifactName := artifactName
		curOptions := curOptions
		go func() {
			defer waitForArtifacts.Done()
			artifactId, err := BuildArtifactE(t, curOptions)

			if err != nil {
				errorsOccurred = append(errorsOccurred, err)
			} else {
				artifactNameToArtifactId[artifactName] = artifactId
			}
		}()
	}

	waitForArtifacts.Wait()

	return artifactNameToArtifactId, customerrors.NewMultiError(errorsOccurred...)
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
//
func extractArtifactID(packerLogOutput string) (string, error) {
	re := regexp.MustCompile(`.+artifact,\d+?,id,(?:.+?:|)(.+)`)
	matches := re.FindStringSubmatch(packerLogOutput)

	if len(matches) == 2 {
		return matches[1], nil
	}
	return "", errors.New("Could not find Artifact ID pattern in Packer output")
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

	return append(args, options.Template)
}
