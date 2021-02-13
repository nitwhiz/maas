package main

import (
	"github.com/alecthomas/kong"
	"github.com/docker/docker/client"
)

type CLI struct {
	WorkingDirectory WorkingDirectory `kong:"type='existingdir',short='w',help='Set the working directory for the command'"`
	Create           CreateCmd        `kong:"cmd,help='Create a new server.'"`
	Build            BuildCmd         `kong:"cmd,help='(Re-)create the server container.'"`
	Up               UpCmd            `kong:"cmd,help='Start server.'"`
	Down             DownCmd          `kong:"cmd,help='Stop server.'"`
	List             ListCmd          `kong:"cmd,help='List server containers.'"`
	Logs             LogsCmd          `kong:"cmd,help='Get server container logs.'"`
}

type Context struct {
	docker *client.Client
	cli    *CLI
}

func InitCommands(docker *client.Client) {
	cli := CLI{}

	ctx := kong.Parse(&cli,
		kong.Name("maas"),
		kong.Description("Minecraft As A Service"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
	)

	err := ctx.Run(&Context{
		docker: docker,
		cli:    &cli,
	})

	ctx.FatalIfErrorf(err)
}
