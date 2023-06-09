<p align="center">
<img style="align: center; padding-left: 10px; padding-right: 10px; padding-bottom: 10px;" width="238px" height="238px" src="./logo.png" />
</p>

<h2 align="center">Pack - package manager</h2>

![Generic badge](https://img.shields.io/badge/status-alpha-red.svg)
[![Generic badge](https://img.shields.io/badge/license-gpl-orange.svg)](https://fmnx.su/core/pack/src/branch/main/LICENSE)
[![Generic badge](https://img.shields.io/badge/fmnx-repo-006db0.svg)](https://fmnx.su/core/pack)
[![Generic badge](https://img.shields.io/badge/codeberg-repo-45a3fb.svg)](https://codeberg.org/fmnx/pack)
[![Generic badge](https://img.shields.io/badge/github-repo-white.svg)](https://github.com/fmnx-io/pack)
[![Generic badge](https://img.shields.io/badge/docker-info-blue.svg)](https://fmnx.su/core/-/packages/container/pack/latest)

> **Warning!** Project is in alpha stage, API's might be changed.

Pack is utility that aims to simplify user interaction with pacman, automate some operations and provide additional functionality for software delivery.

Pack can be used to create registries that serve as regular arch package mirrors, and provides automated form of interaction them.

Also pack has slightly reworked API to make some operations faster. For example, flag `-q`, or `--quick`, can be used as shortcut alternative to `--noconfirm`. Run `pack -Sh`, `pack -Rh`, `pack -Ph`, `pack -Oh` to get full description.

### Installation

Single line installation script for all arch based distributions:

```sh
git clone https://fmnx.su/core/pack && cd pack && makepkg -sfri --needed --noconfirm
```

Alternatively, you can install pack with `go`:

```sh
go install fmnx.su/core/pack
```

### Examples

- Package installation:

```sh
pack -S nano git vim
```

- Full system upgrade:

```sh
pack -Syuq
```

- Package build (you should be in directory containing valid `PKGBUILD`):

```sh
pack -Bqs
```

- Package delivery (push package):

```sh
pack -Pf example.com/group/package
```

- Run pack registry (pack server):

```sh
sudo pack -O
```

### Registry

It is recommend to use container environment for security purposes, also that makes deploy much easier with `docker-compose`. Docker compose example:

```yml
# This compose can be used as a starting point to launch pack server.
# By default pack works properly only over https, but you can provide http
# flag to {-S --sync} command (hidden by default) for local testing.

services:
  registry:
    image: fmnx.su/core/pack
    container_name: pack
    ports:
      - 80:80
    volumes:
      - ./docker:/data
    command: |
      -O
      --port 80
      --cert /data/cert.pem
      --gpgdir /data/keys
      --dir /data/pkgs
    # --key /data/key.pem
    # --dir /data/pkgs
```

To get complete working registry, you need to pass directory with public GPG keys, so pushed package signatures could be validated with. You should name files containing GPG keys with accordance to email. Directory, provided by `--gpgdir` flag should contain keys that will be used for push validation.

After setup you can test `pack` locally with push and sync commands. Example:

```sh
pack -Bqs
pack -P --protocol http localhost
pack -S --protocol http package
```

If're not familiar with docker, you can launch server with `-O` command on your machine.

### Libraries

Pack provides 3 libraries for go language:

- [`pacman wrapper`](pacman/README.md) - high level wrapper over pacman package manager.
- [`pack client`](pack/README.md) - top level client fucnctions to interact with package manager.
- [`pack server`](server/README.md) - set of interfaces required to embed pack API into existing go projects.
