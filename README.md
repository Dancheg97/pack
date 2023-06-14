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

Also pack has slightly reworked API to make some operations faster. For example, flag `-q`, or `--quick`, can be used as shortcut alternative to `--noconfirm`.

Run `pack -Sh`, `pack -Rh`, `pack -Ph`, `pack -Oh` to get full description for pack commands.

![](push.png)

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

- Run pack registry:

```sh
sudo pack -O
```

### Registry

It is recommend to use container environment for security purposes, also that makes deploy much easier with `docker-compose`. Docker compose example:

```yml
# This compose can be used as a starting point to launch pack registry.
# To use pack via http (for local testing), you can provide -w option to
# syncronization and push operations on packages.

services:
  pack:
    image: fmnx.su/core/pack
    container_name: pack
    ports:
      - 80:80
    volumes:
      - ./docker:/data
    command: |
      -O
      --port 80
      --gpgdir /data/keys
      --dir /data/pkgs
    # --key /data/key.pem
    # --cert /data/cert.pem
```

To get complete working registry, you need to pass directory containing public GPG armored keys related to emails, so that pushed package signatures could be validated. You should name files containing GPG keys with accordance to email. Example:

```
/gpgdir
  -email1@example.com
  -email2@example.com
  -email3@example.com
```

After setup you can test `pack` locally with push and sync commands. Example:

```sh
git clone https://aur.archlinux.org/flutter
cd flutter
pack -Bqs
pack -Pw localhost/flutter
pack -R flutter
pack -Skyw localhost/flutter
```

If do not like docker, you can run registry with `pack -O` command on your machine.

### Libraries

Pack provides 3 libraries for go language:

- [`pacman wrapper`](pacman/README.md) - high level wrapper over pacman package manager.
- [`pack client`](pack/README.md) - top level client fucnctions to interact with package manager.
- [`pack registry`](registry/README.md) - set of interfaces required to embed pack API into existing go projects.
