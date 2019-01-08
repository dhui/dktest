package mockdockerclient

import (
	"context"
	"io"
	"time"
)

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
)

// ContainerAPIClient is a mock implementation of the Docker's client.ContainerAPIClient interface
type ContainerAPIClient struct {
	CreateResp  *container.ContainerCreateCreatedBody
	StartErr    error
	StopErr     error
	RemoveErr   error
	InspectResp *types.ContainerJSON
	Logs        io.ReadCloser
}

// ContainerAttach is a mock implementation of Docker's client.ContainerAPIClient.ContainerAttach()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerAttach(context.Context, string,
	types.ContainerAttachOptions) (types.HijackedResponse, error) {
	return types.HijackedResponse{}, nil
}

// ContainerCommit is a mock implementation of Docker's client.ContainerAPIClient.ContainerCommit()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerCommit(context.Context, string,
	types.ContainerCommitOptions) (types.IDResponse, error) {
	return types.IDResponse{}, nil
}

// ContainerCreate is a mock implementation of Docker's client.ContainerAPIClient.ContainerCreate()
func (c *ContainerAPIClient) ContainerCreate(context.Context, *container.Config, *container.HostConfig,
	*network.NetworkingConfig, string) (container.ContainerCreateCreatedBody, error) {
	if c.CreateResp == nil {
		return container.ContainerCreateCreatedBody{}, Err
	}
	return *c.CreateResp, nil
}

// ContainerDiff is a mock implementation of Docker's client.ContainerAPIClient.ContainerDiff()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerDiff(context.Context,
	string) ([]container.ContainerChangeResponseItem, error) {
	return nil, nil
}

// ContainerExecAttach is a mock implementation of Docker's client.ContainerAPIClient.ContainerExecAttach()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerExecAttach(context.Context, string,
	types.ExecStartCheck) (types.HijackedResponse, error) {
	return types.HijackedResponse{}, nil
}

// ContainerExecCreate is a mock implementation of Docker's client.ContainerAPIClient.ContainerExecCreate()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerExecCreate(context.Context, string,
	types.ExecConfig) (types.IDResponse, error) {
	return types.IDResponse{}, nil
}

// ContainerExecInspect is a mock implementation of Docker's client.ContainerAPIClient.ContainerExecInspect()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerExecInspect(context.Context,
	string) (types.ContainerExecInspect, error) {
	return types.ContainerExecInspect{}, nil
}

// ContainerExecResize is a mock implementation of Docker's client.ContainerAPIClient.ContainerExecResize()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerExecResize(context.Context, string,
	types.ResizeOptions) error {
	return nil
}

// ContainerExecStart is a mock implementation of Docker's client.ContainerAPIClient.ContainerExecStart()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerExecStart(context.Context, string,
	types.ExecStartCheck) error {
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
	types.ContainerListOptions) ([]types.Container, error) {
	return nil, nil
}

// ContainerLogs is a mock implementation of Docker's client.ContainerAPIClient.ContainerLogs()
func (c *ContainerAPIClient) ContainerLogs(context.Context, string,
	types.ContainerLogsOptions) (io.ReadCloser, error) {
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
	types.ContainerRemoveOptions) error {
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
func (c *ContainerAPIClient) ContainerResize(context.Context, string, types.ResizeOptions) error {
	return nil
}

// ContainerRestart is a mock implementation of Docker's client.ContainerAPIClient.ContainerRestart()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerRestart(context.Context, string, *time.Duration) error {
	return nil
}

// ContainerStatPath is a mock implementation of Docker's client.ContainerAPIClient.ContainerStatPath()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerStatPath(context.Context, string,
	string) (types.ContainerPathStat, error) {
	return types.ContainerPathStat{}, nil
}

// ContainerStats is a mock implementation of Docker's client.ContainerAPIClient.ContainerStats()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerStats(context.Context, string,
	bool) (types.ContainerStats, error) {
	return types.ContainerStats{}, nil
}

// ContainerStart is a mock implementation of Docker's client.ContainerAPIClient.ContainerStart()
func (c *ContainerAPIClient) ContainerStart(context.Context, string,
	types.ContainerStartOptions) error {
	return c.StartErr
}

// ContainerStop is a mock implementation of Docker's client.ContainerAPIClient.ContainerStop()
func (c *ContainerAPIClient) ContainerStop(context.Context, string, *time.Duration) error {
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
	container.WaitCondition) (<-chan container.ContainerWaitOKBody, <-chan error) {
	return nil, nil
}

// CopyFromContainer is a mock implementation of Docker's client.ContainerAPIClient.CopyFromContainer()
//
// TODO: properly implement
func (c *ContainerAPIClient) CopyFromContainer(context.Context, string, string) (io.ReadCloser,
	types.ContainerPathStat, error) {
	return nil, types.ContainerPathStat{}, nil
}

// CopyToContainer is a mock implementation of Docker's client.ContainerAPIClient.CopyToContainer()
//
// TODO: properly implement
func (c *ContainerAPIClient) CopyToContainer(context.Context, string, string, io.Reader,
	types.CopyToContainerOptions) error {
	return nil
}

// ContainersPrune is a mock implementation of Docker's client.ContainerAPIClient.ContainersPrune()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainersPrune(context.Context, filters.Args) (types.ContainersPruneReport, error) {
	return types.ContainersPruneReport{}, nil
}
