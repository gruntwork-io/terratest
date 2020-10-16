package podman

import (
	"encoding/json"
	"fmt"
	"strconv"
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
	ID      string `json:"Id"`
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

type inspectOutputV2 struct {
	ID      string
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
		Ports inspectPortsV2
	}
	HostConfig struct {
		Binds []string
	}
}

type inspectPortsV2 map[string][]PortV2

type PortV2 struct {
	HostIP   string `json:"HostIp"`
	HostPort string
}

type Version struct {
	Client struct {
		Version string
	}
}

// Inspect runs the 'podman inspect {container id}' command and returns a ContainerInspect
// struct, converted from the output JSON, along with any errors
func Inspect(t *testing.T, id string) *ContainerInspect {
	out, err := InspectE(t, id)
	require.NoError(t, err)

	return out
}

func getVersion(t *testing.T) (*Version, error) {
	var (
		version Version
		err     error
	)

	cmd := shell.Command{
		Command: "podman",
		Args:    []string{"version", "--format", "json"},
		// inspect is a short-running command, don't print the output.
		Logger: logger.Discard,
	}

	out, err := shell.RunCommandAndGetStdOutE(t, cmd)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(out), &version)
	if err != nil {
		return nil, err
	}

	return &version, nil
}

func convertV2Output(containersV2 []inspectOutputV2) (*inspectOutput, error) {
	containerV2 := containersV2[0]
	container := inspectOutput{}
	container.ID = containerV2.ID
	container.Created = containerV2.Created
	container.Name = containerV2.Name
	container.State = containerV2.State
	container.HostConfig = containerV2.HostConfig

	var ports []Port

	ports, err := transformContainerPorts(containerV2.NetworkSettings.Ports)
	if err != nil {
		return nil, err
	}

	container.NetworkSettings.Ports = ports

	return &container, nil
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

	version, err := getVersion(t)
	if err != nil {
		return nil, err
	}

	var containers []inspectOutput

	if !strings.HasPrefix(version.Client.Version, "1") {
		// Podman 2.x.x has a different inspect output
		// It is now similar to Docker inspect output
		// We will convert V2 to v1 here
		var containersV2 []inspectOutputV2

		err = json.Unmarshal([]byte(out), &containersV2)
		if err != nil {
			return nil, err
		}

		var containerOutput *inspectOutput

		containerOutput, err = convertV2Output(containersV2)
		if err != nil {
			return nil, err
		}

		containers = append(containers, *containerOutput)
	} else {
		err = json.Unmarshal([]byte(out), &containers)
		if err != nil {
			return nil, err
		}
	}

	if len(containers) == 0 {
		return nil, fmt.Errorf("no container found with ID %s", id)
	}

	container := containers[0]

	return transformContainer(container)
}

// transformContainerPorts converts 'docker inspect' output JSON into a more friendly and testable format
func transformContainer(container inspectOutput) (*ContainerInspect, error) {
	name := strings.TrimLeft(container.Name, "/")

	volumes := transformContainerVolumes(container.HostConfig.Binds)

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
func transformContainerVolumes(binds []string) []VolumeBind {
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

// transformContainerPorts converts podman's v2.x ports from the following json into a more testable format
// {
//   "80/tcp": [
//     {
// 	     "HostIp": ""
//       "HostPort": "8080"
//     }
//   ]
// }
func transformContainerPorts(cPorts inspectPortsV2) ([]Port, error) {
	var ports []Port

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
