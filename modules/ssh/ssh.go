// Package ssh allows to manage SSH connections and send commands through them.
package ssh

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gruntwork-io/terratest/modules/core/v2/files"
	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/retry"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/hashicorp/go-multierror"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

const (
	// defaultSSHPort is the standard SSH port number.
	defaultSSHPort = 22

	// defaultDirPermissions is the default directory permissions used when creating local directories.
	defaultDirPermissions = 0o755

	// sshConnectionTimeout is the timeout for establishing an SSH connection.
	sshConnectionTimeout = 10 * time.Second
)

// ErrNoAuthMethod is returned when no authentication method (key pair, agent, or password) is configured on a [Host].
var ErrNoAuthMethod = errors.New("no authentication method defined")

// Host is a remote host. Set one or more authentication methods on the host;
// the first valid method will be used.
type Host struct {
	// SshKeyPair is the SSH key pair to use for authentication. Disabled by default.
	SshKeyPair *KeyPair //nolint:staticcheck,revive // preserving existing field name
	// OverrideSshAgent enables an in-process [SSHAgent] for connections to this host. Disabled by default.
	OverrideSshAgent *SSHAgent //nolint:staticcheck,revive // preserving existing field name
	// Hostname is the host name or IP address.
	Hostname string
	// SshUserName is the SSH user name.
	SshUserName string //nolint:staticcheck,revive // preserving existing field name
	// Password is the plain text password for authentication. Blank by default.
	Password string
	// CustomPort is the port number to use to connect to the host. Port 22 is used if unset.
	CustomPort int
	// SshAgent enables authentication using the existing local SSH agent. Disabled by default.
	SshAgent bool //nolint:staticcheck,revive // preserving existing field name
}

// SCPDownloadOptions configures the parameters for downloading files from a remote host via SCP.
type SCPDownloadOptions struct {
	// RemoteDir is the directory on the remote machine to copy files from.
	RemoteDir string
	// LocalDir is the directory on the local machine to copy files to.
	LocalDir string
	// FileNameFilters are file name patterns to match. May include bash-style wildcards (e.g., *.log).
	FileNameFilters []string
	// RemoteHost is the connection information for the remote machine.
	RemoteHost Host
	// MaxFileSizeMB is the maximum file size in megabytes to download. Files larger than this are skipped.
	MaxFileSizeMB int
}

// GetPort returns the port to use for SSH connections. If [Host.CustomPort] is set,
// it returns that value; otherwise, it returns the default SSH port 22.
func (h *Host) GetPort() int {
	if h.CustomPort == 0 {
		return defaultSSHPort
	}

	return h.CustomPort
}

// SCPFileToContext uploads the contents using SCP to the given host.
// This will fail the test if the connection fails.
// The ctx parameter supports cancellation and timeouts.
func SCPFileToContext(t testing.TestingT, ctx context.Context, host *Host, mode os.FileMode, remotePath, contents string) {
	err := SCPFileToContextE(t, ctx, host, mode, remotePath, contents)
	if err != nil {
		t.Fatal(err)
	}
}

// SCPFileToContextE uploads the contents using SCP to the given host and returns an error if the process fails.
// The ctx parameter supports cancellation and timeouts.
func SCPFileToContextE(t testing.TestingT, ctx context.Context, host *Host, mode os.FileMode, remotePath, contents string) error {
	authMethods, err := createAuthMethodsForHost(ctx, host)
	if err != nil {
		return err
	}

	dir, file := filepath.Split(remotePath)

	hostOptions := SSHConnectionOptions{
		Username:    host.SshUserName,
		Address:     host.Hostname,
		Port:        host.GetPort(),
		Command:     "/usr/bin/scp -t " + dir,
		AuthMethods: authMethods,
	}

	scp := sendScpCommandsToCopyFile(mode, file, contents)

	sshSession := &SSHSession{
		Options:  &hostOptions,
		JumpHost: &JumpHostSession{},
		Input:    &scp,
	}

	defer sshSession.Cleanup(t)

	_, err = runSSHCommand(ctx, t, sshSession)

	return err
}

