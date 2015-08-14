package main

import (
	"bytes"
	"os"

	"github.com/fsouza/go-dockerclient"

	"github.com/wantedly/risu/schema"
)

const (
	DefaultDockerEndpoint = "unix:///var/run/docker.sock"
)

func dockerPush(build schema.Build) error {

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
	// push build image
	pushOpts := docker.PushImageOptions{
		Name:         build.Name,
		Tag:          "latest",
		Registry:     "quay.io",
		OutputStream: outputbuf,
	}
	authConfig := docker.AuthConfiguration{
		Username:      "spesnova",
		Password:      "XXXXXXXXXXXXXXXXXXXXX",
		Email:         "spesnova@gmail.com",
		ServerAddress: "quay.io",
	}
	if err := client.PushImage(pushOpts, authConfig); err != nil {
		return err
	}
	os.Stdout.Write(outputbuf.Bytes())
	return nil
}
