package server

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"os"
	"strconv"
	"strings"
)

type PullPrinterFunc func(reader *io.ReadCloser) error

type BuildOptions struct {
	PullPrinter PullPrinterFunc
}

func (s *Server) getContainerEnv() []string {
	env := []string{
		fmt.Sprintf("TYPE=%s", strings.ToUpper(s.Settings.Version.Type)),
		fmt.Sprintf("VERSION=%s", s.Settings.Version.GameVersion),
		"EULA=TRUE",
	}

	env = append(env, s.VMConfig.Environment...)

	return env
}

func (s *Server) ensureImageAvailability(docker *client.Client) (io.ReadCloser, error) {
	images, err := docker.ImageList(context.Background(), types.ImageListOptions{})

	if err != nil {
		return nil, err
	}

	imageExists := false

	for _, image := range images {
		for _, tag := range image.RepoTags {
			if tag == s.VMConfig.Image {
				imageExists = true
				break
			}
		}
	}

	if !imageExists {
		reader, err := docker.ImagePull(context.Background(), s.VMConfig.Image, types.ImagePullOptions{})

		if err != nil {
			return nil, err
		}

		return reader, nil
	}

	return nil, nil
}

// Build builds the server container
func (s *Server) Build(docker *client.Client, opts BuildOptions) error {
	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	env := s.getContainerEnv()

	configHash, err := s.GetConfigHash()

	if err != nil {
		return err
	}

	containerConfig := container.Config{
		Labels: map[string]string{
			"com.github.nitwhiz.maas.configPath": s.ConfigPath,
			"com.github.nitwhiz.maas.configHash": configHash,
		},
		Image: s.VMConfig.Image,
		Env:   env,
		ExposedPorts: nat.PortSet{
			"25565/tcp": struct{}{},
		},
	}

	hostConfig := container.HostConfig{
		AutoRemove: false,
		LogConfig: container.LogConfig{
			Type: "json-file",
			Config: map[string]string{
				"max-size": "10m",
				"max-file": "3",
			},
		},
		Binds: []string{
			fmt.Sprintf("%s:/data", cwd),
		},
		PortBindings: nat.PortMap{
			"25565/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: strconv.Itoa(s.VMConfig.ExposedPort),
				},
			},
		},
	}

	networkingConfig := network.NetworkingConfig{}

	if foundContainer, err := s.FindContainer(docker); err == nil {
		_ = docker.ContainerRemove(context.Background(), foundContainer.ID, types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		})
	}

	cwdSegments := strings.Split(cwd, "/")
	containerName := cwdSegments[len(cwdSegments)-1]

	pullReader, err := s.ensureImageAvailability(docker)

	if err != nil {
		return err
	}

	if opts.PullPrinter != nil {
		err := opts.PullPrinter(&pullReader)

		if err != nil {
			return err
		}
	}

	_, err = docker.ContainerCreate(context.Background(), &containerConfig, &hostConfig, &networkingConfig, nil, containerName)

	if err != nil {
		return err
	}

	return nil
}