// SCPFileFromContext downloads the file from remotePath on the given host using SCP.
// This will fail the test if the connection fails.
// The ctx parameter supports cancellation and timeouts.
func SCPFileFromContext(t testing.TestingT, ctx context.Context, host *Host, remotePath string, localDestination *os.File, useSudo bool) {
	err := SCPFileFromContextE(t, ctx, host, remotePath, localDestination, useSudo)
	if err != nil {
		t.Fatal(err)
	}
}

// SCPFileFromContextE downloads the file from remotePath on the given host using SCP
// and returns an error if the process fails.
// The ctx parameter supports cancellation and timeouts.
func SCPFileFromContextE(t testing.TestingT, ctx context.Context, host *Host, remotePath string, localDestination *os.File, useSudo bool) error {
	authMethods, err := createAuthMethodsForHost(ctx, host)
	if err != nil {
		return err
	}

	dir := filepath.Dir(remotePath)

	hostOptions := SSHConnectionOptions{
		Username:    host.SshUserName,
		Address:     host.Hostname,
		Port:        host.GetPort(),
		Command:     "/usr/bin/scp -t " + dir,
		AuthMethods: authMethods,
	}

	sshSession := &SSHSession{
		Options:  &hostOptions,
		JumpHost: &JumpHostSession{},
	}

	defer sshSession.Cleanup(t)

	return copyFileFromRemote(ctx, t, sshSession, localDestination, remotePath, useSudo)
}

// SCPDirFromContext downloads all the files from remotePath on the given host using SCP.
// This will fail the test if the connection fails.
// The ctx parameter supports cancellation and timeouts.
func SCPDirFromContext(t testing.TestingT, ctx context.Context, options *SCPDownloadOptions, useSudo bool) {
	err := SCPDirFromContextE(t, ctx, options, useSudo)
	if err != nil {
		t.Fatal(err)
	}
}

// SCPDirFromContextE downloads all the files from remotePath on the given host using SCP
// and returns an error if the process fails. Only files within remotePath will
// be downloaded. This function will not recursively download subdirectories or follow
// symlinks.
// The ctx parameter supports cancellation and timeouts.
func SCPDirFromContextE(t testing.TestingT, ctx context.Context, options *SCPDownloadOptions, useSudo bool) error {
	authMethods, err := createAuthMethodsForHost(ctx, &options.RemoteHost)
	if err != nil {
		return err
	}

	hostOptions := SSHConnectionOptions{
		Username:    options.RemoteHost.SshUserName,
		Address:     options.RemoteHost.Hostname,
		Port:        options.RemoteHost.GetPort(),
		Command:     "/usr/bin/scp -t " + options.RemoteDir,
		AuthMethods: authMethods,
	}

	sshSession := &SSHSession{
		Options:  &hostOptions,
		JumpHost: &JumpHostSession{},
	}

	defer sshSession.Cleanup(t)

	filesInDir, err := listFileInRemoteDir(ctx, t, sshSession, options, useSudo)
	if err != nil {
		return err
	}

	if !files.FileExists(options.LocalDir) {
		err := os.MkdirAll(options.LocalDir, defaultDirPermissions)
		if err != nil {
			return err
		}
	}

	errorsOccurred := new(multierror.Error)

	for _, fullRemoteFilePath := range filesInDir {
		fileName := filepath.Base(fullRemoteFilePath)
		localFilePath := filepath.Join(options.LocalDir, fileName)

		localFile, err := os.Create(localFilePath)
		if err != nil {
			return err
		}

		logger.Default.Logf(t, "Copying remote file: %s to local path %s", fullRemoteFilePath, localFilePath)

		err = copyFileFromRemote(ctx, t, sshSession, localFile, fullRemoteFilePath, useSudo)

		if closeErr := localFile.Close(); closeErr != nil && err == nil {
			err = closeErr
		}

		errorsOccurred = multierror.Append(errorsOccurred, err)
	}

	return errorsOccurred.ErrorOrNil()
}

// CheckSSHConnectionContext checks that you can connect via SSH to the given host.
// This will fail the test if the connection fails.
// The ctx parameter supports cancellation and timeouts.
func CheckSSHConnectionContext(t testing.TestingT, ctx context.Context, host *Host) {
	err := CheckSSHConnectionContextE(t, ctx, host)
	if err != nil {
		t.Fatal(err)
	}
}

