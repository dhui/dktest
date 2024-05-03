package dktest

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
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

func pullImage(ctx context.Context, lgr Logger, dc client.ImageAPIClient, imgName, platform string) error {
	lgr.Log("Pulling image:", imgName)
	// lgr.Log(dc.ImageList(ctx, types.ImageListOptions{All: true}))

	resp, err := dc.ImagePull(ctx, imgName, types.ImagePullOptions{Platform: platform})
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

	if _, err := dc.ImageRemove(ctx, imgName, types.ImageRemoveOptions{Force: true, PruneChildren: true}); err != nil {
		lgr.Log("Failed to remove image: ", err.Error())
	}
}

func runImage(ctx context.Context, lgr Logger, dc client.ContainerAPIClient, imgName string,
	opts Options) (ContainerInfo, error) {
	c := ContainerInfo{Name: genContainerName(), ImageName: imgName}
	createResp, err := dc.ContainerCreate(ctx, &container.Config{
		Image:        imgName,
		Labels:       map[string]string{label: "true"},
		Env:          opts.env(),
		Entrypoint:   opts.Entrypoint,
		Cmd:          opts.Cmd,
		Volumes:      opts.volumes(),
		Hostname:     opts.Hostname,
		ExposedPorts: opts.ExposedPorts,
	}, &container.HostConfig{
		PublishAllPorts: true,
		PortBindings:    opts.PortBindings,
		ShmSize:         opts.ShmSize,
		Mounts:          opts.Mounts,
	}, &network.NetworkingConfig{},
		nil,
		c.Name)
	if err != nil {
		return c, err
	}
	c.ID = createResp.ID
	lgr.Log("Created container:", c.String())

	if err := dc.ContainerStart(ctx, createResp.ID, types.ContainerStartOptions{}); err != nil {
		return c, err
	}
	lgr.Log("Started container:", c.String())

	if !opts.PortRequired {
		return c, nil
	}

	inspectResp, err := dc.ContainerInspect(ctx, c.ID)
	if err != nil {
		return c, err
	}
	lgr.Log("Inspected container:", c.String())

	if inspectResp.NetworkSettings == nil {
		return c, errNoNetworkSettings
	}
	c.Ports = inspectResp.NetworkSettings.Ports

	return c, nil
}

func stopContainer(ctx context.Context, lgr Logger, dc client.ContainerAPIClient, c ContainerInfo,
	logStdout, logStderr bool) {
	if logStdout || logStderr {
		if logs, err := dc.ContainerLogs(ctx, c.ID, types.ContainerLogsOptions{
			Timestamps: true, ShowStdout: logStdout, ShowStderr: logStderr,
		}); err == nil {
			b, err := ioutil.ReadAll(logs)
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

	if err := dc.ContainerStop(ctx, c.ID, container.StopOptions{}); err != nil {
		lgr.Log("Error stopping container:", c.String(), "error:", err)
	}
	lgr.Log("Stopped container:", c.String())

	if err := dc.ContainerRemove(ctx, c.ID,
		types.ContainerRemoveOptions{RemoveVolumes: true, Force: true}); err != nil {
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
	dc, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.41"))
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

	if err := pullImage(pullCtx, logger, dc, imgName, opts.Platform); err != nil {
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
