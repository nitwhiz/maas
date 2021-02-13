package main

import "github.com/nitwhiz/maas/internal/server"

type BuildCmd struct {
}

func (c *BuildCmd) Run(ctx *Context) error {
	s, err := GetServerFromCwd()

	if err != nil {
		return err
	}

	err = s.Down(ctx.docker, server.DownOptions{Container: false})

	if _, ok := err.(*server.NoContainerFoundError); !ok {
		return err
	}

	return s.Build(ctx.docker)
}
