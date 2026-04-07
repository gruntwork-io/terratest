package ssh

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"net"
	"os"
	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
	"golang.org/x/crypto/ssh/agent"
)

// SSHAgent is an in-process SSH agent that can be used for SSH authentication in tests.
type SSHAgent struct {
	agent      agent.Agent
	ln         net.Listener
	stop       chan bool
	stopped    chan bool
	socketDir  string
	socketFile string
}

// NewSSHAgent creates an SSH agent, starts it in the background, and returns control back to the main thread.
// You should stop the agent to clean up files afterwards by calling defer s.Stop().
// The ctx parameter is used when establishing the Unix socket listener.
func NewSSHAgent(t testing.TestingT, ctx context.Context, socketDir string, socketFile string) (*SSHAgent, error) {
	var err error

	s := &SSHAgent{
		stop:       make(chan bool),
		stopped:    make(chan bool),
		socketDir:  socketDir,
		socketFile: socketFile,
		agent:      agent.NewKeyring(),
	}

	s.ln, err = (&net.ListenConfig{}).Listen(ctx, "unix", s.socketFile)
	if err != nil {
		return nil, err
	}

	go s.run(t)

	return s, nil
}

// SocketFile returns the path to the SSH agent's Unix socket file.
func (s *SSHAgent) SocketFile() string {
	return s.socketFile
}

// SocketDir returns the path to the directory containing the SSH agent's Unix socket.
func (s *SSHAgent) SocketDir() string {
	return s.socketDir
}

// Agent returns the underlying ssh agent.Agent used by this SSHAgent.
func (s *SSHAgent) Agent() agent.Agent {
	return s.agent
}

// run is the SSH agent listener and handler loop.
func (s *SSHAgent) run(t testing.TestingT) {
	defer close(s.stopped)

	for {
		select {
		case <-s.stop:
			return
		default:
			c, err := s.ln.Accept()
			if err != nil {
				select {
				// When s.Stop() closes the listener, s.ln.Accept() returns an error that can be ignored
				// since the agent is in stopping process.
				case <-s.stop:
					return
				// When s.ln.Accept() returns a legit error, we print it and continue accepting further requests.
				default:
					logger.Default.Logf(t, "could not accept connection to agent %v", err)

					continue
				}
			} else {
				go func(c net.Conn) {
					defer func() { _ = c.Close() }()

					err := agent.ServeAgent(s.agent, c)
					if err != nil {
						logger.Default.Logf(t, "could not serve ssh agent %v", err)
					}
				}(c)
			}
		}
	}
}

// Stop stops the SSH agent, closes its listener, and removes the socket directory.
func (s *SSHAgent) Stop() {
	close(s.stop)
	_ = s.ln.Close()
	<-s.stopped
	_ = os.RemoveAll(s.socketDir)
}

// SSHAgentWithKeyPair creates and returns an in-memory SSH agent with the given KeyPair already added.
// You should stop the agent to clean up files afterwards by calling defer sshAgent.Stop().
// This will fail the test if there is an error.
// The ctx parameter is used when establishing the Unix socket listener.
func SSHAgentWithKeyPair(t testing.TestingT, ctx context.Context, keyPair *KeyPair) *SSHAgent {
	sshAgent, err := SSHAgentWithKeyPairE(t, ctx, keyPair)
	if err != nil {
		t.Fatal(err)
	}

	return sshAgent
}

// SSHAgentWithKeyPairE creates and returns an in-memory SSH agent with the given KeyPair already added.
// You should stop the agent to clean up files afterwards by calling defer sshAgent.Stop().
// The ctx parameter is used when establishing the Unix socket listener.
func SSHAgentWithKeyPairE(t testing.TestingT, ctx context.Context, keyPair *KeyPair) (*SSHAgent, error) {
	return SSHAgentWithKeyPairsE(t, ctx, []*KeyPair{keyPair})
}

// SSHAgentWithKeyPairs creates and returns an in-memory SSH agent with the given KeyPairs already added.
// You should stop the agent to clean up files afterwards by calling defer sshAgent.Stop().
// This will fail the test if there is an error.
// The ctx parameter is used when establishing the Unix socket listener.
func SSHAgentWithKeyPairs(t testing.TestingT, ctx context.Context, keyPairs []*KeyPair) *SSHAgent {
	sshAgent, err := SSHAgentWithKeyPairsE(t, ctx, keyPairs)
	if err != nil {
		t.Fatal(err)
	}

	return sshAgent
}

// SSHAgentWithKeyPairsE creates and returns an in-memory SSH agent with the given KeyPairs already added.
// You should stop the agent to clean up files afterwards by calling defer sshAgent.Stop().
// The ctx parameter is used when establishing the Unix socket listener.
func SSHAgentWithKeyPairsE(t testing.TestingT, ctx context.Context, keyPairs []*KeyPair) (*SSHAgent, error) {
	logger.Default.Logf(t, "Generating SSH Agent with given KeyPair(s)")

	// Instantiate a temporary SSH agent.
	socketDir, err := os.MkdirTemp("", "ssh-agent-")
	if err != nil {
		return nil, err
	}

	socketFile := filepath.Join(socketDir, "ssh_auth.sock")

	sshAgent, err := NewSSHAgent(t, ctx, socketDir, socketFile)
	if err != nil {
		return nil, err
	}

	// Add given ssh keys to the newly created agent.
	for _, keyPair := range keyPairs {
		// Create SSH key for the agent using the given SSH key pair(s).
		block, _ := pem.Decode([]byte(keyPair.PrivateKey))

		privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}

		key := agent.AddedKey{PrivateKey: privateKey}

		if err := sshAgent.agent.Add(key); err != nil {
			return nil, err
		}
	}

	return sshAgent, nil
}
