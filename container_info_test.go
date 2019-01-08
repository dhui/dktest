package dktest_test

import (
	"testing"
)

import (
	"github.com/docker/go-connections/nat"
)

import (
	"github.com/dhui/dktest"
)

func getTestContainerInfo(t *testing.T) dktest.ContainerInfo {
	_, portMap, err := nat.ParsePortSpecs([]string{"8080:80", "3737:37/udp"})
	if err != nil {
		t.Fatal(err)
	}

	return dktest.ContainerInfo{
		ID:        "testID",
		Name:      "testContainerInfo",
		ImageName: "testImageName",
		Ports:     portMap,
	}
}

func TestContainerInfoString(t *testing.T) {
	// real test cases in internal tests for portMapToStrings()
	ci := getTestContainerInfo(t)
	expected := `dktest.ContainerInfo{ID:"testID", Name:"testContainerInfo", ImageName:"testImageName", Ports:[80/tcp -> :8080 37/udp -> :3737]}`
	if s := ci.String(); s != expected {
		t.Error("ContainerInfo String() doesn't match expected:", s, "!=", expected)
	}
}

// nolint:unparam
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

func TestContainerInfoPort(t *testing.T) {
	// real test cases in internal tests for mapPort()
	ci := getTestContainerInfo(t)
	ip, port, err := ci.Port(80)
	expectMapping(t, ip, port, err, "127.0.0.1", "8080", nil)
}

func TestContainerInfoUDPPort(t *testing.T) {
	// real test cases in internal tests for mapPort()
	ci := getTestContainerInfo(t)
	ip, port, err := ci.UDPPort(37)
	expectMapping(t, ip, port, err, "127.0.0.1", "3737", nil)
}

func TestContainerInfoFirstPort(t *testing.T) {
	// real test cases in internal tests for firstPort()
	ci := getTestContainerInfo(t)
	ip, port, err := ci.FirstPort()
	expectMapping(t, ip, port, err, "127.0.0.1", "8080", nil)
}

func TestContainerInfoFirstUDPPort(t *testing.T) {
	// real test cases in internal tests for firstPort()
	ci := getTestContainerInfo(t)
	ip, port, err := ci.FirstUDPPort()
	expectMapping(t, ip, port, err, "127.0.0.1", "3737", nil)
}
