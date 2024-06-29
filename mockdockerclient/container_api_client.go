package mockdockerclient

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"

	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

// ContainerAPIClient is a mock implementation of the Docker's client.ContainerAPIClient interface
type ContainerAPIClient struct {
	CreateResp  *container.CreateResponse
	StartErr    error
	StopErr     error
	RemoveErr   error
	InspectResp *types.ContainerJSON
	Logs        io.ReadCloser
}

var _ client.ContainerAPIClient = (*ContainerAPIClient)(nil)

// ContainerAttach is a mock implementation of Docker's client.ContainerAPIClient.ContainerAttach()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerAttach(context.Context, string,
	container.AttachOptions) (types.HijackedResponse, error) {
	return types.HijackedResponse{}, nil
}

// ContainerCommit is a mock implementation of Docker's client.ContainerAPIClient.ContainerCommit()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerCommit(context.Context, string,
	container.CommitOptions) (types.IDResponse, error) {
	return types.IDResponse{}, nil
}

// ContainerCreate is a mock implementation of Docker's client.ContainerAPIClient.ContainerCreate()
func (c *ContainerAPIClient) ContainerCreate(context.Context, *container.Config, *container.HostConfig,
	*network.NetworkingConfig, *v1.Platform, string) (container.CreateResponse, error) {
	if c.CreateResp == nil {
		return container.CreateResponse{}, Err
	}
	return *c.CreateResp, nil
}

// ContainerDiff is a mock implementation of Docker's client.ContainerAPIClient.ContainerDiff()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerDiff(context.Context,
	string) ([]container.FilesystemChange, error) {
	return nil, nil
}

// ContainerExecAttach is a mock implementation of Docker's client.ContainerAPIClient.ContainerExecAttach()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerExecAttach(context.Context, string,
	container.ExecStartOptions) (types.HijackedResponse, error) {
	return types.HijackedResponse{}, nil
}

// ContainerExecCreate is a mock implementation of Docker's client.ContainerAPIClient.ContainerExecCreate()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerExecCreate(context.Context, string,
	container.ExecOptions) (types.IDResponse, error) {
	return types.IDResponse{}, nil
}

// ContainerExecInspect is a mock implementation of Docker's client.ContainerAPIClient.ContainerExecInspect()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerExecInspect(context.Context,
	string) (container.ExecInspect, error) {
	return container.ExecInspect{}, nil
}

// ContainerExecResize is a mock implementation of Docker's client.ContainerAPIClient.ContainerExecResize()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerExecResize(context.Context, string,
	container.ResizeOptions) error {
	return nil
}

// ContainerExecStart is a mock implementation of Docker's client.ContainerAPIClient.ContainerExecStart()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerExecStart(context.Context, string,
	container.ExecStartOptions) error {
	return nil
}

// ContainerExport is a mock implementation of Docker's client.ContainerAPIClient.ContainerExport()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerExport(context.Context, string) (io.ReadCloser, error) {
	return nil, nil
}

// ContainerInspect is a mock implementation of Docker's client.ContainerAPIClient.ContainerInspect()
func (c *ContainerAPIClient) ContainerInspect(context.Context, string) (types.ContainerJSON, error) {
	if c.InspectResp == nil {
		return types.ContainerJSON{}, Err
	}
	return *c.InspectResp, nil
}

// ContainerInspectWithRaw is a mock implementation of Docker's client.ContainerAPIClient.ContainerInspectWithRaw()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerInspectWithRaw(context.Context, string,
	bool) (types.ContainerJSON, []byte, error) {
	return types.ContainerJSON{}, nil, nil
}

// ContainerKill is a mock implementation of Docker's client.ContainerAPIClient.ContainerKill()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerKill(context.Context, string, string) error {
	return nil
}

// ContainerList is a mock implementation of Docker's client.ContainerAPIClient.ContainerList()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerList(context.Context,
	container.ListOptions) ([]types.Container, error) {
	return nil, nil
}