// CheckSSHConnectionContextE checks that you can connect via SSH to the given host
// and returns an error if the connection fails.
// The ctx parameter supports cancellation and timeouts.
func CheckSSHConnectionContextE(t testing.TestingT, ctx context.Context, host *Host) error {
	_, err := CheckSSHCommandContextE(t, ctx, host, "'exit'")

	return err
}

// CheckSSHConnectionWithRetryContext attempts to connect via SSH until max retries has been exceeded.
// This will fail the test if the connection fails.
// The ctx parameter supports cancellation and timeouts.
func CheckSSHConnectionWithRetryContext(t testing.TestingT, ctx context.Context, host *Host, retries int, sleepBetweenRetries time.Duration, f ...func(testing.TestingT, context.Context, *Host) error) {
	err := CheckSSHConnectionWithRetryContextE(t, ctx, host, retries, sleepBetweenRetries, f...)
	if err != nil {
		t.Fatal(err)
	}
}

// CheckSSHConnectionWithRetryContextE attempts to connect via SSH until max retries has been exceeded
// and returns an error if the connection fails.
// The ctx parameter supports cancellation and timeouts.
func CheckSSHConnectionWithRetryContextE(t testing.TestingT, ctx context.Context, host *Host, retries int, sleepBetweenRetries time.Duration, f ...func(testing.TestingT, context.Context, *Host) error) error {
	handler := CheckSSHConnectionContextE

	if len(f) > 0 {
		handler = f[0]
	}

	_, err := retry.DoWithRetryContextE(t, ctx, "Checking SSH connection to "+host.Hostname, retries, sleepBetweenRetries, func() (string, error) {
		return "", handler(t, ctx, host)
	})

	return err
}

// CheckSSHCommandContext checks that you can connect via SSH to the given host and run the given command.
// Returns the stdout/stderr. This will fail the test if the connection fails.
// The ctx parameter supports cancellation and timeouts.
func CheckSSHCommandContext(t testing.TestingT, ctx context.Context, host *Host, command string) string {
	out, err := CheckSSHCommandContextE(t, ctx, host, command)
	if err != nil {
		t.Fatal(err)
	}

	return out
}

// CheckSSHCommandContextE checks that you can connect via SSH to the given host and run the given command.
// Returns the stdout/stderr.
// The ctx parameter supports cancellation and timeouts.
func CheckSSHCommandContextE(t testing.TestingT, ctx context.Context, host *Host, command string) (string, error) {
	authMethods, err := createAuthMethodsForHost(ctx, host)
	if err != nil {
		return "", err
	}

	hostOptions := SSHConnectionOptions{
		Username:    host.SshUserName,
		Address:     host.Hostname,
		Port:        host.GetPort(),
		Command:     command,
		AuthMethods: authMethods,
	}

	sshSession := &SSHSession{
		Options:  &hostOptions,
		JumpHost: &JumpHostSession{},
	}

	defer sshSession.Cleanup(t)

	return runSSHCommand(ctx, t, sshSession)
}

// CheckSSHCommandWithRetryContext checks that you can connect via SSH to the given host and run the given command
// until max retries have been exceeded. Returns the stdout/stderr.
// This will fail the test if the connection fails.
// The ctx parameter supports cancellation and timeouts.
func CheckSSHCommandWithRetryContext(t testing.TestingT, ctx context.Context, host *Host, command string, retries int, sleepBetweenRetries time.Duration, f ...func(testing.TestingT, context.Context, *Host, string) (string, error)) string {
	out, err := CheckSSHCommandWithRetryContextE(t, ctx, host, command, retries, sleepBetweenRetries, f...)
	if err != nil {
		t.Fatal(err)
	}

	return out
}

// CheckSSHCommandWithRetryContextE checks that you can connect via SSH to the given host and run the given command
// until max retries has been exceeded. Returns an error if the command fails after max retries has been exceeded.
// The ctx parameter supports cancellation and timeouts.
func CheckSSHCommandWithRetryContextE(t testing.TestingT, ctx context.Context, host *Host, command string, retries int, sleepBetweenRetries time.Duration, f ...func(testing.TestingT, context.Context, *Host, string) (string, error)) (string, error) {
	handler := CheckSSHCommandContextE

	if len(f) > 0 {
		handler = f[0]
	}

	return retry.DoWithRetryContextE(t, ctx, "Checking SSH connection to "+host.Hostname, retries, sleepBetweenRetries, func() (string, error) {
		return handler(t, ctx, host, command)
	})
}

