package ssh

import (
	"io"
	"net"
	"reflect"
	"slices"
	"strconv"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
	"golang.org/x/crypto/ssh"
)

// SSHConnectionOptions are the options for an SSH connection.
type SSHConnectionOptions struct {
	// JumpHost is the optional jump host connection options for tunneling.
	JumpHost *SSHConnectionOptions
	// Username is the SSH user name.
	Username string
	// Address is the host address.
	Address string
	// Command is the command to run on the remote host.
	Command string
	// AuthMethods are the SSH authentication methods to use.
	AuthMethods []ssh.AuthMethod
	// Port is the SSH port number.
	Port int
}

// SshConnectionOptions is a backwards-compatible alias for [SSHConnectionOptions].
//
// Deprecated: Use [SSHConnectionOptions] instead.
type SshConnectionOptions = SSHConnectionOptions //nolint:staticcheck,revive // preserving deprecated type name

// ConnectionString returns the connection string for an SSH connection.
func (options *SSHConnectionOptions) ConnectionString() string {
	return net.JoinHostPort(options.Address, strconv.Itoa(options.Port))
}

// SSHSession is a container object for all resources created by an SSH session. The reason we need this is so that we
// can do a single defer in a top-level method that calls the Cleanup method to go through and ensure all of these
// resources are released and cleaned up.
type SSHSession struct {
	// Options are the SSH connection options.
	Options *SSHConnectionOptions
	// Client is the SSH client.
	Client *ssh.Client
	// Session is the SSH session.
	Session *ssh.Session
	// JumpHost is the optional jump host session for tunneling.
	JumpHost *JumpHostSession
	// Input is an optional function that writes to the session's stdin pipe.
	Input *func(io.WriteCloser)
}

// SshSession is a backwards-compatible alias for [SSHSession].
//
// Deprecated: Use [SSHSession] instead.
type SshSession = SSHSession //nolint:staticcheck,revive // preserving deprecated type name

// Cleanup cleans up an existing SSH session.
func (sshSession *SSHSession) Cleanup(t testing.TestingT) {
	if sshSession == nil {
		return
	}

	// Closing the session may result in an EOF error if it's already closed (e.g. due to hitting CTRL + D), so
	// don't report those errors, as there is nothing actually wrong in that case.
	Close(t, sshSession.Session, io.EOF.Error())
	Close(t, sshSession.Client)
	sshSession.JumpHost.Cleanup(t)
}

// JumpHostSession is a session with a jump host used for tunneling SSH connections.
type JumpHostSession struct {
	// JumpHostClient is the SSH client for the jump host.
	JumpHostClient *ssh.Client
	// HostVirtualConnection is the virtual connection to the target host through the jump host.
	HostVirtualConnection net.Conn
	// HostConnection is the SSH connection to the target host.
	HostConnection ssh.Conn
}

// Cleanup cleans the jump host session up.
func (jumpHost *JumpHostSession) Cleanup(t testing.TestingT) {
	if jumpHost == nil {
		return
	}

	// Closing a connection may result in an EOF error if it's already closed (e.g. due to hitting CTRL + D), so
	// don't report those errors, as there is nothing actually wrong in that case.
	Close(t, jumpHost.HostConnection, io.EOF.Error())
	Close(t, jumpHost.HostVirtualConnection, io.EOF.Error())
	Close(t, jumpHost.JumpHostClient)
}

// Closeable can be closed.
type Closeable interface {
	Close() error
}

// Close closes a Closeable.
func Close(t testing.TestingT, closeable Closeable, ignoreErrors ...string) {
	if interfaceIsNil(closeable) {
		return
	}

	if err := closeable.Close(); err != nil && !slices.Contains(ignoreErrors, err.Error()) {
		logger.Default.Logf(t, "Error closing %s: %s", closeable, err.Error())
	}
}

// interfaceIsNil checks whether the given interface value is nil. A direct nil comparison does not work for interface
// values that wrap a typed nil pointer, so reflection is used.
// See https://go.dev/doc/faq#nil_error for details.
func interfaceIsNil(i interface{}) bool {
	return i == nil || reflect.ValueOf(i).IsNil()
}
