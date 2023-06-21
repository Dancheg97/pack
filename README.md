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
