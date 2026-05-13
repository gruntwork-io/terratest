package docker

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Login(t *testing.T) {

	err := LoginE(t, LoginOptions{Registry: "registry-1.docker.io", Login: "1", Password: "1"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "incorrect username or password")
}
