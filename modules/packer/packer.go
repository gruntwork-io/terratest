// Package packer allows to interact with Packer.
package packer

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/shell"
)

// Options are the options for Packer.
type Options struct {
	Template string            // The path to the Packer template
	Vars     map[string]string // The custom vars to pass when running the build command
	Only     string            // If specified, only run the build of this name
	Env      map[string]string // Custom environment variables to set when running Packer
}

// BuildAmi builds the given Packer template and return the generated AMI ID.
func BuildAmi(t *testing.T, options *Options) string {
	amiID, err := BuildAmiE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return amiID
}

// BuildAmiE builds the given Packer template and return the generated AMI ID.
func BuildAmiE(t *testing.T, options *Options) (string, error) {
	logger.Logf(t, "Running Packer to generate AMI for template %s", options.Template)

	cmd := shell.Command{
		Command: "packer",
		Args:    formatPackerArgs(options),
		Env:     options.Env,
	}

	output, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		return "", err
	}

	return extractAMIID(output)
}

// The Packer machine-readable log output should contain an entry of this format:
//
// <timestamp>,<builder>,artifact,<index>,id,<region>:<ami_id>
//
// For example:
//
// 1456332887,amazon-ebs,artifact,0,id,us-east-1:ami-b481b3de
func extractAMIID(packerLogOutput string) (string, error) {
	re := regexp.MustCompile(`.+artifact,\d+?,id,.+?:(.+)`)
	matches := re.FindStringSubmatch(packerLogOutput)

	if len(matches) == 2 {
		return matches[1], nil
	}
	return "", errors.New("Could not find AMI ID pattern in Packer output")
}

// Convert the inputs to a format palatable to packer. The build command should have the format:
//
// packer build [OPTIONS] template
func formatPackerArgs(options *Options) []string {
	args := []string{"build", "-machine-readable"}

	for key, value := range options.Vars {
		args = append(args, "-var", fmt.Sprintf("%s=%s", key, value))
	}

	if options.Only != "" {
		args = append(args, fmt.Sprintf("-only=%s", options.Only))
	}

	return append(args, options.Template)
}
