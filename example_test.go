package dktest_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

import (
	"github.com/dhui/dktest"
	_ "github.com/lib/pq"
)

func Example_nginx() {
	dockerImageName := "nginx:alpine"
	readyFunc := func(c dktest.ContainerInfo) bool {
		ip, port, err := c.FirstPort()
		if err != nil {
			return false
		}
		u := url.URL{Scheme: "http", Host: ip + ":" + port}
		if resp, err := http.Get(u.String()); err != nil {
			return false
		} else if resp.StatusCode != 200 {
			return false
		}
		return true
	}

	// dktest.Run() should be used within a test
	dktest.Run(&testing.T{}, dockerImageName, dktest.Options{PortRequired: true, ReadyFunc: readyFunc},
		func(t *testing.T, c dktest.ContainerInfo) {
			// test code here
		})

	// Output:
}

func Example_postgres() {
	dockerImageName := "postgres:alpine"
	readyFunc := func(c dktest.ContainerInfo) bool {
		ip, port, err := c.FirstPort()
		if err != nil {
			return false
		}
		connStr := fmt.Sprintf("host=%s port=%s user=postgres dbname=postgres sslmode=disable", ip, port)
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			return false
		}
		defer db.Close() // nolint:errcheck
		return db.Ping() == nil
	}

	// dktest.Run() should be used within a test
	dktest.Run(&testing.T{}, dockerImageName, dktest.Options{PortRequired: true, ReadyFunc: readyFunc},
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
			defer db.Close() // nolint:errcheck
			if err := db.Ping(); err != nil {
				t.Fatal(err)
			}
			// Test using db
		})

	// Output:
}
