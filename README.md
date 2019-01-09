# dktest

[![Build Status](https://img.shields.io/travis/dhui/dktest/master.svg)](https://travis-ci.org/dhui/dktest) [![Code Coverage](https://img.shields.io/codecov/c/github/dhui/dktest.svg)](https://codecov.io/gh/dhui/dktest) [![GoDoc](https://godoc.org/github.com/dhui/dktest?status.svg)](https://godoc.org/github.com/dhui/dktest) [![Go Report Card](https://goreportcard.com/badge/github.com/dhui/dktest)](https://goreportcard.com/report/github.com/dhui/dktest) [![GitHub Release](https://img.shields.io/github/release/dhui/dktest/all.svg)](https://github.com/dhui/dktest/releases) ![Supported Go versions](https://img.shields.io/badge/Go-1.11-lightgrey.svg)

`dktest` is short for **d**oc**k**er**test**.

`dktest` makes it stupidly easy to write integration tests in Go using Docker. Pulling images, starting containers, and cleaning up (even if your tests panic) is handled for you automatically!

## API

`Run()` is the workhorse

```golang
type ContainerInfo struct {
    ID        string
    Name      string
    ImageName string
    IP        string
    Port      string
}

type Options struct {
    Timeout   time.Duration
    ReadyFunc func(ContainerInfo) bool
    Env       map[string]string
    // If you prefer to specify your port bindings as a string, use nat.ParsePortSpecs()
    PortBindings nat.PortMap
    PortRequired bool
}

func Run(t *testing.T, imgName string, opts Options, testFunc func(*testing.T, ContainerInfo))
```

## Example Usage

```golang
import (
    "context"
    "testing"
)

import (
    "github.com/dhui/dktest"
    _ "github.com/lib/pq"
)

func pgReady(ctx context.Context, c dktest.ContainerInfo) bool {
    ip, port, err := c.FirstPort()
    if err != nil {
        return false
    }
    connStr := fmt.Sprintf("host=%s port=%s user=postgres dbname=postgres sslmode=disable", ip, port)
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return false
    }
    defer db.Close()
    return db.PingContext(ctx) == nil
}

func Test(t *testing.T) {
    dktest.Run(t, "postgres:alpine", dktest.Options{PortRequired: true, ReadyFunc: pgReady},
        func(t *testing.T, c dktest.ContainerInfo) {
        ip, port, err := c.FirstPort()
        if err != nil {
            t.Fatal(err)
        }
        connStr := fmt.Sprintf("host=%s port=%s user=postgres dbname=postgres sslmode=disable", ip, port)
        db, err := sql.Open("postgres", connStr)
        if err != nil {
            t.Fatal(err)
        }
        defer db.Close()
        if err := db.Ping(); err != nil {
            t.Fatal(err)
        }
        // Test using db
    })
}
```

For more examples, see the [docs](https://godoc.org/github.com/dhui/dktest).

## Debugging tests

Running `go test` with the `-v` option will display the container lifecycle log statements
along with the container ID.

### Short lived tests/containers

Run `go test` with the `-v` option and specify the `LogStdout` and/or `LogStderr` `Options`
to see the container's logs.

### Interactive tests/containers

Run `go test` with the `-v` option to get the container ID and check the container's logs with
`docker logs -f $CONTAINER_ID`.

## Cleaning up dangling containers

In the unlikely scenario where `dktest` leaves dangling containers,
you can find and removing them by using the `dktest` label:

```shell
# list dangling containers
$ docker ps -a --filter label=dktest
# stop dangling containers
$ docker ps --filter label=dktest | awk '{print $1}' | grep -v CONTAINER | xargs docker stop
# remove dangling containers
$ docker container prune --filter label=dktest
```

## Roadmap

* [x] Support multiple ports in `ContainerInfo`
* [ ] Use non-default network
* [ ] Add more `Options`
  * [ ] Volume mounts
  * [ ] Network config
* [ ] Support testing against multiple containers. It can be faked for now by nested/recursive `Run()` calls but that serializes the containers' startup time.

## Comparison to [dockertest](https://github.com/ory/dockertest)

### Why `dktest` is better

* Uses the [official Docker SDK](https://github.com/docker/docker)
  * [docker/docker](https://github.com/docker/docker) (aka [moby/moby](https://github.com/moby/moby)) uses [import path checking](https://golang.org/cmd/go/#hdr-Import_path_checking), so needs to be imported as `github.com/docker/docker`
* Designed to run in the Go testing environment
  * Smaller API surface
  * Running Docker containers are automatically cleaned up
* Has better test coverage
* Uses package management (Go modules) properly. e.g. not [manually vendored](https://github.com/ory/dockertest/pull/122)

### Why `dockertest` is better

* Has been around longer and API is more stable
* More options for configuring Docker containers
* Has more Github stars and contributors
