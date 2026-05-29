package dktest

import (
	"context"
	"fmt"
	"net/netip"
	"strings"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/api/types/network"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

// Options contains the configurable options for running tests in the docker image
type Options struct {
	// PullTimeout is the timeout used when pulling images
	PullTimeout time.Duration
	// PullRegistryAuth is the base64 encoded credentials for the registry
	PullRegistryAuth string
	// Timeout is the timeout used when starting a container and checking if it's ready
	Timeout time.Duration
	// ReadyTimeout is the timeout used for each container ready check.
	// e.g. each invocation of the ReadyFunc
	ReadyTimeout time.Duration
	// CleanupTimeout is the timeout used when stopping and removing a container
	CleanupTimeout time.Duration
	// CleanupImage specifies whether or not the image should be removed after the test run.
	// If the image is used by multiple tests, you'll want to cleanup the image yourself.
	CleanupImage bool
	ReadyFunc    func(context.Context, ContainerInfo) bool
	Env          map[string]string
	Entrypoint   []string
	Cmd          []string
	// If you prefer to specify your port bindings as a string, use nat.ParsePortSpecs()
	PortBindings nat.PortMap
	PortRequired bool
	LogStdout    bool
	LogStderr    bool
	ShmSize      int64
	Volumes      []string
	Mounts       []mount.Mount
	Hostname     string
	// Platform specifies the platform of the docker image that is pulled, e.g. "linux/amd64"
	Platform     string
	ExposedPorts nat.PortSet
}

// parsePlatform turns a platform string into a [v1.Platform]. This was added
// during the migration to [github.com/moby/moby] as a way to keep the public [Options] api static
// platform string convention is $os/$arch/$variant
func parsePlatform(p string) (v1.Platform, error) {
	splitPlat := strings.Split(p, "/")
	if len(splitPlat) < 2 {
		return v1.Platform{}, fmt.Errorf("invalid platform (%s): os and architecture must be provided $os/$architecture", p)
	}
	plat := v1.Platform{
		OS:           splitPlat[0],
		Architecture: splitPlat[1],
	}
	if len(splitPlat) == 3 {
		plat.Variant = splitPlat[2]
	}
	return plat, nil
}

// convertPortSet converts a [nat.PortSet] to a [network.PortSet]. This was added during the migration to
// [github.com/moby/moby] as a way to keep the public [Options] api static.
func convertPortSet(s nat.PortSet) network.PortSet {
	if len(s) == 0 {
		return nil
	}
	out := make(network.PortSet, len(s))
	for p := range s {
		out[network.MustParsePort(string(p))] = struct{}{}
	}
	return out
}

// convertPortMap converts a [nat.PortMap] to a [network.PortMap]. This was added during the migration to
// [github.com/moby/moby] as a way to keep the public [Options] api static.
func convertPortMap(m nat.PortMap) network.PortMap {
	if len(m) == 0 {
		return nil
	}
	out := make(network.PortMap, len(m))
	for p, bindings := range m {
		networkPort := network.MustParsePort(string(p))
		networkBindings := make([]network.PortBinding, len(bindings))
		for i, b := range bindings {
			var hostIP netip.Addr
			if b.HostIP != "" {
				hostIP, _ = netip.ParseAddr(b.HostIP)
			}
			networkBindings[i] = network.PortBinding{
				HostIP:   hostIP,
				HostPort: b.HostPort,
			}
		}
		out[networkPort] = networkBindings
	}
	return out
}

func (o *Options) init() {
	if o.PullTimeout <= 0 {
		o.PullTimeout = DefaultPullTimeout
	}
	if o.Timeout <= 0 {
		o.Timeout = DefaultTimeout
	}
	if o.ReadyTimeout <= 0 {
		o.ReadyTimeout = DefaultReadyTimeout
	}
	if o.CleanupTimeout <= 0 {
		o.CleanupTimeout = DefaultCleanupTimeout
	}
}

func (o *Options) volumes() map[string]struct{} {
	volumes := make(map[string]struct{})
	for _, v := range o.Volumes {
		volumes[v] = struct{}{}
	}
	return volumes
}

func (o *Options) env() []string {
	env := make([]string, 0, len(o.Env))
	for k, v := range o.Env {
		env = append(env, k+"="+v)
	}
	return env
}
