package dktest

import (
	"strconv"
	"testing"
)

import (
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
)

func TestMapHost(t *testing.T) {
	testCases := []struct {
		host               string
		expectedMappedHost string
	}{
		{host: "", expectedMappedHost: "127.0.0.1"},
		{host: "0.0.0.0", expectedMappedHost: "127.0.0.1"},
		{host: "0.0.0.1", expectedMappedHost: "0.0.0.1"},
		{host: "localhost", expectedMappedHost: "localhost"},
		{host: "not a host", expectedMappedHost: "not a host"},
	}

	for _, tc := range testCases {
		t.Run(tc.host, func(t *testing.T) {
			if h := mapHost(tc.host); h != tc.expectedMappedHost {
				t.Error("mapped host does not match expected:", h, "!=", tc.expectedMappedHost)
			}
		})
	}
}

func expectMapping(t *testing.T, ip, port string, err error, expectedIP, expectedPort string, expectedErr error) {
	if ip != expectedIP {
		t.Error("ip does not match expected:", ip, "!=", expectedIP)
	}
	if port != expectedPort {
		t.Error("port does not match expected:", port, "!=", expectedPort)
	}
	if err != expectedErr {
		t.Error("err does not match expected:", err, "!=", expectedErr)
	}
}

func TestMapPort(t *testing.T) {
	_, portMap, err := nat.ParsePortSpecs([]string{"9000:8000", "10000-11000:9000-10000"})
	if err != nil {
		t.Fatal(err)
	}
	portBindingsForRange := make([]nat.PortBinding, 0, 1000)
	for i := 10000; i <= 11000; i++ {
		portBindingsForRange = append(portBindingsForRange, nat.PortBinding{HostPort: strconv.Itoa(i)})
	}
	portMapWithRange := nat.PortMap{
		"9000-10000": portBindingsForRange,
	}

	testCases := []struct {
		name         string
		portMap      nat.PortMap
		port         nat.Port
		expectedIP   string
		expectedPort string
		expectedErr  error
	}{
		{name: "invalid search port", portMap: portMap, port: "", expectedErr: errNoPort},
		{name: "wrong protocol", portMap: portMap, port: "8000/udp", expectedErr: errNoPort},
		{name: "success - single port", portMap: portMap, port: "8000/tcp",
			expectedIP: "127.0.0.1", expectedPort: "9000"},
		{name: "port range - parsed", portMap: portMap, port: "9050/tcp",
			expectedIP: "127.0.0.1", expectedPort: "10050"},
		{name: "port range - manual - success", portMap: portMapWithRange, port: "9050/tcp",
			expectedIP: "127.0.0.1", expectedPort: "10050"},
		{name: "port range - manual - malformed range", portMap: nat.PortMap{"foobar": []nat.PortBinding{}},
			port: "9050/tcp", expectedErr: errNoPort},
		{name: "port range - manual - invalid range", portMap: nat.PortMap{"10000-9000": []nat.PortBinding{}},
			port: "9050/tcp", expectedErr: errNoPort},
		{name: "port range - manual - not in range", portMap: nat.PortMap{"2000-3000": []nat.PortBinding{}},
			port: "9050/tcp", expectedErr: errNoPort},
		{name: "port range - manual - invalid mapping", portMap: nat.PortMap{"9000-10000": []nat.PortBinding{}},
			port: "9050/tcp", expectedErr: errNoPort},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ip, port, err := mapPort(tc.portMap, tc.port)
			expectMapping(t, ip, port, err, tc.expectedIP, tc.expectedPort, tc.expectedErr)
		})
	}
}

func TestFirstPort(t *testing.T) {
	_, portMap, err := nat.ParsePortSpecs([]string{"9000:8000"})
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name         string
		portMap      nat.PortMap
		proto        string
		expectedIP   string
		expectedPort string
		expectedErr  error
	}{
		{name: "invalid proto", portMap: portMap, proto: "", expectedErr: errNoPort},
		{name: "wrong proto", portMap: portMap, proto: "udp", expectedErr: errNoPort},
		{name: "success", portMap: portMap, proto: "tcp", expectedIP: "127.0.0.1", expectedPort: "9000"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ip, port, err := firstPort(tc.portMap, tc.proto)
			expectMapping(t, ip, port, err, tc.expectedIP, tc.expectedPort, tc.expectedErr)
		})
	}
}

func TestPortMapToStrings(t *testing.T) {
	_, portMap, err := nat.ParsePortSpecs([]string{"9000-9010:8000-8010", "8000:7000"})
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name     string
		portMap  nat.PortMap
		expected []string
	}{
		{name: "malformed port range", portMap: nat.PortMap{"foobar": []nat.PortBinding{}}},
		{name: "success", portMap: portMap, expected: []string{
			"8010/tcp -> :9010",
			"8009/tcp -> :9009",
			"8008/tcp -> :9008",
			"8007/tcp -> :9007",
			"8006/tcp -> :9006",
			"8005/tcp -> :9005",
			"8004/tcp -> :9004",
			"8003/tcp -> :9003",
			"8002/tcp -> :9002",
			"8001/tcp -> :9001",
			"8000/tcp -> :9000",
			"7000/tcp -> :8000",
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := portMapToStrings(tc.portMap)
			assert.Equal(t, tc.expected, s)
		})
	}
}
