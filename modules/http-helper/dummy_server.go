package http_helper

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// RunDummyServer runs a dummy HTTP server on a unique port that will return the given text. Returns the Listener for the server, the
// port it's listening on, or an error if something went wrong while trying to start the listener. Make sure to call
// the Close() method on the Listener when you're done!
func RunDummyServer(t testing.TestingT, text string) (net.Listener, int) {
	listener, port, err := RunDummyServerE(t, text)
	if err != nil {
		t.Fatal(err)
	}
	return listener, port
}

// RunDummyServerE runs a dummy HTTP server on a unique port that will return the given text. Returns the Listener for the server, the
// port it's listening on, or an error if something went wrong while trying to start the listener. Make sure to call
// the Close() method on the Listener when you're done!
func RunDummyServerE(t testing.TestingT, text string) (net.Listener, int, error) {
	port := getNextPort()

	// Create new serve mux so that multiple handlers can be created
	server := http.NewServeMux()
	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, text)
	})

	logger.Logf(t, "Starting dummy HTTP server in port %d that will return the text '%s'", port, text)

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, 0, fmt.Errorf("error listening: %s", err)
	}

	go http.Serve(listener, server)

	return listener, port, err
}

// RunDummyServerWithHandlers runs a dummy HTTP server on a unique port that will serve given handlers. Returns the Listener for the server,
// the port it's listening on, or an error if something went wrong while trying to start the listener. Make sure to call
// the Close() method on the Listener when you're done!
func RunDummyServerWithHandlers(t testing.TestingT, handlers map[string]func(http.ResponseWriter, *http.Request)) (net.Listener, int) {
	listener, port, err := RunDummyServerWithHandlersE(t, handlers)
	if err != nil {
		t.Fatal(err)
	}
	return listener, port
}

// RunDummyServerWithHandlersE runs a dummy HTTP server on a unique port that will server given handlers. Returns the Listener for the server,
// the port it's listening on, or an error if something went wrong while trying to start the listener. Make sure to call
// the Close() method on the Listener when you're done!
func RunDummyServerWithHandlersE(t testing.TestingT, handlers map[string]func(http.ResponseWriter, *http.Request)) (net.Listener, int, error) {
	port := getNextPort()

	for path, handler := range handlers {
		http.HandleFunc(path, handler)
	}

	logger.Logf(t, "Starting dummy HTTP server in port %d", port)

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, 0, fmt.Errorf("error listening: %s", err)
	}

	go http.Serve(listener, nil)

	return listener, port, err
}

// DO NOT ACCESS THIS VARIABLE DIRECTLY. See getNextPort() below.
var testServerPort int32 = 8080

// Since we run tests in parallel, we need to ensure that each test runs on a different port. This function returns a
// unique port by atomically incrementing the testServerPort variable.
func getNextPort() int {
	return int(atomic.AddInt32(&testServerPort, 1))
}
