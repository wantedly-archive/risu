package main

import (
	"archive/tar"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
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

func dockerCopy(build schema.Build) (string, error) {
	dockerEndpoint := os.Getenv("DOCKER_HOST")

	if dockerEndpoint == "" {
		dockerEndpoint = DefaultDockerEndpoint
	}

	client, err := docker.NewClient(dockerEndpoint)

	if err != nil {
		return "", err
	}

	// docker run
	container, err := client.CreateContainer(
		docker.CreateContainerOptions{
			Config: &docker.Config{
				Hostname:        "",
				Domainname:      "",
				User:            "",
				AttachStdin:     false,
				AttachStdout:    false,
				AttachStderr:    false,
				Tty:             false,
				OpenStdin:       false,
				StdinOnce:       false,
				Env:             nil,
				Cmd:             []string{"sleep", "3600"},
				Entrypoint:      []string{},
				Image:           build.ImageName,
				Labels:          nil,
				Volumes:         nil,
				WorkingDir:      "",
				NetworkDisabled: false,
				MacAddress:      "",
				ExposedPorts:    nil,
			},
			HostConfig: &docker.HostConfig{},
		})

	if err != nil {
		return "", err
	}

	if err = client.StartContainer(container.ID, &docker.HostConfig{}); err != nil {
		return "", err
	}

	// docker cp
	saveBaseDir := c.DefaultInflatedCacheDir + getCacheKey(build.SourceRepo) + "/"

	for _, cacheDirectory := range build.CacheDirectories {
		outputbuf := bytes.NewBuffer(nil)

		if err = client.CopyFromContainer(
			docker.CopyFromContainerOptions{
				Container:    container.ID,
				OutputStream: outputbuf,
				Resource:     cacheDirectory["container"],
			}); err != nil {
			return "", err
		}

		reader := tar.NewReader(outputbuf)

		for {
			header, err := reader.Next()

			if err == io.EOF {
				break
			}

			if err != nil {
				return "", err
			}

			buffer := new(bytes.Buffer)
			outPath := saveBaseDir + header.Name

			switch header.Typeflag {
			case tar.TypeDir:
				if _, err = os.Stat(outPath); err != nil {
					os.MkdirAll(outPath, 0755)
				}

			case tar.TypeReg, tar.TypeRegA:
				if _, err = io.Copy(buffer, reader); err != nil {
					return "", err
				}

				if err = ioutil.WriteFile(outPath, buffer.Bytes(), 0644); err != nil {
					return "", err
				}
			}
		}
	}

	cacheSavedDirectories := []string{saveBaseDir}

	if err = putCache(build, cacheSavedDirectories); err != nil {
		return "", err
	}

	// TODO: Stop & Remove container

	return saveBaseDir, nil
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

func putCache(build schema.Build, cacheSavedDirectories []string) error {
	cache := c.NewCache(os.Getenv("CACHE_BACKEND"))

	if err := cache.Put(getCacheKey(build.SourceRepo), cacheSavedDirectories); err != nil {
		return err
	}

	return nil
}

func getCacheKey(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))

	return hex.EncodeToString(hasher.Sum(nil))[0:12]
}
