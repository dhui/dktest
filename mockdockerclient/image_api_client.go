package mockdockerclient

import (
	"context"
	"io"
)

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/registry"
)

// ImageAPIClient is a mock implementation of the Docker's client.ImageAPIClient interface
type ImageAPIClient struct {
	PullResp io.ReadCloser
}

// ImageBuild is a mock implementation of Docker's client.ImageAPIClient.ImageBuild()
//
// TODO: properly implement
func (c *ImageAPIClient) ImageBuild(context.Context, io.Reader,
	types.ImageBuildOptions) (types.ImageBuildResponse, error) {
	return types.ImageBuildResponse{}, nil
}

// BuildCachePrune is a mock implementation of Docker's client.ImageAPIClient.BuildCachePrune()
//
// TODO: properly implement
func (c *ImageAPIClient) BuildCachePrune(context.Context,
	types.BuildCachePruneOptions) (*types.BuildCachePruneReport, error) {
	return nil, nil
}

// BuildCancel is a mock implementation of Docker's client.ImageAPIClient.BuildCancel()
//
// TODO: properly implement
func (c *ImageAPIClient) BuildCancel(context.Context, string) error { return nil }

// ImageCreate is a mock implementation of Docker's client.ImageAPIClient.ImageCreate()
//
// TODO: properly implement
func (c *ImageAPIClient) ImageCreate(context.Context, string,
	types.ImageCreateOptions) (io.ReadCloser, error) {
	return nil, nil
}

// ImageHistory is a mock implementation of Docker's client.ImageAPIClient.ImageHistory()
//
// TODO: properly implement
func (c *ImageAPIClient) ImageHistory(context.Context, string) ([]image.HistoryResponseItem, error) {
	return nil, nil
}

// ImageImport is a mock implementation of Docker's client.ImageAPIClient.ImageImport()
//
// TODO: properly implement
func (c *ImageAPIClient) ImageImport(context.Context, types.ImageImportSource, string,
	types.ImageImportOptions) (io.ReadCloser, error) {
	return nil, nil
}

// ImageInspectWithRaw is a mock implementation of Docker's client.ImageAPIClient.ImageInspectWithRaw()
//
// TODO: properly implement
func (c *ImageAPIClient) ImageInspectWithRaw(context.Context, string) (types.ImageInspect, []byte, error) {
	return types.ImageInspect{}, nil, nil
}

// ImageList is a mock implementation of Docker's client.ImageAPIClient.ImageList()
//
// TODO: properly implement
func (c *ImageAPIClient) ImageList(context.Context, types.ImageListOptions) ([]types.ImageSummary, error) {
	return nil, nil
}

// ImageLoad is a mock implementation of Docker's client.ImageAPIClient.ImageLoad()
//
// TODO: properly implement
func (c *ImageAPIClient) ImageLoad(context.Context, io.Reader, bool) (types.ImageLoadResponse, error) {
	return types.ImageLoadResponse{}, nil
}

// ImagePull is a mock implementation of Docker's client.ImageAPIClient.ImagePull()
func (c *ImageAPIClient) ImagePull(context.Context, string, types.ImagePullOptions) (io.ReadCloser, error) {
	if c.PullResp == nil {
		return nil, Err
	}
	return c.PullResp, nil
}

// ImagePush is a mock implementation of Docker's client.ImageAPIClient.ImagePush()
//
// TODO: properly implement
func (c *ImageAPIClient) ImagePush(context.Context, string, types.ImagePushOptions) (io.ReadCloser, error) {
	return nil, nil
}

// ImageRemove is a mock implementation of Docker's client.ImageAPIClient.ImageRemove()
//
// TODO: properly implement
func (c *ImageAPIClient) ImageRemove(context.Context, string,
	types.ImageRemoveOptions) ([]types.ImageDeleteResponseItem, error) {
	return nil, nil
}

// ImageSearch is a mock implementation of Docker's client.ImageAPIClient.ImageSearch()
//
// TODO: properly implement
func (c *ImageAPIClient) ImageSearch(context.Context, string,
	types.ImageSearchOptions) ([]registry.SearchResult, error) {
	return nil, nil
}

// ImageSave is a mock implementation of Docker's client.ImageAPIClient.ImageSave()
//
// TODO: properly implement
func (c *ImageAPIClient) ImageSave(context.Context, []string) (io.ReadCloser, error) {
	return nil, nil
}

// ImageTag is a mock implementation of Docker's client.ImageAPIClient.ImageTag()
//
// TODO: properly implement
func (c *ImageAPIClient) ImageTag(context.Context, string, string) error { return nil }

// ImagesPrune is a mock implementation of Docker's client.ImageAPIClient.ImagesPrune()
//
// TODO: properly implement
func (c *ImageAPIClient) ImagesPrune(context.Context, filters.Args) (types.ImagesPruneReport, error) {
	return types.ImagesPruneReport{}, nil
}
