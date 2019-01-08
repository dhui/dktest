package dktest

import (
	"context"
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
)

var (
	// DefaultTimeout is the default timeout to use
	DefaultTimeout = time.Minute
)

const (
	label = "dktest"
)

func pullImage(ctx context.Context, lgr logger, dc client.ImageAPIClient, imgName string) error {
	lgr.Log("Pulling image:", imgName)
	// lgr.Log(dc.ImageList(ctx, types.ImageListOptions{All: true}))

	resp, err := dc.ImagePull(ctx, imgName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Close(); err != nil {
			lgr.Log("Failed to close image:", err)
		}
	}()

	// Log response
	b := strings.Builder{}
	if err := jsonmessage.DisplayJSONMessagesStream(resp, &b, 0, false, nil); err == nil {
		lgr.Log(b.String())
	} else {
		lgr.Log("Error parsing image pull response:", err)
	}

	return nil
}

func runImage(ctx context.Context, lgr logger, dc client.ContainerAPIClient, imgName string,
	opts Options) (ContainerInfo, error) {
	c := ContainerInfo{Name: genContainerName(), ImageName: imgName}
	createResp, err := dc.ContainerCreate(ctx, &container.Config{
		Image:      imgName,
		Labels:     map[string]string{label: "true"},
		Env:        opts.env(),
		Entrypoint: opts.Entrypoint,
		Cmd:        opts.Cmd,
	}, &container.HostConfig{
		PublishAllPorts: true,
		PortBindings:    opts.PortBindings,
	}, &network.NetworkingConfig{}, c.Name)
	if err != nil {
		return c, err
	}
	c.ID = createResp.ID
	lgr.Log("Created container.", c.String())

	if err := dc.ContainerStart(ctx, createResp.ID, types.ContainerStartOptions{}); err != nil {
		return c, err
	}

	if !opts.PortRequired {
		return c, nil
	}

	inspectResp, err := dc.ContainerInspect(ctx, c.ID)
	if err != nil {
		return c, err
	}
	if inspectResp.NetworkSettings == nil {
		return c, errNoNetworkSettings
	}
	c.Ports = inspectResp.NetworkSettings.Ports

	return c, nil
}

func stopContainer(ctx context.Context, lgr logger, dc client.ContainerAPIClient, c ContainerInfo,
	logStdout, logStderr bool) {
	if logStdout || logStderr {
		if logs, err := dc.ContainerLogs(ctx, c.ID, types.ContainerLogsOptions{
			Timestamps: true, ShowStdout: logStdout, ShowStderr: logStderr,
		}); err == nil {
			b, err := ioutil.ReadAll(logs)
			defer logs.Close()
			if err == nil {
				lgr.Log(string(b))
			} else {
				lgr.Log("Error reading container logs:", err)
			}
		} else {
			lgr.Log("Error fetching container logs:", err)
		}
	}

	if err := dc.ContainerStop(ctx, c.ID, nil); err != nil {
		lgr.Log("Error stopping container:", c.String(), "error:", err)
	}
	lgr.Log("Stopped container:", c.String())

	if err := dc.ContainerRemove(ctx, c.ID,
		types.ContainerRemoveOptions{RemoveVolumes: true, Force: true}); err != nil {
		lgr.Log("Error removing container:", c.String(), "error:", err)
	}
	lgr.Log("Removed container:", c.String())
}

func waitContainerReady(ctx context.Context, lgr logger, c ContainerInfo,
	readyFunc func(ContainerInfo) bool) bool {
	if readyFunc == nil {
		return true
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if readyFunc(c) {
				return true
			}
		case <-ctx.Done():
			lgr.Log("Container was never ready.", c.String())
			return false
		}
	}
}

// Run runs the given test function once the specified Docker image is running in a container
func Run(t *testing.T, imgName string, opts Options, testFunc func(*testing.T, ContainerInfo)) {
	dc, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.39"))
	if err != nil {
		t.Fatal("Failed to get Docker client:", err)
	}

	opts.init()
	ctx, cancelFunc := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancelFunc()

	if err := pullImage(ctx, t, dc, imgName); err != nil {
		t.Fatal("Failed to pull image:", imgName, "error:", err)
	}

	func() {
		c, err := runImage(ctx, t, dc, imgName, opts)
		if err != nil {
			t.Fatal("Failed to run image:", imgName, "error:", err)
		}
		defer stopContainer(ctx, t, dc, c, opts.LogStdout, opts.LogStderr)

		if waitContainerReady(ctx, t, c, opts.ReadyFunc) {
			testFunc(t, c)
		} else {
			t.Fatal("Container was never ready before timing out:", c.String())
		}
	}()
}
