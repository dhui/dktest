package dktest

import (
	"time"
)

import (
	"github.com/docker/go-connections/nat"
)

// Options contains the configurable options for running tests in the docker image
type Options struct {
	Timeout   time.Duration
	ReadyFunc func(ContainerInfo) bool
	Env       map[string]string
	// If you prefer to specify your port bindings as a string, use nat.ParsePortSpecs()
	PortBindings nat.PortMap
	PortRequired bool
}

func (o *Options) init() {
	if o.Timeout <= 0 {
		o.Timeout = DefaultTimeout
	}
}

func (o *Options) env() []string {
	env := make([]string, 0, len(o.Env))
	for k, v := range o.Env {
		env = append(env, k+"="+v)
	}
	return env
}
