package main

import (
	"bytes"
	"os"
	"strings"

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

	nametag := strings.Split(build.Name, ":")

	outputbuf := bytes.NewBuffer(nil)
	// push build image
	pushOpts := docker.PushImageOptions{
		Name:         nametag[0],
		Tag:          nametag[1],
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
