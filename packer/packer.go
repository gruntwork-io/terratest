package packer

import (
	"errors"
	"github.com/gruntwork-io/terratest/shell"
	"log"
	"regexp"
	"strings"
)

type PackerOptions struct {
	Template string            // The path to the Packer template
	Vars     map[string]string // The custom vars to pass when running the build command
	Only     string            // If specified, only run the build of this name
	Env      map[string]string // Custom environment variables to set when running Packer
}

// Build the given Packer template and return the generated AMI ID
func BuildAmi(options PackerOptions, logger *log.Logger) (string, error) {
	logger.Printf("Running Packer to generate AMI for template %s", options.Template)

	cmd := shell.Command{
		Command: "packer",
		Args:    formatPackerArgs(options),
		Env:     options.Env,
	}

	output, err := shell.RunCommandAndGetOutput(cmd, logger)
	if err != nil {
		return "", err
	}

	return extractAmiId(output)
}

// The Packer machine-readable log output should contain an entry of this format:
//
// <timestamp>,<builder>,artifact,<index>,id,<region>:<ami_id>
//
// For example:
//
// 1456332887,amazon-ebs,artifact,0,id,us-east-1:ami-b481b3de
func extractAmiId(packerLogOutput string) (string, error) {
	re := regexp.MustCompile(".+artifact,\\d+?,id,.+?:(.+)")
	matches := re.FindAllStringSubmatch(packerLogOutput, -1)

	for _, curMatch := range matches {
		if len(curMatch) == 2 && !strings.Contains(curMatch[0], "sha256") {
			return curMatch[1], nil
		}
	}
	return "", errors.New("Could not find AMI ID pattern in Packer output")
}

// Convert the inputs to a format palatable to packer. The build command should have the format:
//
// packer build [OPTIONS] template
func formatPackerArgs(options PackerOptions) []string {
	args := []string{"build", "-machine-readable"}

	for key, value := range options.Vars {
		args = append(args, "-var", key+"="+value)
	}

	if options.Only != "" {
		args = append(args, "-only="+options.Only)
	}

	return append(args, options.Template)
}
