package dktest

import (
	"net/netip"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/moby/moby/api/types/network"
	"github.com/stretchr/testify/assert"
)

func TestOptionsInit(t *testing.T) {
	timeout := 9 * time.Second

	testCases := []struct {
		name     string
		opts     Options
		expected Options
	}{
		{name: "default timeouts used", opts: Options{},
			expected: Options{
				PullTimeout:    DefaultPullTimeout,
				Timeout:        DefaultTimeout,
				ReadyTimeout:   DefaultReadyTimeout,
				CleanupTimeout: DefaultCleanupTimeout,
			},
		},
		{name: "default timeouts not used",
			opts: Options{
				PullTimeout:    timeout,
				Timeout:        timeout,
				ReadyTimeout:   timeout,
				CleanupTimeout: timeout,
			},
			expected: Options{
				PullTimeout:    timeout,
				Timeout:        timeout,
				ReadyTimeout:   timeout,
				CleanupTimeout: timeout,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.opts.init()
			assert.Equal(t, tc.expected, tc.opts, "Expected post-init Options to match expected")
		})
	}
}

func TestOptionsEnv(t *testing.T) {
	testCases := []struct {
		name        string
		env         map[string]string
		expectedEnv []string
	}{
		{name: "nil", env: nil, expectedEnv: nil},
		{name: "empty", env: nil, expectedEnv: nil},
		{name: "1 var", env: map[string]string{"foo": "bar"}, expectedEnv: []string{"foo=bar"}},
		{name: "1 var - empty value", env: map[string]string{"foo": ""}, expectedEnv: []string{"foo="}},
		{name: "1 var - empty key", env: map[string]string{"": "bar"}, expectedEnv: []string{"=bar"}},
		{name: "3 vars", env: map[string]string{"foo": "bar", "hello": "world", "dead": "beef"},
			expectedEnv: []string{"foo=bar", "hello=world", "dead=beef"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := Options{Env: tc.env}
			assert.ElementsMatch(t, tc.expectedEnv, opts.env(), "Options environment to match expected")
		})
	}
}

func TestConvertPortSet(t *testing.T) {
	testCases := []struct {
		name     string
		input    nat.PortSet
		expected network.PortSet
	}{
		{name: "nil", input: nil, expected: nil},
		{name: "empty", input: nat.PortSet{}, expected: nil},
		{name: "single port", input: nat.PortSet{
			"80/tcp": {},
		}, expected: network.PortSet{
			network.MustParsePort("80/tcp"): {},
		}},
		{name: "multiple ports", input: nat.PortSet{
			"80/tcp": {},
			"53/udp": {},
		}, expected: network.PortSet{
			network.MustParsePort("80/tcp"): {},
			network.MustParsePort("53/udp"): {},
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := convertPortSet(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestConvertPortMap(t *testing.T) {
	testCases := []struct {
		name     string
		input    nat.PortMap
		expected network.PortMap
	}{
		{name: "nil", input: nil, expected: nil},
		{name: "empty", input: nat.PortMap{}, expected: nil},
		{name: "single port no host ip", input: nat.PortMap{
			"80/tcp": {{HostPort: "8080"}},
		}, expected: network.PortMap{
			network.MustParsePort("80/tcp"): {{HostPort: "8080"}},
		}},
		{name: "with host ip", input: nat.PortMap{
			"80/tcp": {{HostIP: "10.0.0.1", HostPort: "8080"}},
		}, expected: network.PortMap{
			network.MustParsePort("80/tcp"): {{HostIP: netip.MustParseAddr("10.0.0.1"), HostPort: "8080"}},
		}},
		{name: "with 0.0.0.0", input: nat.PortMap{
			"80/tcp": {{HostIP: "0.0.0.0", HostPort: "8080"}},
		}, expected: network.PortMap{
			network.MustParsePort("80/tcp"): {{HostIP: netip.MustParseAddr("0.0.0.0"), HostPort: "8080"}},
		}},
		{name: "multiple ports and bindings", input: nat.PortMap{
			"80/tcp": {
				{HostPort: "8080"},
				{HostIP: "192.168.1.1", HostPort: "9090"},
			},
			"53/udp": {{HostPort: "5353"}},
		}, expected: network.PortMap{
			network.MustParsePort("80/tcp"): {
				{HostPort: "8080"},
				{HostIP: netip.MustParseAddr("192.168.1.1"), HostPort: "9090"},
			},
			network.MustParsePort("53/udp"): {{HostPort: "5353"}},
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := convertPortMap(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}
