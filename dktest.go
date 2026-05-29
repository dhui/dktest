package dktest

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
	"github.com/moby/moby/client/pkg/jsonmessage"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

var (
	// DefaultPullTimeout is the default timeout used when pulling images
	DefaultPullTimeout = time.Minute
	// DefaultTimeout is the default timeout used when starting a container and checking if it's ready
	DefaultTimeout = time.Minute
	// DefaultReadyTimeout is the default timeout used for each container ready check.
	// e.g. each invocation of the ReadyFunc
	DefaultReadyTimeout = 2 * time.Second
	// DefaultCleanupTimeout is the default timeout used when stopping and removing a container
	DefaultCleanupTimeout = 15 * time.Second
)

const (
	label = "dktest"
)

func pullImage(ctx context.Context, lgr Logger, dc client.ImageAPIClient, registryAuth string, imgName string, platform string) error {
	lgr.Log("Pulling image:", imgName)
	// lgr.Log(dc.ImageList(ctx, types.ImageListOptions{All: true}))

	var platforms []v1.Platform
	if len(platform) > 0 {
		p, err := parsePlatform(platform)
		if err != nil {
			return err
		}
		platforms = []v1.Platform{p}
	}
	resp, err := dc.ImagePull(ctx, imgName, client.ImagePullOptions{
		Platforms:    platforms,
		RegistryAuth: registryAuth,
	})
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Close(); err != nil {
			lgr.Log("Failed to close image response:", err)
		}
	}()

	// Log response
	b := strings.Builder{}
	if err := jsonmessage.DisplayJSONMessagesStream(resp, &b, 0, false, nil); err == nil {
		lgr.Log("Image pull response:", b.String())
	} else {
		lgr.Log("Error parsing image pull response:", err)
	}

	return nil
}

func removeImage(ctx context.Context, lgr Logger, dc client.ImageAPIClient, imgName string) {
	lgr.Log("Removing image:", imgName)

	if _, err := dc.ImageRemove(ctx, imgName, client.ImageRemoveOptions{Force: true, PruneChildren: true}); err != nil {
		lgr.Log("Failed to remove image: ", err.Error())
	}
}

func runImage(ctx context.Context, lgr Logger, dc client.ContainerAPIClient, imgName string,
	opts Options) (ContainerInfo, error) {
	c := ContainerInfo{Name: genContainerName(), ImageName: imgName}
	createResp, err := dc.ContainerCreate(ctx,
		client.ContainerCreateOptions{
			Config: &container.Config{
				Image:        imgName,
				Labels:       map[string]string{label: "true"},
				Env:          opts.env(),
				Entrypoint:   opts.Entrypoint,
				Cmd:          opts.Cmd,
				Volumes:      opts.volumes(),
				Hostname:     opts.Hostname,
				ExposedPorts: convertPortSet(opts.ExposedPorts),
			},
			HostConfig: &container.HostConfig{
				PublishAllPorts: true,
				PortBindings:    convertPortMap(opts.PortBindings),
				ShmSize:         opts.ShmSize,
				Mounts:          opts.Mounts,
			},
			NetworkingConfig: &network.NetworkingConfig{},
			Name:             c.Name,
		})
	if err != nil {
		return c, err
	}
	c.ID = createResp.ID
	lgr.Log("Created container:", c.String())

	if _, err := dc.ContainerStart(ctx, createResp.ID, client.ContainerStartOptions{}); err != nil {
		return c, err
	}
	lgr.Log("Started container:", c.String())

	if !opts.PortRequired {
		return c, nil
	}

	inspectResp, err := dc.ContainerInspect(ctx, c.ID, client.ContainerInspectOptions{})
	if err != nil {
		return c, err
	}
	lgr.Log("Inspected container:", c.String())

	if inspectResp.Container.NetworkSettings == nil {
		return c, errNoNetworkSettings
	}
	c.Ports = toNatPortMap(inspectResp.Container.NetworkSettings.Ports)

	return c, nil
}

