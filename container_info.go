package dktest

import (
	"fmt"
	"strconv"
)

import (
	"github.com/docker/go-connections/nat"
)

func mapHost(h string) string {
	switch h {
	case "", "0.0.0.0":
		return "127.0.0.1"
	default:
		return h
	}
}

func mapPort(portMap nat.PortMap, port nat.Port) (hostIP string, hostPort string, err error) {
	// Single port mapped
	portBindings, ok := portMap[port]
	if ok {
		for _, pb := range portBindings {
			return mapHost(pb.HostIP), pb.HostPort, nil
		}
	}

	// Search for port mapped in a range
	portInt := port.Int()
	proto := port.Proto()
	for p, portBindings := range portMap {
		if p.Proto() != proto {
			continue
		}
		start, end, err := p.Range()
		if err != nil {
			continue
		}
		if portInt < start || portInt > end {
			continue
		}
		offset := portInt - start
		if offset >= len(portBindings) {
			continue
		}
		pb := portBindings[offset]
		return mapHost(pb.HostIP), pb.HostPort, nil
	}

	return "", "", errNoPort
}

// firstPort gets the first port from the nat.PortMap.
// Since the underlying type is a map, the first port returned will not be consistent
func firstPort(portMap nat.PortMap, proto string) (hostIP string, hostPort string, err error) {
	for p, portBindings := range portMap {
		if p.Proto() != proto {
			continue
		}
		for _, pb := range portBindings {
			return mapHost(pb.HostIP), pb.HostPort, nil
		}
	}
	return "", "", errNoPort
}

func portMapToStrings(portMap nat.PortMap) []string {
	var portBindingStrs []string
	ports := make([]nat.Port, 0, len(portMap))
	for p := range portMap {
		ports = append(ports, p)
	}

	nat.SortPortMap(ports, portMap)

	for _, p := range ports {
		start, end, err := p.Range()
		if err != nil {
			continue
		}

		portBindings, ok := portMap[p]
		if !ok {
			// Very unlikely to happen
			continue
		}

		l := min(end-start+1, len(portBindings))
		proto := p.Proto()
		for i := 0; i < l; i++ {
			pb := portBindings[i]
			portBindingStrs = append(portBindingStrs, strconv.Itoa(start+i)+"/"+proto+" -> "+
				pb.HostIP+":"+pb.HostPort)
		}
	}
	return portBindingStrs
}

// ContainerInfo holds information about a running Docker container
type ContainerInfo struct {
	ID        string
	Name      string
	ImageName string
	Ports     nat.PortMap
}

// String gets the string representation for the ContainerInfo. This is intended for debugging purposes.
func (c ContainerInfo) String() string {
	return fmt.Sprintf("dktest.ContainerInfo{ID:%q, Name:%q, ImageName:%q, Ports:%v}", c.ID, c.Name, c.ImageName,
		portMapToStrings(c.Ports))
}

// Port gets the specified published/bound/mapped TCP port
func (c ContainerInfo) Port(containerPort uint16) (hostIP string, hostPort string, err error) {
	port, err := nat.NewPort("tcp", strconv.Itoa(int(containerPort)))
	if err != nil {
		return "", "", err
	}
	return mapPort(c.Ports, port)
}

// UDPPort gets the specified published/bound/mapped UDP port
func (c ContainerInfo) UDPPort(containerPort uint16) (hostIP string, hostPort string, err error) {
	port, err := nat.NewPort("udp", strconv.Itoa(int(containerPort)))
	if err != nil {
		return "", "", err
	}
	return mapPort(c.Ports, port)
}

// FirstPort gets the first published/bound/mapped TCP port. It is always safer to use Port().
// This provided as a convenience method and should only be used with Docker images that only expose a single port.
// If the Docker image exposes multiple ports, then the "first" port will not always be the same.
func (c ContainerInfo) FirstPort() (hostIP string, hostPort string, err error) {
	return firstPort(c.Ports, "tcp")
}

// FirstUDPPort gets the first published/bound/mapped UDP port. It is always safer to use UDPPort().
// This provided as a convenience method and should only be used with Docker images that only expose a single port.
// If the Docker image exposes multiple ports, then the "first" port will not always be the same.
func (c ContainerInfo) FirstUDPPort() (hostIP string, hostPort string, err error) {
	return firstPort(c.Ports, "udp")
}
