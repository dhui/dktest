package dktest_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/dhui/dktest"
	"github.com/docker/go-connections/nat"
)

const (
	testImage        = "alpine"
	testNetworkImage = "nginx:alpine"
)

// ready functions
func nginxReady(ctx context.Context, c dktest.ContainerInfo) bool {
	ip, port, err := c.FirstPort()
	if err != nil {
		return false
	}
	u := url.URL{Scheme: "http", Host: ip + ":" + port}
	fmt.Println(u.String())
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		fmt.Println("req err", err)
		return false
	}
	req = req.WithContext(ctx)
	if resp, err := http.DefaultClient.Do(req); err != nil {
		fmt.Println("do err:", err)
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

func TestRunContext(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		err := dktest.RunContext(context.Background(), t, testImage, dktest.Options{}, func(dktest.ContainerInfo) error {
			return nil
		})
		if err != nil {
			t.Fatal("failed", err)
		}
	})

	t.Run("test func returns error", func(t *testing.T) {
		var errForTest = errors.New("testFunc failed")
		err := dktest.RunContext(context.Background(), t, testImage, dktest.Options{}, func(dktest.ContainerInfo) error {
			return errForTest
		})
		if err == nil {
			t.Fatal("expected error")
		}
		if errors.Unwrap(err) != errForTest {
			t.Fatal("test func error not propagated with cause, got error:", err)
		}
	})
}

func TestRunParallel(t *testing.T) {
	numTests := 10
	for i := 0; i < numTests; i++ {
		t.Run("", func(t *testing.T) {
			t.Parallel()
			dktest.Run(t, testImage, dktest.Options{}, noop)
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
