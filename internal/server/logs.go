package server

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
)

type LogOpts struct {
	Timestamps bool
	Follow     bool
}

// GetLogs retrieves the logs of a server container
func (s *Server) GetLogs(docker *client.Client, opts LogOpts) (io.ReadCloser, error) {
	var out io.ReadCloser

	foundContainer, err := s.FindContainer(docker)

	if err != nil {
		return out, err
	}

	out, err = docker.ContainerLogs(context.Background(), foundContainer.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: opts.Timestamps,
		Follow:     opts.Follow,
		Details:    false,
	})

	if err != nil {
		return out, err
	}

	return out, nil
}
