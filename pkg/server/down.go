package server

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"time"
)

type DownOptions struct {
	Container bool
}

// Down stops a running server container, optionally removing the container
func (s *Server) Down(docker *client.Client, opts DownOptions) error {
	foundContainer, err := s.FindContainer(docker)

	if err != nil {
		return err
	}

	timeout := time.Second * 10

	if err := docker.ContainerStop(context.Background(), foundContainer.ID, &timeout); err != nil {
		return err
	}

	if opts.Container {
		if err := docker.ContainerRemove(context.Background(), foundContainer.ID, types.ContainerRemoveOptions{}); err != nil {
			return err
		}
	}

	return nil
}
