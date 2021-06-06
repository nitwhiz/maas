package main

import (
	"fmt"
	"github.com/nitwhiz/maas/pkg/server"
)

type UpCmd struct {
	Follow bool `kong:"short='f',help='Follow container log output'"`
}

func (c *UpCmd) Run(ctx *Context) error {
	s, err := GetServerFromCwd()

	if err != nil {
		return err
	}

	upErr := s.Up(ctx.docker)

	if _, ok := upErr.(*server.NoContainerFoundError); ok {
		fmt.Println("no container found. building server container ...")

		err := s.Build(ctx.docker, server.BuildOptions{
			PullPrinter: printPullProgress,
		})

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

		err = s.Build(ctx.docker, server.BuildOptions{
			PullPrinter: printPullProgress,
		})

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
			Tail:       "25",
		}

		err = logsCmd.Run(ctx)

		if err != nil {
			return err
		}
	}

	return upErr
}
