package http_helper //nolint:staticcheck // package name determined by directory

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// RunDummyServerContext runs a dummy HTTP server on a unique port that will return the given text. Returns the Listener
// for the server, the port it's listening on, or an error if something went wrong while trying to start the listener.
// The provided context is used when establishing the network listener. Make sure to call the Close() method on the
// Listener when you're done!
func RunDummyServerContext(t testing.TestingT, ctx context.Context, text string) (net.Listener, int) {
	listener, port, err := RunDummyServerContextE(t, ctx, text)
	if err != nil {
		t.Fatal(err)
	}

	return listener, port
}

// RunDummyServerContextE runs a dummy HTTP server on a unique port that will return the given text. Returns the Listener
// for the server, the port it's listening on, or an error if something went wrong while trying to start the listener.
// The provided context is used when establishing the network listener. Make sure to call the Close() method on the
// Listener when you're done!
func RunDummyServerContextE(t testing.TestingT, ctx context.Context, text string) (net.Listener, int, error) {
	port := getNextPort()

	// Create new serve mux so that multiple handlers can be created
	server := http.NewServeMux()
	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, text)
	})

	logger.Default.Logf(t, "Starting dummy HTTP server in port %d that will return the text '%s'", port, text)

	listener, err := (&net.ListenConfig{}).Listen(ctx, "tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, 0, fmt.Errorf("error listening: %w", err)
	}

	go func() { _ = http.Serve(listener, server) }()

	return listener, port, nil
}

// RunDummyServer runs a dummy HTTP server on a unique port that will return the given text. Returns the Listener for
// the server, the port it's listening on, or an error if something went wrong while trying to start the listener.
// Make sure to call the Close() method on the Listener when you're done!
//
// Deprecated: Use [RunDummyServerContext] instead.
func RunDummyServer(t testing.TestingT, text string) (net.Listener, int) {
	return RunDummyServerContext(t, context.Background(), text)
}

// RunDummyServerE runs a dummy HTTP server on a unique port that will return the given text. Returns the Listener for
// the server, the port it's listening on, or an error if something went wrong while trying to start the listener.
// Make sure to call the Close() method on the Listener when you're done!
//
// Deprecated: Use [RunDummyServerContextE] instead.
func RunDummyServerE(t testing.TestingT, text string) (net.Listener, int, error) {
	return RunDummyServerContextE(t, context.Background(), text)
}

// RunDummyServerWithHandlersContext runs a dummy HTTP server on a unique port that will serve the given handlers.
// Returns the Listener for the server, the port it's listening on, or an error if something went wrong while trying
// to start the listener. The provided context is used when establishing the network listener. Make sure to call the
// Close() method on the Listener when you're done!
func RunDummyServerWithHandlersContext(t testing.TestingT, ctx context.Context, handlers map[string]func(http.ResponseWriter, *http.Request)) (net.Listener, int) {
	listener, port, err := RunDummyServerWithHandlersContextE(t, ctx, handlers)
	if err != nil {
		t.Fatal(err)
	}

	return listener, port
}

// RunDummyServerWithHandlersContextE runs a dummy HTTP server on a unique port that will serve the given handlers.
// Returns the Listener for the server, the port it's listening on, or an error if something went wrong while trying
// to start the listener. The provided context is used when establishing the network listener. Make sure to call the
// Close() method on the Listener when you're done!
func RunDummyServerWithHandlersContextE(t testing.TestingT, ctx context.Context, handlers map[string]func(http.ResponseWriter, *http.Request)) (net.Listener, int, error) {
	port := getNextPort()

	server := http.NewServeMux()
	for path, handler := range handlers {
		server.HandleFunc(path, handler)
	}

	logger.Default.Logf(t, "Starting dummy HTTP server in port %d", port)

	listener, err := (&net.ListenConfig{}).Listen(ctx, "tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, 0, fmt.Errorf("error listening: %w", err)
	}

	go func() { _ = http.Serve(listener, server) }()

	return listener, port, nil
}

// RunDummyServerWithHandlers runs a dummy HTTP server on a unique port that will serve the given handlers. Returns the
// Listener for the server, the port it's listening on, or an error if something went wrong while trying to start the
// listener. Make sure to call the Close() method on the Listener when you're done!
//
// Deprecated: Use [RunDummyServerWithHandlersContext] instead.
func RunDummyServerWithHandlers(t testing.TestingT, handlers map[string]func(http.ResponseWriter, *http.Request)) (net.Listener, int) {
	return RunDummyServerWithHandlersContext(t, context.Background(), handlers)
}

// RunDummyServerWithHandlersE runs a dummy HTTP server on a unique port that will serve the given handlers. Returns the
// Listener for the server, the port it's listening on, or an error if something went wrong while trying to start the
// listener. Make sure to call the Close() method on the Listener when you're done!
//
// Deprecated: Use [RunDummyServerWithHandlersContextE] instead.
func RunDummyServerWithHandlersE(t testing.TestingT, handlers map[string]func(http.ResponseWriter, *http.Request)) (net.Listener, int, error) {
	return RunDummyServerWithHandlersContextE(t, context.Background(), handlers)
}

// DO NOT ACCESS THIS VARIABLE DIRECTLY. See getNextPort() below.
var testServerPort int32 = 8080

// Since we run tests in parallel, we need to ensure that each test runs on a different port. This function returns a
// unique port by atomically incrementing the testServerPort variable.
func getNextPort() int {
	return int(atomic.AddInt32(&testServerPort, 1))
}
