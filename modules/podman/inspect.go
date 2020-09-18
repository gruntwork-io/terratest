package podman

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/stretchr/testify/require"
)

// ContainerInspect defines the output of the Inspect method, with the options returned by 'podman inspect'
// converted into a more friendly and testable interface
type ContainerInspect struct {
	// ID of the inspected container
	ID string

	// Name of the inspected container
	Name string

	// time.Time that the container was created
	Created time.Time

	// String representing the container's status
	Status string

	// Whether the container is currently running or not
	Running bool

	// Container's exit code
	ExitCode uint8

	// String with the container's error message, if there is any
	Error string

	// Ports exposed by the container
	Ports []Port

	// Volume bindings made to the container
	Binds []VolumeBind

	// Health check
	Health HealthCheck
}

// Port represents a single port mapping exported by the container
type Port struct {
	HostPort      uint16
	ContainerPort uint16
	Protocol      string
}

// VolumeBind represents a single volume binding made to the container
type VolumeBind struct {
	Source      string
	Destination string
}

// HealthCheck represents the current health history of the container
type HealthCheck struct {
	// Health check status
	Status string

	// Current count of failing health checks
	FailingStreak uint8

	// Log of failures
	Log []HealthLog
}

// HealthLog represents the output of a single Health check of the container
type HealthLog struct {
	// Start time of health check
	Start string

	// End time of health check
	End string

	// Exit code of health check
	ExitCode uint8

	// Output of health check
	Output string
}

// inspectOutput defines options that will be returned by 'podman inspect', in JSON format.
// Not all options are included here, only the ones that we might need
type inspectOutput struct {
	Id      string
	Created string
	Name    string
	State   struct {
		Health   HealthCheck
		Status   string
		Running  bool
		ExitCode uint8
		Error    string
	}
	NetworkSettings struct {
		Ports []Port
	}
	HostConfig struct {
		Binds []string
	}
}

// Inspect runs the 'podman inspect {container id}' command and returns a ContainerInspect
// struct, converted from the output JSON, along with any errors
func Inspect(t *testing.T, id string) *ContainerInspect {
	out, err := InspectE(t, id)
	require.NoError(t, err)

	return out
}

// InspectE runs the 'podman inspect {container id}' command and returns a ContainerInspect
// struct, converted from the output JSON, along with any errors
func InspectE(t *testing.T, id string) (*ContainerInspect, error) {
	cmd := shell.Command{
		Command: "podman",
		Args:    []string{"container", "inspect", id},
		// inspect is a short-running command, don't print the output.
		Logger: logger.Discard,
	}

	out, err := shell.RunCommandAndGetStdOutE(t, cmd)
	if err != nil {
		return nil, err
	}

	var containers []inspectOutput
	err = json.Unmarshal([]byte(out), &containers)
	if err != nil {
		return nil, err
	}

	if len(containers) == 0 {
		return nil, fmt.Errorf("no container found with ID %s", id)
	}

	container := containers[0]

	return transformContainer(t, container)
}

// transformContainerPorts converts 'podman inspect' output JSON into a more friendly and testable format
func transformContainer(t *testing.T, container inspectOutput) (*ContainerInspect, error) {
	name := strings.TrimLeft(container.Name, "/")

	volumes := transformContainerVolumes(container)

	created, err := time.Parse(time.RFC3339Nano, container.Created)
	if err != nil {
		return nil, err
	}

	inspect := ContainerInspect{
		ID:       container.Id,
		Name:     name,
		Created:  created,
		Status:   container.State.Status,
		Running:  container.State.Running,
		ExitCode: container.State.ExitCode,
		Error:    container.State.Error,
		Ports:    container.NetworkSettings.Ports,
		Binds:    volumes,
		Health: HealthCheck{
			Status:        container.State.Health.Status,
			FailingStreak: container.State.Health.FailingStreak,
			Log:           container.State.Health.Log,
		},
	}

	return &inspect, nil
}

// transformContainerVolumes converts Podman's volume bindings from the
// format "/foo/bar:/foo/baz" into a more testable one
func transformContainerVolumes(container inspectOutput) []VolumeBind {
	binds := container.HostConfig.Binds
	volumes := make([]VolumeBind, 0, len(binds))

	for _, bind := range binds {
		var source, dest string

		split := strings.Split(bind, ":")

		// Considering it as an unbound volume
		dest = split[0]

		if len(split) == 2 {
			source = split[0]
			dest = split[1]
		}

		volumes = append(volumes, VolumeBind{
			Source:      source,
			Destination: dest,
		})
	}

	return volumes
}
