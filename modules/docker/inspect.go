package docker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// ErrNoContainerFound is returned when no container is found with the given ID.
var ErrNoContainerFound = errors.New("no container found")

// ContainerInspect defines the output of the Inspect method, with the options returned by 'docker inspect'
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

	// String with the container's error message, if there is any
	Error string

	// Ports exposed by the container
	Ports []Port

	// Volume bindings made to the container
	Binds []VolumeBind

	// Health check
	Health HealthCheck

	// Container's exit code
	ExitCode uint8

	// Whether the container is currently running or not
	Running bool
}

// Port represents a single port mapping exported by the container
type Port struct {
	Protocol      string
	HostPort      uint16
	ContainerPort uint16
}

// VolumeBind represents a single volume binding made to the container
type VolumeBind struct {
	Source      string
	Destination string
}

// HealthCheck represents the current health history of the container
type HealthCheck struct {
	// Health check status
	Status string `json:"Status"`

	// Log of failures
	Log []HealthLog `json:"Log"`

	// Current count of failing health checks
	FailingStreak uint8 `json:"FailingStreak"`
}

// HealthLog represents the output of a single Health check of the container
type HealthLog struct {
	// Start time of health check
	Start string `json:"Start"`

	// End time of health check
	End string `json:"End"`

	// Output of health check
	Output string `json:"Output"`

	// Exit code of health check
	ExitCode uint8 `json:"ExitCode"`
}

// volumeBindParts is the expected number of colon-separated parts in a host
// volume bind string (source:destination).
const volumeBindParts = 2

// inspectOutput defines options that will be returned by 'docker inspect', in JSON format.
// Not all options are included here, only the ones that we might need
type inspectOutput struct {
	NetworkSettings struct {
		Ports map[string][]struct {
			HostIP   string `json:"HostIp"`
			HostPort string `json:"HostPort"`
		} `json:"Ports"`
	} `json:"NetworkSettings"`
	HostConfig struct {
		Binds []string `json:"Binds"`
	} `json:"HostConfig"`
	ID      string `json:"Id"`
	Created string `json:"Created"`
	Name    string `json:"Name"`
	State   struct {
		Status   string      `json:"Status"`
		Error    string      `json:"Error"`
		Health   HealthCheck `json:"Health"`
		ExitCode uint8       `json:"ExitCode"`
		Running  bool        `json:"Running"`
	} `json:"State"`
}

// Inspect runs the 'docker inspect {container id}' command and returns a ContainerInspect
// struct, converted from the output JSON, along with any errors
//
// Deprecated: Use [InspectContext] instead.
func Inspect(t testing.TestingT, id string) *ContainerInspect {
	t.Helper()

	return InspectContext(t, context.Background(), id)
}

// InspectContext runs the 'docker inspect {container id}' command and returns a [ContainerInspect] struct,
// converted from the output JSON. This will fail the test if there are any errors. The ctx parameter supports
// cancellation and timeouts.
func InspectContext(t testing.TestingT, ctx context.Context, id string) *ContainerInspect {
	t.Helper()

	out, err := InspectContextE(t, ctx, id)
	require.NoError(t, err)

	return out
}

// InspectE runs the 'docker inspect {container id}' command and returns a ContainerInspect
// struct, converted from the output JSON, along with any errors
//
// Deprecated: Use [InspectContextE] instead.
func InspectE(t testing.TestingT, id string) (*ContainerInspect, error) {
	t.Helper()

	return InspectContextE(t, context.Background(), id)
}

// InspectContextE runs the 'docker inspect {container id}' command and returns a [ContainerInspect] struct,
// converted from the output JSON, along with any errors. The ctx parameter supports cancellation and timeouts.
func InspectContextE(t testing.TestingT, ctx context.Context, id string) (*ContainerInspect, error) {
	t.Helper()

	cmd := &shell.Command{
		Command: "docker",
		Args:    []string{"container", "inspect", id},
		// inspect is a short-running command, don't print the output.
		Logger: logger.Discard,
	}

	out, err := shell.RunCommandContextAndGetStdOutE(t, ctx, cmd)
	if err != nil {
		return nil, err
	}

	var containers []inspectOutput

	err = json.Unmarshal([]byte(out), &containers)
	if err != nil {
		return nil, err
	}

	if len(containers) == 0 {
		return nil, fmt.Errorf("%w: %s", ErrNoContainerFound, id)
	}

	container := &containers[0]

	return transformContainer(t, container)
}

// transformContainer converts 'docker inspect' output JSON into a more friendly and testable format
func transformContainer(t testing.TestingT, container *inspectOutput) (*ContainerInspect, error) {
	t.Helper()

	name := strings.TrimLeft(container.Name, "/")

	ports, err := transformContainerPorts(container)
	if err != nil {
		return nil, err
	}

	volumes := transformContainerVolumes(container)

	created, err := time.Parse(time.RFC3339Nano, container.Created)
	if err != nil {
		return nil, err
	}

	inspect := ContainerInspect{
		ID:       container.ID,
		Name:     name,
		Created:  created,
		Status:   container.State.Status,
		Running:  container.State.Running,
		ExitCode: container.State.ExitCode,
		Error:    container.State.Error,
		Ports:    ports,
		Binds:    volumes,
		Health: HealthCheck{
			Status:        container.State.Health.Status,
			FailingStreak: container.State.Health.FailingStreak,
			Log:           container.State.Health.Log,
		},
	}

	return &inspect, nil
}

// transformContainerPorts converts Docker's ports from the following json into a more testable format
//
//	{
//	  "80/tcp": [
//	    {
//		     "HostIp": ""
//	      "HostPort": "8080"
//	    }
//	  ]
//	}
func transformContainerPorts(container *inspectOutput) ([]Port, error) {
	var ports []Port

	cPorts := container.NetworkSettings.Ports

	for key, portBinding := range cPorts {
		split := strings.Split(key, "/")

		containerPort, err := strconv.ParseUint(split[0], 10, 16)
		if err != nil {
			return nil, err
		}

		var protocol string
		if len(split) > 1 {
			protocol = split[1]
		}

		for _, port := range portBinding {
			hostPort, err := strconv.ParseUint(port.HostPort, 10, 16)
			if err != nil {
				return nil, err
			}

			ports = append(ports, Port{
				HostPort:      uint16(hostPort),
				ContainerPort: uint16(containerPort),
				Protocol:      protocol,
			})
		}
	}

	return ports, nil
}

// GetExposedHostPort returns an exposed host port according to requested container port. Returns 0 if the requested port is not exposed.
func (c *ContainerInspect) GetExposedHostPort(containerPort uint16) uint16 {
	for _, port := range c.Ports {
		if port.ContainerPort == containerPort {
			return port.HostPort
		}
	}

	return uint16(0)
}

// transformContainerVolumes converts Docker's volume bindings from the
// format "/foo/bar:/foo/baz" into a more testable one
func transformContainerVolumes(container *inspectOutput) []VolumeBind {
	binds := container.HostConfig.Binds
	volumes := make([]VolumeBind, 0, len(binds))

	for _, bind := range binds {
		var source, dest string

		split := strings.Split(bind, ":")

		// Considering it as an unbound volume
		dest = split[0]

		if len(split) == volumeBindParts {
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
