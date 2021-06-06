# maas: Minecraft As A Service

An easy way to manage dockerized minecraft servers.
_Or just a CLI to manage [itzg/docker-minecraft-server](https://github.com/itzg/docker-minecraft-server) container._

## Getting Started

1. Download the binary from Releases
2. Move the binary to `/usr/local/bin`
3. Change into your home directory: `cd ~`
4. Create a new vanilla 1.16.5 server listening on port 42000: `maas create --game-version 1.16.5 --port 42000 --name myserver`
5. Change into the server data directory: `cd myserver`
6. Build & start the server: `maas up`

## Usage

`maas` may work with other server containers, too.
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

Use `--game-version` to set the game version the server should run:

```
> maas create --game-version 1.14

wonderful_dolphin
```

If you don't specify a name, the `create` command generates a name and outputs it. The data directory for the server will have this name, too.
(`wonderful_dolphin` in the example above)

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
  - `Image` is the docker image to use.
  - `Environment` is a list of `KEY=VALUE` strings used as environment in the container.
  - `ExposedPort` is the port the server listens on.
- `Settings`
  - `Version.Type` is the server type (spigot, vanilla, paper, ...).
  - `Version.GameVersion` is the minecraft game version to be used.

See [itzg/docker-minecraft-server/README.md](https://github.com/itzg/docker-minecraft-server/blob/master/README.md) on how to set up the `Environment`.
E.g.

```json
{
  "Environment": [
    "SNOOPER_ENABLED=false",
    "FTB_MODPACK_ID=36",
    "FTB_MODPACK_VERSION_ID=38",
    "MEMORY=3G",
    "OVERRIDE_SERVER_PROPERTIES=true",
    "MOTD=Hello World!",
    "PVP=false",
    "LEVEL_TYPE=BIOMESOP",
    "MAX_PLAYERS=30"
  ]
}
```

### Starting a server

`maas up` starts the server and builds the container, if necessary.
If you changed the `maas.json` config and there existed a server container before, it's rebuild before starting.

Optionally follow the logs immediately with `--follow`.

### Stopping a server

`maas down` stops the server. Optionally remove the server container with `--remove-container`.

### Building the server container

`maas build` (Re-)builds the server container. This will pull the specified image, if necessary.

### List existing maas containers

`maas list` shows all server containers on the current system and their status:

```
> maas list

NAME           PORT    CONTAINER ID   STATUS                       DATA DIR
clean_piglin   32265   51bcd2ba4039   Up 7 minutes (healthy)       /srv/minecraft/clean_piglin
myserver       42000   cffe70f9f08e   Exited (0) 31 minutes ago    /srv/minecraft/myserver
lemon_mule     31476   dbb8837e8e20   Up About an hour (healthy)   /srv/minecraft/lemon_mule
```

This works independent of the current working directory.

### Show Logs

`maas logs` shows the latest server logs. See `maas logs --help` for options.

## Check out versions

`maas versions` allows you to browse through all minecraft versions found in the [manifest file](https://launchermeta.mojang.com/mc/game/version_manifest_v2.json).
The downloaded manifest file is cached for 4 hours. You can refresh the cached file any time by passing `--force-download` to the `versions` command.

Show the last 10 releases:

```
> maas versions --type release

ID       TYPE      RELEASE TIME                SHA1
1.16.5   release   2021-01-14T16:05:32+00:00   436877ffaef948954053e1a78a366b8b7c204a91
1.16.4   release   2020-10-29T15:49:37+00:00   b8adadb9b21e7be96994d4485f7576377f883a0d
1.16.3   release   2020-09-10T13:42:37+00:00   1e27bd30d3bfd10e072c5ff7ea1bc2a22e11b687
1.16.2   release   2020-08-11T10:13:46+00:00   f27887ca787ae699c2284da64c2a28abda52bfe0
1.16.1   release   2020-06-24T10:31:40+00:00   2f7a599ac111edf5ba9547b665a641e1510e8580
1.16     release   2020-06-23T16:20:52+00:00   f658f2176d0bf1f6b397a23f9b4e40fa9bafbcf3
1.15     release   2019-12-09T13:13:38+00:00   bb4fc0f28e197db137d184eb963f6df918aa7ea1
1.15.2   release   2020-01-17T10:03:52+00:00   1d78c44115b99bf28a30be482979832537dc4328
1.15.1   release   2019-12-16T10:29:47+00:00   693ed0b8f1ea164fa547ca73ca5b0bd9d0693f28
1.14.4   release   2019-07-19T09:25:47+00:00   361b6f18d422c3cd7323c268201f5404b06194e4
(more)
```

Show old alpha versions:

```
> maas versions --type old_alpha
  
ID             TYPE        RELEASE TIME                SHA1
inf-20100618   old_alpha   2010-06-15T22:00:00+00:00   065ce5795aaf172080a4975cefac0d248bee7a3b
c0.30_01c      old_alpha   2009-12-21T22:00:00+00:00   0bb9bdebc3e124818fd31779a4bb394283050a02
rd-161348      old_alpha   2009-05-16T11:48:00+00:00   a937d17cca60af0a7d45d04b49a849af16b08a28
rd-160052      old_alpha   2009-05-15T22:52:00+00:00   c33dd04acfdbf34dcdfcca64db8545339ea24f02
rd-20090515    old_alpha   2009-05-14T22:00:00+00:00   1bcd01f323df5c5092e9f0967b3d310d5bc0013a
rd-132328      old_alpha   2009-05-13T21:28:00+00:00   77baa48d9cbbc6c3165c294e5bcdab2ca6903d57
rd-132211      old_alpha   2009-05-13T20:11:00+00:00   0f2a46082313d0ec67972f9f63c3fa6591f9bb85
a1.2.6         old_alpha   2010-12-02T22:00:00+00:00   d385e176aa7d3d3702bac78ad1ba906a77de13df
a1.2.5         old_alpha   2010-11-30T22:00:00+00:00   491a4961f00770bd130206c013795f35af948493
a1.2.4_01      old_alpha   2010-11-29T22:00:00+00:00   e802a257031c5b9297c971599cc2573c2efece2c
(more)

```

Search for all 1.9 release versions:

```
> maas versions --search 1.9 --type release

ID      TYPE      RELEASE TIME                SHA1
1.9     release   2016-02-29T13:49:54+00:00   fab85b386a3de3009e5944b0183ce063faa09691
1.9.4   release   2016-05-10T10:17:16+00:00   7f40b382dedcfe9eca74a5b14d615075ec34c108
1.9.3   release   2016-05-10T08:33:35+00:00   e8bab05ecee645e3c9b962f532ca7fd6c52e120e
1.9.2   release   2016-03-30T15:23:55+00:00   3cc8cee91366290508c8767e8826c6352c2b89c5
1.9.1   release   2016-03-30T13:43:07+00:00   a7c5c055718d8e7d709f3f2338b4e8f1125b5aae
```

