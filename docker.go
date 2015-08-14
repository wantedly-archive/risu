package main

import (
	"bytes"
	"os"

	"github.com/fsouza/go-dockerclient"

	c "github.com/wantedly/risu/cache"
	"github.com/wantedly/risu/schema"
)

const (
	DefaultDockerEndpoint = "unix:///var/run/docker.sock"
)

func DockerBuild(build *schema.Build) error {
	cache := c.NewCache(os.Getenv("CACHE_BACKEND"))
	inflatedCachePath, err := cache.Get(build.ID.String())

	if err != nil {
		return err
	}

	if inflatedCachePath != "" {
		// put cache to repository
	}

	var dockerEndpoint string

	if os.Getenv("DOCKER_HOST") != "" {
		dockerEndpoint = os.Getenv("DOCKER_HOST")
	}

	if dockerEndpoint == "" {
		dockerEndpoint = DefaultDockerEndpoint
	}

	client, err := docker.NewClient(dockerEndpoint)

	if err != nil {
		return err
	}

	outputbuf := bytes.NewBuffer(nil)
	opts := docker.BuildImageOptions{
		Name:                build.Name,
		NoCache:             false,
		SuppressOutput:      true,
		RmTmpContainer:      true,
		ForceRmTmpContainer: true,
		Dockerfile:          build.Dockerfile,
		OutputStream:        outputbuf,
		ContextDir:          "", // TODO: Set `git clone` destination
	}

	if err := client.BuildImage(opts); err != nil {
		return err
	}

	return nil
}