// CheckPrivateSSHConnectionContext attempts to connect to privateHost (which is not addressable from the Internet) via a
// separate publicHost (which is addressable from the Internet) and then executes "command" on privateHost and returns
// its output. It is useful for checking that it's possible to SSH from a Bastion Host to a private instance.
// This will fail the test if the connection fails.
// The ctx parameter supports cancellation and timeouts.
func CheckPrivateSSHConnectionContext(t testing.TestingT, ctx context.Context, publicHost *Host, privateHost *Host, command string) string {
	out, err := CheckPrivateSSHConnectionContextE(t, ctx, publicHost, privateHost, command)
	if err != nil {
		t.Fatal(err)
	}

	return out
}

// CheckPrivateSSHConnectionContextE attempts to connect to privateHost (which is not addressable from the Internet) via a
// separate publicHost (which is addressable from the Internet) and then executes "command" on privateHost and returns
// its output. It is useful for checking that it's possible to SSH from a Bastion Host to a private instance.
// The ctx parameter supports cancellation and timeouts.
func CheckPrivateSSHConnectionContextE(t testing.TestingT, ctx context.Context, publicHost *Host, privateHost *Host, command string) (string, error) {
	jumpHostAuthMethods, err := createAuthMethodsForHost(ctx, publicHost)
	if err != nil {
		return "", err
	}

	jumpHostOptions := SSHConnectionOptions{
		Username:    publicHost.SshUserName,
		Address:     publicHost.Hostname,
		Port:        publicHost.GetPort(),
		AuthMethods: jumpHostAuthMethods,
	}

	hostAuthMethods, err := createAuthMethodsForHost(ctx, privateHost)
	if err != nil {
		return "", err
	}

	hostOptions := SSHConnectionOptions{
		Username:    privateHost.SshUserName,
		Address:     privateHost.Hostname,
		Port:        privateHost.GetPort(),
		Command:     command,
		AuthMethods: hostAuthMethods,
		JumpHost:    &jumpHostOptions,
	}

	sshSession := &SSHSession{
		Options:  &hostOptions,
		JumpHost: &JumpHostSession{},
	}

	defer sshSession.Cleanup(t)

	return runSSHCommand(ctx, t, sshSession)
}

// FetchContentsOfFilesContext connects to the given host via SSH and fetches the contents of the files at the given filePaths.
// If useSudo is true, then the contents will be retrieved using sudo. Returns a map from file path to contents.
// This will fail the test if the connection fails.
// The ctx parameter supports cancellation and timeouts.
func FetchContentsOfFilesContext(t testing.TestingT, ctx context.Context, host *Host, useSudo bool, filePaths ...string) map[string]string {
	out, err := FetchContentsOfFilesContextE(t, ctx, host, useSudo, filePaths...)
	if err != nil {
		t.Fatal(err)
	}

	return out
}

// FetchContentsOfFilesContextE connects to the given host via SSH and fetches the contents of the files at the given filePaths.
// If useSudo is true, then the contents will be retrieved using sudo. Returns a map from file path to contents.
// The ctx parameter supports cancellation and timeouts.
func FetchContentsOfFilesContextE(t testing.TestingT, ctx context.Context, host *Host, useSudo bool, filePaths ...string) (map[string]string, error) {
	filePathToContents := map[string]string{}

	for _, filePath := range filePaths {
		contents, err := FetchContentsOfFileContextE(t, ctx, host, useSudo, filePath)
		if err != nil {
			return nil, err
		}

		filePathToContents[filePath] = contents
	}

	return filePathToContents, nil
}

// FetchContentsOfFileContext connects to the given host via SSH and fetches the contents of the file at the given filePath.
// If useSudo is true, then the contents will be retrieved using sudo. Returns the contents of that file.
// This will fail the test if the connection fails.
// The ctx parameter supports cancellation and timeouts.
func FetchContentsOfFileContext(t testing.TestingT, ctx context.Context, host *Host, useSudo bool, filePath string) string {
	out, err := FetchContentsOfFileContextE(t, ctx, host, useSudo, filePath)
	if err != nil {
		t.Fatal(err)
	}

	return out
}

// shellQuote wraps a path in single quotes, escaping any embedded single quotes,
// so that paths containing spaces or shell metacharacters work correctly when
// passed to commands like `cat` and `dd if=`.
func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

