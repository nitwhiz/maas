package main

import (
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/nitwhiz/maas/internal/server"
	"os"
)

type LogsCmd struct {
	Timestamps bool `kong:"short='t',help='Log timestamps.'"`
	Follow     bool `kong:"short='f',help='Follow Logs.'"`
}

func (c *LogsCmd) Run(ctx *Context) error {
	s, err := GetServerFromCwd()

	if err != nil {
		return err
	}

	out, err := s.GetLogs(ctx.docker, server.LogOpts{
		Timestamps: c.Timestamps,
		Follow:     c.Follow,
	})

	if err != nil {
		return err
	}

	defer out.Close()

	_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	if err != nil {
		return err
	}

	return nil
}
