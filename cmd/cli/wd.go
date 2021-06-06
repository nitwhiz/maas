package main

import (
	"github.com/nitwhiz/maas/pkg/server"
	"os"
	"path"
)

type WorkingDirectory string

func (wd WorkingDirectory) AfterApply() error {
	return os.Chdir(string(wd))
}

func GetServerFromCwd() (server.Server, error) {
	var s server.Server

	cwd, err := os.Getwd()

	if err != nil {
		return s, err
	}

	maasFile := path.Join(cwd, "maas.json")

	return server.FromConfig(maasFile, server.ConfigOptions{})
}
