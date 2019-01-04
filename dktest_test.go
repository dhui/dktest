package dktest_test

import (
	"net/http"
	"net/url"
	"testing"
)

import (
	"github.com/docker/go-connections/nat"
)

import (
	"github.com/dhui/dktest"
)

const (
	testImage        = "alpine:3.8"
	testNetworkImage = "nginx:alpine"
)

// ready functions
func alwaysReady(dktest.ContainerInfo) bool { return true }
func neverReady(dktest.ContainerInfo) bool  { return false }
func nginxReady(c dktest.ContainerInfo) bool {
	u := url.URL{Scheme: "http", Host: c.IP + ":" + c.Port}
	if resp, err := http.Get(u.String()); err != nil {
		return false
	} else if resp.StatusCode != 200 {
		return false
	}
	return true
}

// test functions
func noop(*testing.T, dktest.ContainerInfo) {}

func TestRun(t *testing.T) {
	dktest.Run(t, testImage, dktest.Options{}, noop)
}

func TestRunParallel(t *testing.T) {
	numTests := 10
	for i := 0; i < numTests; i++ {
		t.Run("", func(t *testing.T) {
			t.Parallel()
			dktest.Run(t, "alpine:3.8", dktest.Options{}, noop)
		})
	}
}

func TestRunWithNetwork(t *testing.T) {
	dktest.Run(t, testNetworkImage, dktest.Options{ReadyFunc: nginxReady, PortRequired: true}, noop)
}

func TestRunWithNetworkPortBinding(t *testing.T) {
	port, err := nat.NewPort("tcp", "80")
	if err != nil {
		t.Fatal("Invalid port:", err)
	}

	dktest.Run(t, testNetworkImage, dktest.Options{ReadyFunc: nginxReady, PortRequired: true,
		PortBindings: nat.PortMap{port: []nat.PortBinding{{HostPort: "8181"}}}}, noop)
}
