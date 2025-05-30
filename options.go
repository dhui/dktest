package dktest

import (
	"context"
	"time"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
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
	// Platform specifies the platform of the docker image that is pulled.
	Platform     string
	ExposedPorts nat.PortSet
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
