// Package ssh allows to manage SSH connections and send commands through them.
package ssh

import (
	"net"
	"testing"
	"time"

	"fmt"
	"io"
	"os"
	"strconv"

	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/logger"
	"golang.org/x/crypto/ssh"
)

// Host is a host on AWS.
type Host struct {
	Hostname    string
	SshUserName string
	SshKeyPair  *KeyPair
}

// ScpFileToE uploads the contents using SCP to the given host and fails the test if the connection fails.
func ScpFileTo(t *testing.T, host Host, mode os.FileMode, remotePath, contents string) {
	err := ScpFileToE(t, host, mode, remotePath, contents)
	if err != nil {
		t.Fatal(err)
	}
}

// ScpFileToE uploads the contents using SCP to the given host and return an error if the process fails.
func ScpFileToE(t *testing.T, host Host, mode os.FileMode, remotePath, contents string) error {
	authMethods, err := createAuthMethodsForHost(host)
	if err != nil {
		return err
	}
	dir, file := filepath.Split(remotePath)

	hostOptions := SshConnectionOptions{
		Username:    host.SshUserName,
		Address:     host.Hostname,
		Port:        22,
		Command:     "/usr/bin/scp -t " + dir,
		AuthMethods: authMethods,
	}

	scp := sendScpCommandsToCopyFile(mode, file, contents)

	sshSession := &SshSession{
		Options:  &hostOptions,
		JumpHost: &JumpHostSession{},
		Input:    &scp,
	}

	defer sshSession.Cleanup(t)

	_, err = runSSHCommand(t, sshSession)
	return err
}

// CheckSshConnection checks that you can connect via SSH to the given host and fail the test if the connection fails.
func CheckSshConnection(t *testing.T, host Host) {
	err := CheckSshConnectionE(t, host)
	if err != nil {
		t.Fatal(err)
	}
}

// CheckSshConnectionE checks that you can connect via SSH to the given host and return an error if the connection fails.
func CheckSshConnectionE(t *testing.T, host Host) error {
	_, err := CheckSshCommandE(t, host, "'exit'")
	return err
}

