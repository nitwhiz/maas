# maas: Minecraft As A Service

An easy way to manage dockerized minecraft servers.

## Getting Started

1. Download the binary from Releases
2. Move the binary to `/usr/local/bin`
3. Change into your home directory: `cd ~`
4. Create a new vanilla 1.16.5 server listening on port 42000: `maas create --game-version 1.16.5 --port 42000 --name myserver`
5. Change into the server data directory: `cd myserver`
6. Build & start the server: `maas up`

## Usage

Disclaimer: `maas` is only tested with [itzg/docker-minecraft-server](https://github.com/itzg/docker-minecraft-server), but may work with other server containers, too.
See the [itzg/docker-minecraft-server/README.md](https://github.com/itzg/docker-minecraft-server/blob/master/README.md) on how to tune the server for your needs. 

```
Usage: maas <command>

Minecraft As A Service

Flags:
  -h, --help    Show context-sensitive help.
  -w, --working-directory=WORKING-DIRECTORY
                Set the working directory for the command

Commands:
  create      Create a new server.
  build       (Re-)create the server container.
  up          Start server.
  down        Stop server.
  list        List server containers.
  logs        Get server container logs.
  versions    Get available versions from manifest.
```

### Create new server

`maas create` generates a `maas.json` in the server data directory. All other `maas` subcommands search for it in the current working directory.
Here is a very basic example for such a file:

```json
{
  "VMConfig": {
    "Image": "itzg/minecraft-server:latest",
    "Environment": [],
    "ExposedPort": 42000
  },
  "Settings": {
    "Version": {
      "Type": "vanilla",
      "GameVersion": "1.16.5"
    }
  }
}
```

- `VMConfig` is used to configure the docker container.
  - `Image` is the docker image to use. It's not pulled automatically!
  - `Environment` is a list of `KEY=VALUE` strings used as environment in the container.
  - `ExposedPort` is the port the server listens on.
- `Settings`
  - `Version.Type` is the server type (spigot, vanilla, paper, ...).
  - `Version.GameVersion` is the minecraft game version to be used.

### Starting a server

`maas up` starts the server and builds the container, if necessary.
If you changed the `maas.json` config and there existed a server container before, it's rebuild before starting.

Optionally follow the logs immediately with `--follow`.

### Stopping a server

`maas down` stops the server. Optionally remove the server container with `--remove-container`.

### Building the server container

`maas build` (Re-)builds the server container.

### List existing maas containers

`maas list` shows all server containers on the current system and their status:

```
NAME           PORT    CONTAINER ID   STATUS                       DATA DIR
clean_piglin   32265   51bcd2ba4039   Up 7 minutes (healthy)       /srv/minecraft/clean_piglin
myserver       42000   cffe70f9f08e   Exited (0) 31 minutes ago    /srv/minecraft/myserver
lemon_mule     31476   dbb8837e8e20   Up About an hour (healthy)   /srv/minecraft/lemon_mule
```

### Show Logs

`maas logs` shows the latest server logs. See `maas logs --help` for options.
