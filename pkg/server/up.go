package server

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// Up starts the server container of this server
func (s *Server) Up(docker *client.Client) error {
	foundContainer, err := s.FindContainer(docker)

	if err != nil {
		return err
	}

	configHash, _ := s.GetConfigHash()

	if foundContainer.Labels["com.github.nitwhiz.maas.configHash"] != configHash {
		return &ConfigMismatchError{}
	}

	return docker.ContainerStart(context.Background(), foundContainer.ID, types.ContainerStartOptions{})
}