// CheckSshCommand checks that you can connect via SSH to the given host and run the given command. Returns the stdout/stderr.
func CheckSshCommand(t *testing.T, host Host, command string) string {
	out, err := CheckSshCommandE(t, host, command)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// CheckSshCommandE checks that you can connect via SSH to the given host and run the given command. Returns the stdout/stderr.
func CheckSshCommandE(t *testing.T, host Host, command string) (string, error) {
	authMethods, err := createAuthMethodsForHost(host)
	if err != nil {
		return "", err
	}

	hostOptions := SshConnectionOptions{
		Username:    host.SshUserName,
		Address:     host.Hostname,
		Port:        22,
		Command:     command,
		AuthMethods: authMethods,
	}

	sshSession := &SshSession{
		Options:  &hostOptions,
		JumpHost: &JumpHostSession{},
	}

	defer sshSession.Cleanup(t)

	return runSSHCommand(t, sshSession)
}

// CheckPrivateSshConnection attempts to connect to privateHost (which is not addressable from the Internet) via a
// separate publicHost (which is addressable from the Internet) and then executes "command" on privateHost and returns
// its output. It is useful for checking that it's possible to SSH from a Bastion Host to a private instance.
func CheckPrivateSshConnection(t *testing.T, publicHost Host, privateHost Host, command string) string {
	out, err := CheckPrivateSshConnectionE(t, publicHost, privateHost, command)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// CheckPrivateSshConnectionE attempts to connect to privateHost (which is not addressable from the Internet) via a
// separate publicHost (which is addressable from the Internet) and then executes "command" on privateHost and returns
// its output. It is useful for checking that it's possible to SSH from a Bastion Host to a private instance.
func CheckPrivateSshConnectionE(t *testing.T, publicHost Host, privateHost Host, command string) (string, error) {
	jumpHostAuthMethods, err := createAuthMethodsForHost(publicHost)
	if err != nil {
		return "", err
	}

	jumpHostOptions := SshConnectionOptions{
		Username:    publicHost.SshUserName,
		Address:     publicHost.Hostname,
		Port:        22,
		AuthMethods: jumpHostAuthMethods,
	}

	hostAuthMethods, err := createAuthMethodsForHost(privateHost)
	if err != nil {
		return "", err
	}

	hostOptions := SshConnectionOptions{
		Username:    privateHost.SshUserName,
		Address:     privateHost.Hostname,
		Port:        22,
		Command:     command,
		AuthMethods: hostAuthMethods,
		JumpHost:    &jumpHostOptions,
	}

	sshSession := &SshSession{
		Options:  &hostOptions,
		JumpHost: &JumpHostSession{},
	}

	defer sshSession.Cleanup(t)

	return runSSHCommand(t, sshSession)
}

func runSSHCommand(t *testing.T, sshSession *SshSession) (string, error) {
	logger.Logf(t, "Running command %s on %s@%s", sshSession.Options.Command, sshSession.Options.Username, sshSession.Options.Address)
	if err := setUpSSHClient(sshSession); err != nil {
		return "", err
	}

	if err := setUpSSHSession(sshSession); err != nil {
		return "", err
	}

	if sshSession.Input != nil {
		w, err := sshSession.Session.StdinPipe()
		if err != nil {
			return "", err
		}
		go func() {
			defer w.Close()
			(*sshSession.Input)(w)
		}()
	}

	bytes, err := sshSession.Session.CombinedOutput(sshSession.Options.Command)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func setUpSSHClient(sshSession *SshSession) error {
	if sshSession.Options.JumpHost == nil {
		return fillSSHClientForHost(sshSession)
	}
	return fillSSHClientForJumpHost(sshSession)
}

func fillSSHClientForHost(sshSession *SshSession) error {
	client, err := createSSHClient(sshSession.Options)

	if err != nil {
		return err
	}

	sshSession.Client = client
	return nil
}

func fillSSHClientForJumpHost(sshSession *SshSession) error {
	jumpHostClient, err := createSSHClient(sshSession.Options.JumpHost)
	if err != nil {
		return err
	}
	sshSession.JumpHost.JumpHostClient = jumpHostClient

	hostVirtualConn, err := jumpHostClient.Dial("tcp", sshSession.Options.ConnectionString())
	if err != nil {
		return err
	}
	sshSession.JumpHost.HostVirtualConnection = hostVirtualConn

	hostConn, hostIncomingChannels, hostIncomingRequests, err := ssh.NewClientConn(hostVirtualConn, sshSession.Options.ConnectionString(), createSSHClientConfig(sshSession.Options))
	if err != nil {
		return err
	}
	sshSession.JumpHost.HostConnection = hostConn

	sshSession.Client = ssh.NewClient(hostConn, hostIncomingChannels, hostIncomingRequests)
	return nil
}

func setUpSSHSession(sshSession *SshSession) error {
	session, err := sshSession.Client.NewSession()
	if err != nil {
		return err
	}

	sshSession.Session = session
	return nil
}

func createSSHClient(options *SshConnectionOptions) (*ssh.Client, error) {
	sshClientConfig := createSSHClientConfig(options)
	return ssh.Dial("tcp", options.ConnectionString(), sshClientConfig)
}

func createSSHClientConfig(hostOptions *SshConnectionOptions) *ssh.ClientConfig {
	clientConfig := &ssh.ClientConfig{
		User: hostOptions.Username,
		Auth: hostOptions.AuthMethods,
		// Do not do a host key check, as Terratest is only used for testing, not prod
		HostKeyCallback: NoOpHostKeyCallback,
		// By default, Go does not impose a timeout, so a SSH connection attempt can hang for a LONG time.
		Timeout: 10 * time.Second,
	}
	clientConfig.SetDefaults()
	return clientConfig
}

// NoOpHostKeyCallback is an ssh.HostKeyCallback that does nothing. Only use this when you're sure you don't want to check the host key at all
// (e.g., only for testing and non-production use cases).
func NoOpHostKeyCallback(hostname string, remote net.Addr, key ssh.PublicKey) error {
	return nil
}

func createAuthMethodsForHost(host Host) ([]ssh.AuthMethod, error) {
	signer, err := ssh.ParsePrivateKey([]byte(host.SshKeyPair.PrivateKey))
	if err != nil {
		return []ssh.AuthMethod{}, err
	}

	return []ssh.AuthMethod{ssh.PublicKeys(signer)}, nil
}

// sendScpCommandsToCopyFile returns a function which will send commands to the SCP binary to output a file on the remote machine.
// A full explanation of the SCP protocol can be found at
// https://web.archive.org/web/20170215184048/https://blogs.oracle.com/janp/entry/how_the_scp_protocol_works
func sendScpCommandsToCopyFile(mode os.FileMode, fileName, contents string) func(io.WriteCloser) {
	return func(input io.WriteCloser) {

		octalMode := "0" + strconv.FormatInt(int64(mode), 8)

		// Create a file at <filename> with Unix permissions set to <octalMost> and the file will be <len(content)> bytes long.
		fmt.Fprintln(input, "C"+octalMode, len(contents), fileName)

		// Actually send the file
		fmt.Fprint(input, contents)

		// End of transfer
		fmt.Fprint(input, "\x00")
	}
}