func stopContainer(ctx context.Context, lgr Logger, dc client.ContainerAPIClient, c ContainerInfo,
	logStdout, logStderr bool) {
	if logStdout || logStderr {
		if logs, err := dc.ContainerLogs(ctx, c.ID, client.ContainerLogsOptions{
			Timestamps: true, ShowStdout: logStdout, ShowStderr: logStderr,
		}); err == nil {
			b, err := io.ReadAll(logs)
			defer func() {
				if err := logs.Close(); err != nil {
					lgr.Log("Error closing logs:", err)
				}
			}()
			if err == nil {
				lgr.Log("Container logs:", string(b))
			} else {
				lgr.Log("Error reading container logs:", err)
			}
		} else {
			lgr.Log("Error fetching container logs:", err)
		}
	}

	if _, err := dc.ContainerStop(ctx, c.ID, client.ContainerStopOptions{}); err != nil {
		lgr.Log("Error stopping container:", c.String(), "error:", err)
	}
	lgr.Log("Stopped container:", c.String())

	if _, err := dc.ContainerRemove(ctx, c.ID,
		client.ContainerRemoveOptions{RemoveVolumes: true, Force: true}); err != nil {
		lgr.Log("Error removing container:", c.String(), "error:", err)
	}
	lgr.Log("Removed container:", c.String())
}

func waitContainerReady(ctx context.Context, lgr Logger, c ContainerInfo,
	readyFunc func(context.Context, ContainerInfo) bool, readyTimeout time.Duration) bool {
	if readyFunc == nil {
		return true
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ready := func() bool {
				readyCtx, canceledFunc := context.WithTimeout(ctx, readyTimeout)
				defer canceledFunc()
				return readyFunc(readyCtx, c)
			}()

			if ready {
				return true
			}
		case <-ctx.Done():
			lgr.Log("Container was never ready:", c.String())
			return false
		}
	}
}

// Run runs the given test function once the specified Docker image is running in a container
func Run(t *testing.T, imgName string, opts Options, testFunc func(*testing.T, ContainerInfo)) {
	err := RunContext(context.Background(), t, imgName, opts, func(containerInfo ContainerInfo) error {
		testFunc(t, containerInfo)
		return nil
	})
	if err != nil {
		t.Fatal("Failed:", err)
	}
}

// RunContext is similar to Run, but takes a parent context and returns an error and doesn't rely on a testing.T.
func RunContext(ctx context.Context, logger Logger, imgName string, opts Options, testFunc func(ContainerInfo) error) (retErr error) {
	dc, err := client.New(client.FromEnv)
	if err != nil {
		return fmt.Errorf("error getting Docker client: %w", err)
	}
	defer func() {
		if err := dc.Close(); err != nil && retErr == nil {
			retErr = fmt.Errorf("error closing Docker client: %w", err)
		}
	}()

	opts.init()
	pullCtx, pullTimeoutCancelFunc := context.WithTimeout(ctx, opts.PullTimeout)
	defer pullTimeoutCancelFunc()

	if err := pullImage(pullCtx, logger, dc, opts.PullRegistryAuth, imgName, opts.Platform); err != nil {
		return fmt.Errorf("error pulling image: %v error: %w", imgName, err)
	}

	return func() error {
		runCtx, runTimeoutCancelFunc := context.WithTimeout(ctx, opts.Timeout)
		defer runTimeoutCancelFunc()

		c, err := runImage(runCtx, logger, dc, imgName, opts)
		if err != nil {
			return fmt.Errorf("error running image: %v error: %w", imgName, err)
		}
		defer func() {
			stopCtx, stopTimeoutCancelFunc := context.WithTimeout(ctx, opts.CleanupTimeout)
			defer stopTimeoutCancelFunc()
			stopContainer(stopCtx, logger, dc, c, opts.LogStdout, opts.LogStderr)
			if opts.CleanupImage {
				removeImage(stopCtx, logger, dc, imgName)
			}
		}()

		if waitContainerReady(runCtx, logger, c, opts.ReadyFunc, opts.ReadyTimeout) {
			if err := testFunc(c); err != nil {
				return fmt.Errorf("error running test func: %w", err)
			}
		} else {
			return fmt.Errorf("timed out waiting for container to get ready: %v", c.String())
		}

		return nil
	}()
}
