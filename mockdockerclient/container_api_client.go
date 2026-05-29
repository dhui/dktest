package mockdockerclient

import (
	"context"
	"io"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

// ContainerAPIClient is a mock implementation of the Docker's client.ContainerAPIClient interface
type ContainerAPIClient struct {
	CreateResp  *client.ContainerCreateResult
	StartErr    error
	StopErr     error
	RemoveErr   error
	InspectResp *client.ContainerInspectResult
	Logs        io.ReadCloser
}

var _ client.ContainerAPIClient = (*ContainerAPIClient)(nil)

// ContainerAttach is a mock implementation of Docker's client.ContainerAPIClient.ContainerAttach()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerAttach(context.Context, string, client.ContainerAttachOptions) (client.ContainerAttachResult, error) {
	return client.ContainerAttachResult{}, nil
}

// ContainerCommit is a mock implementation of Docker's client.ContainerAPIClient.ContainerCommit()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerCommit(context.Context, string,
	client.ContainerCommitOptions) (client.ContainerCommitResult, error) {
	return client.ContainerCommitResult{}, nil
}

// ContainerCreate is a mock implementation of Docker's client.ContainerAPIClient.ContainerCreate()
func (c *ContainerAPIClient) ContainerCreate(context.Context, client.ContainerCreateOptions) (client.ContainerCreateResult, error) {
	if c.CreateResp == nil {
		return client.ContainerCreateResult{}, Err
	}
	return *c.CreateResp, nil
}

// ContainerDiff is a mock implementation of Docker's client.ContainerAPIClient.ContainerDiff()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerDiff(context.Context,
	string, client.ContainerDiffOptions) (client.ContainerDiffResult, error) {
	return client.ContainerDiffResult{}, nil
}

// ContainerExecAttach is a mock implementation of Docker's client.ContainerAPIClient.ContainerExecAttach()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerExecAttach(context.Context, string, client.ExecAttachOptions) (client.ExecAttachResult, error) {
	return client.ExecAttachResult{}, nil
}

// ContainerExecCreate is a mock implementation of Docker's client.ContainerAPIClient.ContainerExecCreate()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerExecCreate(context.Context, string, client.ExecCreateOptions) (client.ExecCreateResult, error) {
	return client.ExecCreateResult{}, nil
}

// ContainerExecInspect is a mock implementation of Docker's client.ContainerAPIClient.ContainerExecInspect()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerExecInspect(context.Context, string, client.ExecInspectOptions) (client.ExecInspectResult, error) {
	return client.ExecInspectResult{}, nil
}

// ContainerExecResize is a mock implementation of Docker's client.ContainerAPIClient.ContainerExecResize()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerExecResize(context.Context, string, client.ExecResizeOptions) (client.ExecResizeResult, error) {
	return client.ExecResizeResult{}, nil
}

// ContainerExecStart is a mock implementation of Docker's client.ContainerAPIClient.ContainerExecStart()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerExecStart(context.Context, string, client.ExecStartOptions) (client.ExecStartResult, error) {
	return client.ExecStartResult{}, nil
}

// ContainerExport is a mock implementation of Docker's client.ContainerAPIClient.ContainerExport()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerExport(context.Context, string, client.ContainerExportOptions) (client.ContainerExportResult, error) {
	return nil, nil
}

// ContainerInspect is a mock implementation of Docker's client.ContainerAPIClient.ContainerInspect()
func (c *ContainerAPIClient) ContainerInspect(context.Context, string, client.ContainerInspectOptions) (client.ContainerInspectResult, error) {
	if c.InspectResp == nil {
		return client.ContainerInspectResult{}, Err
	}
	return *c.InspectResp, nil
}

// ContainerInspectWithRaw is a mock implementation of Docker's client.ContainerAPIClient.ContainerInspectWithRaw()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerInspectWithRaw(context.Context, string,
	bool) (container.InspectResponse, []byte, error) {
	return container.InspectResponse{}, nil, nil
}

// ContainerKill is a mock implementation of Docker's client.ContainerAPIClient.ContainerKill()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerKill(context.Context, string, client.ContainerKillOptions) (client.ContainerKillResult, error) {
	return client.ContainerKillResult{}, nil
}

// ContainerList is a mock implementation of Docker's client.ContainerAPIClient.ContainerList()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerList(context.Context, client.ContainerListOptions) (client.ContainerListResult, error) {
	return client.ContainerListResult{}, nil
}

// ContainerLogs is a mock implementation of Docker's client.ContainerAPIClient.ContainerLogs()
func (c *ContainerAPIClient) ContainerLogs(context.Context, string, client.ContainerLogsOptions) (client.ContainerLogsResult, error) {
	if c.Logs == nil {
		return nil, Err
	}
	return c.Logs, nil
}

