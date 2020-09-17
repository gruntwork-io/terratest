package docker

import (
	"crypto/tls"
	"strconv"
	"strings"
	"testing"
	"time"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/stretchr/testify/require"
)

func TestStop(t *testing.T) {
	t.Parallel()

	// appending timestamp to container name to run tests in parallel
	name := "test-nginx" + strconv.FormatInt(time.Now().UnixNano(), 10)

	// choosing a unique port since 80 may not fly well on test machines
	port := "13030"
	testURL := "http://localhost:" + port

	// for testing the stopping of a docker container
	// we got to run a container first and then stop it
	runOpts := &RunOptions{
		Detach:       true,
		Name:         name,
		Remove:       true,
		OtherOptions: []string{"-p", port + ":80"},
	}
	Run(t, "nginx:1.17-alpine", runOpts)

	// verify nginx is running
	http_helper.HttpGetWithRetryWithCustomValidation(t, testURL, &tls.Config{}, 60, 2*time.Second, verifyNginxIsUp)

	// try to stop it now
	out := Stop(t, []string{name}, &StopOptions{})
	require.Contains(t, out, name)

	// verify nginx is down
	// run a docker ps with name filter
	command := shell.Command{
		Command: "docker",
		Args:    []string{"ps", "-q", "--filter", "name=" + name},
	}
	output := shell.RunCommandAndGetStdOut(t, command)
	require.Empty(t, output)
}

func verifyNginxIsUp(statusCode int, body string) bool {
	return statusCode == 200 && strings.Contains(body, "nginx!")
}