// FetchContentsOfFileContextE connects to the given host via SSH and fetches the contents of the file at the given filePath.
// If useSudo is true, then the contents will be retrieved using sudo. Returns the contents of that file.
// The ctx parameter supports cancellation and timeouts.
func FetchContentsOfFileContextE(t testing.TestingT, ctx context.Context, host *Host, useSudo bool, filePath string) (string, error) {
	command := "cat " + shellQuote(filePath)
	if useSudo {
		command = "sudo " + command
	}

	return CheckSSHCommandContextE(t, ctx, host, command)
}

func listFileInRemoteDir(ctx context.Context, t testing.TestingT, sshSession *SSHSession, options *SCPDownloadOptions, useSudo bool) ([]string, error) {
	logger.Default.Logf(t, "Running command %s on %s@%s", sshSession.Options.Command, sshSession.Options.Username, sshSession.Options.Address)

	var findCommandArgs []string

	if useSudo {
		findCommandArgs = append(findCommandArgs, "sudo")
	}

	findCommandArgs = append(findCommandArgs, "find", options.RemoteDir, "-type", "f")

	filtersLength := len(options.FileNameFilters)

	if options.FileNameFilters != nil && filtersLength > 0 {
		findCommandArgs = append(findCommandArgs, "\\(")

		for i, curFilter := range options.FileNameFilters {

			curFilter = fmt.Sprintf("'%s'", curFilter)
			findCommandArgs = append(findCommandArgs, "-name", curFilter)

			if filtersLength-i > 1 {
				findCommandArgs = append(findCommandArgs, "-o")
			}
		}

		findCommandArgs = append(findCommandArgs, "\\)")
	}

	if options.MaxFileSizeMB != 0 {
		findCommandArgs = append(findCommandArgs, "-size", fmt.Sprintf("-%dM", options.MaxFileSizeMB))
	}

	finalCommandString := strings.Join(findCommandArgs, " ")

	resultString, err := CheckSSHCommandContextE(t, ctx, &options.RemoteHost, finalCommandString)
	if err != nil {
		return nil, err
	}

	if len(resultString) > 0 {
		resultString = resultString[:len(resultString)-1]
	}

	return strings.Split(resultString, "\n"), nil
}

// copyFileFromRemote copies a file from a remote host to a local file.
// Based on code: https://github.com/bramvdbogaerde/go-scp/pull/6/files
func copyFileFromRemote(ctx context.Context, t testing.TestingT, sshSession *SSHSession, file *os.File, remotePath string, useSudo bool) error {

	defer func() { _ = file.Close() }()

	if err := setUpSSHClient(ctx, sshSession); err != nil {
		return err
	}

	if err := setUpSSHSession(sshSession); err != nil {
		return err
	}

	command := "dd if=" + shellQuote(remotePath)
	if useSudo {
		command = "sudo " + command
	}

	logger.Default.Logf(t, "Running command %s on %s@%s", command, sshSession.Options.Username, sshSession.Options.Address)

	defer func() { _ = sshSession.Session.Close() }()

	r, err := sshSession.Session.Output(command)
	if err != nil {
		return fmt.Errorf("error reading from remote stdout: %w", err)
	}

	_, err = file.Write(r)

	return err
}

func runSSHCommand(ctx context.Context, t testing.TestingT, sshSession *SSHSession) (string, error) {
	logger.Default.Logf(t, "Running command %s on %s@%s", sshSession.Options.Command, sshSession.Options.Username, sshSession.Options.Address)

	if err := setUpSSHClient(ctx, sshSession); err != nil {
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
			defer func() { _ = w.Close() }()

			(*sshSession.Input)(w)
		}()
	}

	bytes, err := sshSession.Session.CombinedOutput(sshSession.Options.Command)
	if err != nil {
		return string(bytes), err
	}

	return string(bytes), nil
}

func setUpSSHClient(ctx context.Context, sshSession *SSHSession) error {
	if sshSession.Options.JumpHost == nil {
		return fillSSHClientForHost(ctx, sshSession)
	}

	return fillSSHClientForJumpHost(ctx, sshSession)
}

func fillSSHClientForHost(ctx context.Context, sshSession *SSHSession) error {
	client, err := createSSHClient(ctx, sshSession.Options)
	if err != nil {
		return err
	}

	sshSession.Client = client

	return nil
}

