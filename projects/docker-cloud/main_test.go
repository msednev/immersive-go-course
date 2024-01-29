package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
)
var pool *dockertest.Pool
var resource *dockertest.Resource

func TestMain(m *testing.M) {
	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not construct pool: %v", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("could not connect to Docker: %v", err)
	}

	resource, err = pool.BuildAndRun("docker-cloud", "./", nil)
	if err != nil {
		log.Fatalf("could not start resource: %v", err)
	}

	exitCode := m.Run()
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("cannot cleanup: %v", err)
	}

	os.Exit(exitCode)

}

func TestPing(t *testing.T) {
	var err error

	var resp *http.Response
	err = pool.Retry(func() error {
		var err error
		t.Log(resource.GetPort("80/tcp"))
		resp, err = http.Get(fmt.Sprintf("http://localhost:%s/ping", resource.GetPort("80/tcp")))
		if err != nil {
			t.Log("container not ready, retrying...")
			return err
		}
		return nil
	})

	if err != nil {
		t.Fatalf("response error: %s", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body: %s", err)
	}

	want := "Hello!"
	if got := string(body); want != got {
		t.Fatalf("bad response: wanted %s, got %s", want, got)
	}

}