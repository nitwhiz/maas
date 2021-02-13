package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/nitwhiz/maas/internal/server"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

type ListCmd struct {
}

type ListServer struct {
	Server    server.Server
	Container types.Container
}

func getServers(docker *client.Client) ([]ListServer, error) {
	var srvs []ListServer

	labelFilters := filters.NewArgs()

	labelFilters.Add("label", "com.github.nitwhiz.maas.configPath")

	containers, err := docker.ContainerList(context.Background(), types.ContainerListOptions{
		All:     true,
		Filters: labelFilters,
	})

	if err != nil {
		return srvs, err
	}

	for _, c := range containers {
		configPath := c.Labels["com.github.nitwhiz.maas.configPath"]

		srv, err := server.FromConfig(configPath, server.ConfigOpts{
			IgnoreErrors: true,
		})

		if err == nil {
			srvs = append(srvs, ListServer{
				Server:    srv,
				Container: c,
			})
		}

	}

	return srvs, nil
}

func (c *ListCmd) Run(ctx *Context) error {
	srvs, err := getServers(ctx.docker)

	if err != nil {
		return err
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	_, _ = fmt.Fprintln(writer, "NAME\tPORT\tCONTAINER ID\tSTATUS\tDATA DIR")

	for _, srv := range srvs {
		configPath := srv.Container.Labels["com.github.nitwhiz.maas.configPath"]

		configPathSegments := strings.Split(configPath, "/")
		configPathDir := strings.Join(configPathSegments[:(len(configPathSegments)-1)], "/")

		name := configPathSegments[len(configPathSegments)-2]
		port := srv.Server.VMConfig.ExposedPort
		displayedPort := "???"
		containerId := srv.Container.ID[:12]
		status := srv.Container.Status

		if port != 0 {
			displayedPort = strconv.Itoa(port)
		}

		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\n", name, displayedPort, containerId, status, configPathDir)
	}

	_ = writer.Flush()

	return nil
}