func fillSSHClientForJumpHost(ctx context.Context, sshSession *SSHSession) error {
	jumpHostClient, err := createSSHClient(ctx, sshSession.Options.JumpHost)
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

func setUpSSHSession(sshSession *SSHSession) error {
	session, err := sshSession.Client.NewSession()
	if err != nil {
		return err
	}

	sshSession.Session = session

	return nil
}

func createSSHClient(ctx context.Context, options *SSHConnectionOptions) (*ssh.Client, error) {
	sshClientConfig := createSSHClientConfig(options)

	conn, err := (&net.Dialer{Timeout: sshClientConfig.Timeout}).DialContext(ctx, "tcp", options.ConnectionString())
	if err != nil {
		return nil, err
	}

	c, chans, reqs, err := ssh.NewClientConn(conn, options.ConnectionString(), sshClientConfig)
	if err != nil {
		_ = conn.Close()

		return nil, err
	}

	return ssh.NewClient(c, chans, reqs), nil
}

func createSSHClientConfig(hostOptions *SSHConnectionOptions) *ssh.ClientConfig {
	clientConfig := &ssh.ClientConfig{
		User: hostOptions.Username,
		Auth: hostOptions.AuthMethods,

		HostKeyCallback: NoOpHostKeyCallback,

		Timeout: sshConnectionTimeout,
	}
	clientConfig.SetDefaults()

	return clientConfig
}

// NoOpHostKeyCallback is an ssh.HostKeyCallback that does nothing. Only use this when you're sure you don't want to
// check the host key at all (e.g., only for testing and non-production use cases).
func NoOpHostKeyCallback(hostname string, remote net.Addr, key ssh.PublicKey) error {
	return nil
}

// createAuthMethodsForHost returns an array of authentication methods for the given host.
func createAuthMethodsForHost(ctx context.Context, host *Host) ([]ssh.AuthMethod, error) {
	var methods []ssh.AuthMethod

	if host.OverrideSshAgent != nil {
		conn, err := (&net.Dialer{}).DialContext(ctx, "unix", host.OverrideSshAgent.socketFile)
		if err != nil {
			return methods, fmt.Errorf("failed to dial in-memory ssh agent: %w", err)
		}

		agentClient := agent.NewClient(conn)
		methods = append(methods, ssh.PublicKeysCallback(agentClient.Signers))
	}

	if host.SshAgent {
		socket := os.Getenv("SSH_AUTH_SOCK")

		conn, err := (&net.Dialer{}).DialContext(ctx, "unix", socket)
		if err != nil {
			return methods, err
		}

		agentClient := agent.NewClient(conn)
		methods = append(methods, ssh.PublicKeysCallback(agentClient.Signers))
	}

	if host.SshKeyPair != nil {
		signer, err := ssh.ParsePrivateKey([]byte(host.SshKeyPair.PrivateKey))
		if err != nil {
			return methods, err
		}

		publicKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(host.SshKeyPair.PublicKey))
		if err != nil {
			return methods, err
		}

		if cert, ok := publicKey.(*ssh.Certificate); ok {
			signer, err = ssh.NewCertSigner(cert, signer)
			if err != nil {
				return methods, err
			}
		}

		methods = append(methods, ssh.PublicKeys(signer))
	}

	if len(host.Password) > 0 {
		methods = append(methods, ssh.Password(host.Password))
	}

	if len(methods) < 1 {
		return methods, ErrNoAuthMethod
	}

	return methods, nil
}

// sendScpCommandsToCopyFile returns a function which will send commands to the SCP binary to output a file on the remote machine.
// A full explanation of the SCP protocol can be found at
// https://web.archive.org/web/20170215184048/https://blogs.oracle.com/janp/entry/how_the_scp_protocol_works
func sendScpCommandsToCopyFile(mode os.FileMode, fileName, contents string) func(io.WriteCloser) {
	return func(input io.WriteCloser) {
		octalMode := "0" + strconv.FormatInt(int64(mode), 8)

		_, _ = fmt.Fprintln(input, "C"+octalMode, len(contents), fileName)

		_, _ = fmt.Fprint(input, contents)

		_, _ = fmt.Fprint(input, "\x00")
	}
}
