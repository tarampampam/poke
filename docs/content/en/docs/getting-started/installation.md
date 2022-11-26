---
title: Installation
slug: installation
description: "How to install Poke on Linux, Windows, macOS, or use a Docker image"
lead: "How to install Poke on Linux, Windows, macOS, or use a Docker image."
weight: 110
---

## Docker

You can use Docker to run Poke:

```bash
$ docker run --rm -v "$PWD:/fs:rw" -w /fs ghcr.io/tarampampam/poke ...
```

{{< alert icon="📦" >}}
Using the `:latest` tag for the docker image is highly discouraged because of possible backward-incompatible changes
during major upgrades. Please, use tags in `X.Y.Z` format
{{< /alert >}}

| Registry                          | Image                       |
|-----------------------------------|-----------------------------|
| [GitHub Container Registry][ghcr] | `ghcr.io/tarampampam/poke`  |
| [Docker Hub][docker-hub]          | `tarampampam/poke` (mirror) |

[ghcr]:https://github.com/users/tarampampam/packages/container/package/poke
[docker-hub]:https://hub.docker.com/r/tarampampam/poke/

The following platforms for the image are available:
<!-- docker run --rm mplatform/mquery ghcr.io/tarampampam/poke:latest -->

- `linux/amd64`
- `linux/arm64`
- `linux/arm/v6`
- `linux/arm/v7`

## Binaries Installation

### Linux

Precompiled binary is available on the [releases page][releases]. For example, let's install the
[latest version][latest-release] for **amd64** arch:

```bash
$ curl -SsL -o ./poke https://github.com/tarampampam/poke/releases/latest/download/poke-linux-amd64
$ chmod +x ./poke
```

{{< alert icon="🌍" >}}
Optionally, install the binary file globally: `sudo install -g root -o root -t /usr/local/bin -v ./poke`
{{< /alert >}}

### Debian / Ubuntu

For Debian / Ubuntu, Poke can be installed using a binary `.deb` file provided in each release:

```bash
$ curl -SsL -o ./poke-amd64.deb https://github.com/tarampampam/poke/releases/latest/download/poke-amd64.deb
$ sudo dpkg -i poke-amd64.deb
```

### macOS

{{< alert icon="⏳" text="Coming soon" />}}

### Windows

Visit the [latest releases][latest-release] page, and scroll down to the Assets section.

- Download the `poke-windows-amd64.exe` file and save it under the name `poke.exe`
- Move the executable to the desired directory
- Add this directory to the `PATH` environment variable
- Verify that you have execute permission on the file

Please consult your operating system documentation if you need help setting file permissions or modifying your
`PATH` environment variable.

## Building from sources

To build Poke from source you must:

1. Install [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
2. Install [Go](https://go.dev/doc/install) version 1.19 or later
3. Update your PATH environment variable as described in the [Go documentation](https://go.dev/doc/code#Command)

And then:

```bash
$ git clone --depth 1 --branch "vX.Y.Z" https://github.com/tarampampam/poke.git # replace vX.Y.Z with the latest version
$ cd poke
$ go generate ./...
$ go build -trimpath -ldflags "-s -w" ./cmd/poke/
```

[releases]:https://github.com/tarampampam/poke/releases
[latest-release]:https://github.com/tarampampam/poke/releases/latest
