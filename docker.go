package main

import (
	"bytes"
	"log"
	"os"

	"github.com/fsouza/go-dockerclient"

	"github.com/wantedly/risu/schema"
)

const (
	DefaultDockerEndpoint = "unix:///var/run/docker.sock"
)

func DockerBuild(build *schema.Build, cacheDir string) {
	var dockerEndpoint string

	if os.Getenv("DOCKER_HOST") != "" {
		dockerEndpoint = os.Getenv("DOCKER_HOST")
	}

	if dockerEndpoint == "" {
		dockerEndpoint = DefaultDockerEndpoint
	}

	client, err := docker.NewClient(dockerEndpoint)

	if err != nil {
		log.Fatal(err)
		return
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
		log.Fatal(err)
		return
	}

	os.Stdout.Write(outputbuf.Bytes())
	return
}
