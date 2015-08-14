package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsouza/go-dockerclient"

	c "github.com/wantedly/risu/cache"
	"github.com/wantedly/risu/schema"
)

const (
	DefaultDockerEndpoint = "unix:///var/run/docker.sock"
)

func dockerBuild(build schema.Build) error {
	clonePath := DefaultSourceBaseDir + build.SourceRepo

	if err := addCacheToSrcRepo(build, clonePath); err != nil {
		return err
	}

	dockerEndpoint := os.Getenv("DOCKER_HOST")

	if dockerEndpoint == "" {
		dockerEndpoint = DefaultDockerEndpoint
	}

	client, err := docker.NewClient(dockerEndpoint)

	if err != nil {
		return err
	}

	outputbuf := bytes.NewBuffer(nil)
	opts := docker.BuildImageOptions{
		Name:                build.ImageName,
		NoCache:             false,
		SuppressOutput:      true,
		RmTmpContainer:      true,
		ForceRmTmpContainer: true,
		Dockerfile:          build.Dockerfile,
		OutputStream:        outputbuf,
		ContextDir:          clonePath,
	}

	if err := client.BuildImage(opts); err != nil {
		return err
	}

	return nil
}

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

	dockerImageName := strings.Split(build.ImageName, ":")[0]
	dockerImageTag := strings.Split(build.ImageName, ":")[1]
	dockerRegistry := strings.Split(build.ImageName, "/")[0]

	outputbuf := bytes.NewBuffer(nil)
	pushOpts := docker.PushImageOptions{
		Name:         dockerImageName,
		Tag:          dockerImageTag,
		Registry:     dockerRegistry,
		OutputStream: outputbuf,
	}
	authConfig := docker.AuthConfiguration{
		Username:      os.Getenv("DOCKER_AUTH_USER_NAME"),
		Password:      os.Getenv("DOCKER_AUTH_USER_PASSWORD"),
		Email:         os.Getenv("DOCKER_AUTH_USER_EMAIL"),
		ServerAddress: dockerRegistry,
	}
	if err := client.PushImage(pushOpts, authConfig); err != nil {
		return err
	}
	return nil
}

func addCacheToSrcRepo(build schema.Build, clonePath string) error {
	cache := c.NewCache(os.Getenv("CACHE_BACKEND"))
	inflatedCachePath, err := cache.Get(getCacheKey(build.SourceRepo))

	if err != nil {
		return err
	}

	if inflatedCachePath != "" {
		for _, cacheDirectory := range build.CacheDirectories {
			cachePath := inflatedCachePath + string(filepath.Separator) + cacheDirectory["source"]
			sourcePath := clonePath + string(filepath.Separator) + cacheDirectory["source"]

			if _, err := os.Stat(sourcePath); err == nil {
				if e := os.RemoveAll(sourcePath); e != nil {
					return e
				}
			}

			if err := os.Rename(cachePath, sourcePath); err != nil {
				return err
			}
		}
	}

	return nil
}

func getCacheKey(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))

	return hex.EncodeToString(hasher.Sum(nil))[0:12]
}
