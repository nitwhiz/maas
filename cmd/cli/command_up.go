package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/nitwhiz/maas/internal/cursor"
	"github.com/nitwhiz/maas/internal/server"
	"io"
	"strings"
)

type UpCmd struct {
	Follow bool `kong:"short='f',help='Follow container logs immediately.'"`
}

type pullEvent struct {
	ID             string `json:"id"`
	Status         string `json:"status"`
	Error          string `json:"error,omitempty"`
	Progress       string `json:"progress,omitempty"`
	ProgressDetail struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"progressDetail"`
}

func ensureImageAvailability(ctx *Context, s *server.Server) error {
	images, err := ctx.docker.ImageList(context.Background(), types.ImageListOptions{})

	if err != nil {
		return err
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
		terminalCursor := cursor.Cursor{}
		layers := make([]string, 0)
		oldIndex := len(layers)

		var event *pullEvent

		reader, err := ctx.docker.ImagePull(context.Background(), s.VMConfig.Image, types.ImagePullOptions{})

		if err != nil {
			return err
		}

		decoder := json.NewDecoder(reader)

		terminalCursor.Hide()

		for {
			if err := decoder.Decode(&event); err != nil {
				if err == io.EOF {
					break
				}

				return err
			}

			imageID := event.ID

			if strings.HasPrefix(event.Status, "Digest:") || strings.HasPrefix(event.Status, "Status:") {
				fmt.Printf("%s\n", event.Status)
				continue
			}

			index := 0
			for i, v := range layers {
				if v == imageID {
					index = i + 1
					break
				}
			}

			if index > 0 {
				diff := index - oldIndex

				if diff > 1 {
					down := diff - 1
					terminalCursor.MoveDown(down)
				} else if diff < 1 {
					up := diff*(-1) + 1
					terminalCursor.MoveUp(up)
				}

				oldIndex = index
			} else {
				layers = append(layers, event.ID)
				diff := len(layers) - oldIndex

				if diff > 1 {
					terminalCursor.MoveDown(diff)
				}

				oldIndex = len(layers)
			}

			terminalCursor.ClearLine()

			if event.Status == "Pull complete" {
				fmt.Printf("%s: %s\n", event.ID, event.Status)
			} else {
				fmt.Printf("%s: %s %s\n", event.ID, event.Status, event.Progress)
			}
		}

		terminalCursor.Show()
	}

	return nil
}

func (c *UpCmd) Run(ctx *Context) error {
	s, err := GetServerFromCwd()

	if err != nil {
		return err
	}

	if err = ensureImageAvailability(ctx, &s); err != nil {
		return err
	}

	upErr := s.Up(ctx.docker)

	if _, ok := upErr.(*server.NoContainerFoundError); ok {
		fmt.Println("no container found. building server container ...")

		err := s.Build(ctx.docker)

		if err != nil {
			return err
		}

		fmt.Println("starting server container ...")

		upErr = s.Up(ctx.docker)
	} else if _, ok := upErr.(*server.ConfigMismatchError); ok {
		fmt.Println("config changed. rebuilding server container ...")

		err := s.Down(ctx.docker, server.DownOptions{Container: false})

		if err != nil {
			return err
		}

		err = s.Build(ctx.docker)

		if err != nil {
			return err
		}

		fmt.Println("starting server container ...")

		upErr = s.Up(ctx.docker)
	}

	if upErr == nil && c.Follow {
		logsCmd := LogsCmd{
			Timestamps: true,
			Follow:     true,
			Since:      "5m",
		}

		err = logsCmd.Run(ctx)

		if err != nil {
			return err
		}
	}

	return upErr
}
