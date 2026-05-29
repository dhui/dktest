package dktest

import (
	"context"
	"io"
	"net/netip"
	"testing"
	"time"

	"github.com/dhui/dktest/mockdockerclient"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
)

const (
	imageName = "dktestFakeImageName"
)

var (
	containerInfo = ContainerInfo{}
)

// ready functions
func alwaysReady(context.Context, ContainerInfo) bool { return true }
func neverReady(context.Context, ContainerInfo) bool  { return false }

func testErr(t *testing.T, err error, expectErr bool) {
	t.Helper()
	if err == nil && expectErr {
		t.Error("Expected an error but didn't get one")
	} else if err != nil && !expectErr {
		t.Error("Got unexpected error:", err)
	}
}

func TestPullImage(t *testing.T) {
	successReader := mockdockerclient.MockReader{Err: io.EOF}

	testCases := []struct {
		name      string
		client    mockdockerclient.ImageAPIClient
		platform  string
		expectErr bool
	}{
		{name: "success", client: mockdockerclient.ImageAPIClient{
			PullResp: &mockdockerclient.MockImagePullResponse{ReadCloser: mockdockerclient.MockReadCloser{MockReader: successReader}}}, expectErr: false},
		{name: "with specific platform", client: mockdockerclient.ImageAPIClient{
			PullResp: &mockdockerclient.MockImagePullResponse{ReadCloser: mockdockerclient.MockReadCloser{MockReader: successReader}}},
			platform: "linux/amd64", expectErr: false},
		{name: "pull error", client: mockdockerclient.ImageAPIClient{}, expectErr: true},
		{name: "read error", client: mockdockerclient.ImageAPIClient{
			PullResp: &mockdockerclient.MockImagePullResponse{ReadCloser: mockdockerclient.MockReadCloser{
				MockReader: mockdockerclient.MockReader{Err: mockdockerclient.Err},
			}}}, expectErr: false},
		{name: "close error", client: mockdockerclient.ImageAPIClient{
			PullResp: &mockdockerclient.MockImagePullResponse{ReadCloser: mockdockerclient.MockReadCloser{
				MockReader: successReader,
				MockCloser: mockdockerclient.MockCloser{Err: mockdockerclient.Err},
			}}}, expectErr: false},
	}

	ctx := context.Background()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := tc.client
			err := pullImage(ctx, t, &c, "", imageName, tc.platform)
			testErr(t, err, tc.expectErr)
		})
	}
}

func TestRunImage(t *testing.T) {
	portBindingsNoIP := network.PortMap{
		network.MustParsePort("80/tcp"): []network.PortBinding{{HostPort: "8181"}},
	}
	portBindingsIPZeros := network.PortMap{
		network.MustParsePort("80/tcp"): []network.PortBinding{{HostIP: netip.MustParseAddr("0.0.0.0"), HostPort: "8181"}},
	}
	portBindingsDiffIP := network.PortMap{
		network.MustParsePort("80/tcp"): []network.PortBinding{{HostIP: netip.MustParseAddr("10.0.0.1"), HostPort: "8181"}},
	}

	successCreateResp := &client.ContainerCreateResult{}
	successInspectResp := &client.ContainerInspectResult{}
	successInspectRespWithPortBindingNoIP := &client.ContainerInspectResult{Container: container.InspectResponse{NetworkSettings: &container.NetworkSettings{
		Ports: portBindingsNoIP,
	}}}
	successInspectRespWithPortBindingIPZeros := &client.ContainerInspectResult{Container: container.InspectResponse{NetworkSettings: &container.NetworkSettings{
		Ports: portBindingsIPZeros,
	}}}
	successInspectRespWithPortBindingDiffIP := &client.ContainerInspectResult{Container: container.InspectResponse{NetworkSettings: &container.NetworkSettings{
		Ports: portBindingsDiffIP,
	}}}

	testCases := []struct {
		name      string
		client    mockdockerclient.ContainerAPIClient
		opts      Options
		expectErr bool
	}{
		{name: "success", client: mockdockerclient.ContainerAPIClient{
			CreateResp: successCreateResp, InspectResp: successInspectResp}, expectErr: false},
		{name: "success - with port binding no ip", client: mockdockerclient.ContainerAPIClient{
			CreateResp: successCreateResp, InspectResp: successInspectRespWithPortBindingNoIP}, expectErr: false},
		{name: "success - with port binding ip 0.0.0.0", client: mockdockerclient.ContainerAPIClient{
			CreateResp: successCreateResp, InspectResp: successInspectRespWithPortBindingIPZeros}, expectErr: false},
		{name: "success - with port binding diff ip", client: mockdockerclient.ContainerAPIClient{
			CreateResp: successCreateResp, InspectResp: successInspectRespWithPortBindingDiffIP}, expectErr: false},
		{name: "create error", client: mockdockerclient.ContainerAPIClient{InspectResp: successInspectResp},
			expectErr: true},
		{name: "start error", client: mockdockerclient.ContainerAPIClient{
			CreateResp: successCreateResp, StartErr: mockdockerclient.Err, InspectResp: successInspectResp,
		}, expectErr: true},
		{name: "inspect error", client: mockdockerclient.ContainerAPIClient{CreateResp: successCreateResp},
			opts: Options{PortRequired: true}, expectErr: true},
		{name: "no network settings error", client: mockdockerclient.ContainerAPIClient{
			CreateResp: successCreateResp, InspectResp: successInspectResp}, opts: Options{PortRequired: true},
			expectErr: true},
	}

	ctx := context.Background()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := tc.client
			_, err := runImage(ctx, t, &client, imageName, tc.opts)
			testErr(t, err, tc.expectErr)
		})
	}
}