// ContainerLogs is a mock implementation of Docker's client.ContainerAPIClient.ContainerLogs()
func (c *ContainerAPIClient) ContainerLogs(context.Context, string,
	container.LogsOptions) (io.ReadCloser, error) {
	if c.Logs == nil {
		return nil, Err
	}
	return c.Logs, nil
}

// ContainerPause is a mock implementation of Docker's client.ContainerAPIClient.ContainerPause()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerPause(context.Context, string) error { return nil }

// ContainerRemove is a mock implementation of Docker's client.ContainerAPIClient.ContainerRemove()
func (c *ContainerAPIClient) ContainerRemove(context.Context, string,
	container.RemoveOptions) error {
	return c.RemoveErr
}

// ContainerRename is a mock implementation of Docker's client.ContainerAPIClient.ContainerRename()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerRename(context.Context, string, string) error {
	return nil
}

// ContainerResize is a mock implementation of Docker's client.ContainerAPIClient.ContainerResize()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerResize(context.Context, string, container.ResizeOptions) error {
	return nil
}

// ContainerRestart is a mock implementation of Docker's client.ContainerAPIClient.ContainerRestart()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerRestart(context.Context, string, container.StopOptions) error {
	return nil
}

// ContainerStatPath is a mock implementation of Docker's client.ContainerAPIClient.ContainerStatPath()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerStatPath(context.Context, string,
	string) (container.PathStat, error) {
	return container.PathStat{}, nil
}

// ContainerStats is a mock implementation of Docker's client.ContainerAPIClient.ContainerStats()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerStats(context.Context, string,
	bool) (container.StatsResponseReader, error) {
	return container.StatsResponseReader{}, nil
}

// ContainerStatsOneShot is a mock implementation of Docker's client.ContainerAPIClient.ContainerStatsOneShot()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerStatsOneShot(context.Context, string) (container.StatsResponseReader, error) {
	return container.StatsResponseReader{}, nil
}

// ContainerStart is a mock implementation of Docker's client.ContainerAPIClient.ContainerStart()
func (c *ContainerAPIClient) ContainerStart(context.Context, string,
	container.StartOptions) error {
	return c.StartErr
}

// ContainerStop is a mock implementation of Docker's client.ContainerAPIClient.ContainerStop()
func (c *ContainerAPIClient) ContainerStop(context.Context, string, container.StopOptions) error {
	return c.StopErr
}

// ContainerTop is a mock implementation of Docker's client.ContainerAPIClient.ContainerTop()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerTop(context.Context, string,
	[]string) (container.ContainerTopOKBody, error) {
	return container.ContainerTopOKBody{}, nil
}

// ContainerUnpause is a mock implementation of Docker's client.ContainerAPIClient.ContainerUnpause()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerUnpause(context.Context, string) error {
	return nil
}

// ContainerUpdate is a mock implementation of Docker's client.ContainerAPIClient.ContainerUpdate()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerUpdate(context.Context, string,
	container.UpdateConfig) (container.ContainerUpdateOKBody, error) {
	return container.ContainerUpdateOKBody{}, nil
}

// ContainerWait is a mock implementation of Docker's client.ContainerAPIClient.ContainerWait()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerWait(context.Context, string,
	container.WaitCondition) (<-chan container.WaitResponse, <-chan error) {
	return nil, nil
}

// CopyFromContainer is a mock implementation of Docker's client.ContainerAPIClient.CopyFromContainer()
//
// TODO: properly implement
func (c *ContainerAPIClient) CopyFromContainer(context.Context, string, string) (io.ReadCloser,
	container.PathStat, error) {
	return nil, container.PathStat{}, nil
}

// CopyToContainer is a mock implementation of Docker's client.ContainerAPIClient.CopyToContainer()
//
// TODO: properly implement
func (c *ContainerAPIClient) CopyToContainer(context.Context, string, string, io.Reader,
	container.CopyToContainerOptions) error {
	return nil
}

// ContainersPrune is a mock implementation of Docker's client.ContainerAPIClient.ContainersPrune()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainersPrune(context.Context, filters.Args) (container.PruneReport, error) {
	return container.PruneReport{}, nil
}