// ContainerPause is a mock implementation of Docker's client.ContainerAPIClient.ContainerPause()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerPause(context.Context, string, client.ContainerPauseOptions) (client.ContainerPauseResult, error) {
	return client.ContainerPauseResult{}, nil
}

// ContainerRemove is a mock implementation of Docker's client.ContainerAPIClient.ContainerRemove()
func (c *ContainerAPIClient) ContainerRemove(context.Context, string, client.ContainerRemoveOptions) (client.ContainerRemoveResult, error) {
	return client.ContainerRemoveResult{}, c.RemoveErr
}

// ContainerRename is a mock implementation of Docker's client.ContainerAPIClient.ContainerRename()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerRename(context.Context, string, client.ContainerRenameOptions) (client.ContainerRenameResult, error) {
	return client.ContainerRenameResult{}, nil
}

// ContainerResize is a mock implementation of Docker's client.ContainerAPIClient.ContainerResize()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerResize(context.Context, string, client.ContainerResizeOptions) (client.ContainerResizeResult, error) {
	return client.ContainerResizeResult{}, nil
}

// ContainerRestart is a mock implementation of Docker's client.ContainerAPIClient.ContainerRestart()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerRestart(context.Context, string, client.ContainerRestartOptions) (client.ContainerRestartResult, error) {
	return client.ContainerRestartResult{}, nil
}

// ContainerStatPath is a mock implementation of Docker's client.ContainerAPIClient.ContainerStatPath()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerStatPath(context.Context, string, client.ContainerStatPathOptions) (client.ContainerStatPathResult, error) {
	return client.ContainerStatPathResult{}, nil
}

// ContainerStats is a mock implementation of Docker's client.ContainerAPIClient.ContainerStats()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerStats(context.Context, string, client.ContainerStatsOptions) (client.ContainerStatsResult, error) {
	return client.ContainerStatsResult{}, nil
}

// ContainerStart is a mock implementation of Docker's client.ContainerAPIClient.ContainerStart()
func (c *ContainerAPIClient) ContainerStart(context.Context, string, client.ContainerStartOptions) (client.ContainerStartResult, error) {
	return client.ContainerStartResult{}, c.StartErr
}

// ContainerStop is a mock implementation of Docker's client.ContainerAPIClient.ContainerStop()
func (c *ContainerAPIClient) ContainerStop(context.Context, string, client.ContainerStopOptions) (client.ContainerStopResult, error) {
	return client.ContainerStopResult{}, c.StopErr
}

// ContainerTop is a mock implementation of Docker's client.ContainerAPIClient.ContainerTop()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerTop(context.Context, string, client.ContainerTopOptions) (client.ContainerTopResult, error) {
	return client.ContainerTopResult{}, nil
}

// ContainerUnpause is a mock implementation of Docker's client.ContainerAPIClient.ContainerUnpause()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerUnpause(context.Context, string, client.ContainerUnpauseOptions) (client.ContainerUnpauseResult, error) {
	return client.ContainerUnpauseResult{}, nil
}

// ContainerUpdate is a mock implementation of Docker's client.ContainerAPIClient.ContainerUpdate()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerUpdate(context.Context, string, client.ContainerUpdateOptions) (client.ContainerUpdateResult, error) {
	return client.ContainerUpdateResult{}, nil
}

// ContainerWait is a mock implementation of Docker's client.ContainerAPIClient.ContainerWait()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerWait(context.Context, string, client.ContainerWaitOptions) client.ContainerWaitResult {
	return client.ContainerWaitResult{}
}

// CopyFromContainer is a mock implementation of Docker's client.ContainerAPIClient.CopyFromContainer()
//
// TODO: properly implement
func (c *ContainerAPIClient) CopyFromContainer(context.Context, string, client.CopyFromContainerOptions) (client.CopyFromContainerResult, error) {
	return client.CopyFromContainerResult{}, nil
}

// CopyToContainer is a mock implementation of Docker's client.ContainerAPIClient.CopyToContainer()
//
// TODO: properly implement
func (c *ContainerAPIClient) CopyToContainer(context.Context, string, client.CopyToContainerOptions) (client.CopyToContainerResult, error) {
	return client.CopyToContainerResult{}, nil
}

// ContainerPrune is a mock implementation of Docker's client.ContainerAPIClient.ContainersPrune()
//
// TODO: properly implement
func (c *ContainerAPIClient) ContainerPrune(context.Context, client.ContainerPruneOptions) (client.ContainerPruneResult, error) {
	return client.ContainerPruneResult{}, nil
}