func TestStopContainer(t *testing.T) {
	successReadCloser := mockdockerclient.MockReadCloser{MockReader: mockdockerclient.MockReader{Err: io.EOF}}
	readCloserReadErr := mockdockerclient.MockReadCloser{
		MockReader: mockdockerclient.MockReader{Err: mockdockerclient.Err}}
	readCloserCloseErr := mockdockerclient.MockReadCloser{
		MockReader: mockdockerclient.MockReader{Err: io.EOF},
		MockCloser: mockdockerclient.MockCloser{Err: mockdockerclient.Err}}

	testCases := []struct {
		name   string
		client mockdockerclient.ContainerAPIClient
		log    bool
	}{
		{name: "success", client: mockdockerclient.ContainerAPIClient{}},
		{name: "success - log fetch error", client: mockdockerclient.ContainerAPIClient{}, log: true},
		{name: "success - log fetch success - read error",
			client: mockdockerclient.ContainerAPIClient{Logs: readCloserReadErr}, log: true},
		{name: "success - log fetch success - read success",
			client: mockdockerclient.ContainerAPIClient{Logs: successReadCloser}, log: true},
		{name: "success - log fetch success - close error",
			client: mockdockerclient.ContainerAPIClient{Logs: readCloserCloseErr}, log: true},
		{name: "stop error", client: mockdockerclient.ContainerAPIClient{StopErr: mockdockerclient.Err}},
		{name: "remove error", client: mockdockerclient.ContainerAPIClient{RemoveErr: mockdockerclient.Err}},
	}

	ctx := context.Background()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := tc.client
			stopContainer(ctx, t, &client, containerInfo, tc.log, tc.log)
		})
	}
}

func TestWaitContainerReady(t *testing.T) {
	canceledCtx, cancelFunc := context.WithCancel(context.Background())
	cancelFunc()

	testCases := []struct {
		name        string
		ctx         context.Context
		readyFunc   func(context.Context, ContainerInfo) bool
		expectReady bool
	}{
		{name: "nil readyFunc", ctx: canceledCtx, readyFunc: nil, expectReady: true},
		{name: "ready", ctx: context.Background(), readyFunc: alwaysReady, expectReady: true},
		{name: "not ready", ctx: canceledCtx, readyFunc: neverReady, expectReady: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if ready := waitContainerReady(tc.ctx, t, containerInfo, tc.readyFunc,
				time.Second); ready && !tc.expectReady {
				t.Error("Expected container to not be ready but it was")
			} else if !ready && tc.expectReady {
				t.Error("Expected container to ready but it wasn't")
			}
		})
	}
}
