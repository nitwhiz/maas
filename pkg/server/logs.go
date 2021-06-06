package server

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
)

type LogOptions struct {
	Timestamps bool
	Follow     bool
	Since      string
	Until      string
	Tail       string
}

// GetLogs retrieves the logs of a server container
func (s *Server) GetLogs(docker *client.Client, opts LogOptions) (io.ReadCloser, error) {
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
		Since:      opts.Since,
		Until:      opts.Until,
		Tail:       opts.Tail,
	})

	if err != nil {
		return out, err
	}

	return out, nil
}
