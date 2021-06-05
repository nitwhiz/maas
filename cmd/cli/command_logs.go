package main

import (
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/nitwhiz/maas/internal/server"
	"io"
	"os"
)

type LogsCmd struct {
	Since      string `kong:"short='s',default='15m',help='Only display log entries from after this date. See docker docs.'"`
	Until      string `kong:"short='u',help='Only display log entries from up to this date. See docker docs.'"`
	Timestamps bool   `kong:"short='t',help='Log timestamps.'"`
	Follow     bool   `kong:"short='f',help='Follow Logs.'"`
}

func (c *LogsCmd) Run(ctx *Context) error {
	s, err := GetServerFromCwd()

	if err != nil {
		return err
	}

	out, err := s.GetLogs(ctx.docker, server.LogOpts{
		Timestamps: c.Timestamps,
		Follow:     c.Follow,
		Since:      c.Since,
		Until:      c.Until,
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
