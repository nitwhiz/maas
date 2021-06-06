package main

import (
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/nitwhiz/maas/pkg/server"
	"io"
	"os"
)

type LogsCmd struct {
	Since      string `kong:"short='s',help='Show logs since timestamp (e.g. 2013-01-02T13:23:37Z) or relative (e.g. 42m for 42 minutes)'"`
	Until      string `kong:"short='u',help='Show logs before a timestamp (e.g. 2013-01-02T13:23:37Z) or relative (e.g. 42m for 42 minutes)'"`
	Timestamps bool   `kong:"short='t',help='Show timestamps'"`
	Follow     bool   `kong:"short='f',help='Follow log output'"`
	Tail       string `kong:"short='n',default='25',help='Number of lines to show from the end of the logs'"`
}

func (c *LogsCmd) Run(ctx *Context) error {
	s, err := GetServerFromCwd()

	if err != nil {
		return err
	}

	out, err := s.GetLogs(ctx.docker, server.LogOptions{
		Timestamps: c.Timestamps,
		Follow:     c.Follow,
		Since:      c.Since,
		Until:      c.Until,
		Tail:       c.Tail,
	})

	if err != nil {
		return err
	}

	defer func(out io.ReadCloser) {
		_ = out.Close()
	}(out)

	_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	if err != nil {
		return err
	}

	return nil
}
