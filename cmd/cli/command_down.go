package main

import (
	"github.com/nitwhiz/maas/internal/server"
)

type DownCmd struct {
	RemoveContainer bool `kong:"help='Remove server container.'"`
}

func (d *DownCmd) Run(ctx *Context) error {
	s, err := GetServerFromCwd()

	if err != nil {
		return err
	}

	return s.Down(ctx.docker, server.DownOptions{
		Container: d.RemoveContainer,
	})
}
