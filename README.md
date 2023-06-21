<p align="center">
<img style="align: center; padding-left: 10px; padding-right: 10px; padding-bottom: 10px;" width="238px" height="238px" src="./logo.png" />
</p>

<h2 align="center">Pack - package manager</h2>

![Generic badge](https://img.shields.io/badge/status-alpha-red.svg)
[![Generic badge](https://img.shields.io/badge/license-gpl-orange.svg)](https://fmnx.su/core/pack/src/branch/main/LICENSE)
[![Generic badge](https://img.shields.io/badge/fmnx-repo-006db0.svg)](https://fmnx.su/core/pack)
[![Generic badge](https://img.shields.io/badge/codeberg-repo-45a3fb.svg)](https://codeberg.org/fmnx/pack)
[![Generic badge](https://img.shields.io/badge/github-repo-white.svg)](https://github.com/fmnx-io/pack)
[![Generic badge](https://img.shields.io/badge/arch-package-00bcd4.svg)](https://fmnx.su/core/-/packages/arch/pack)

> **Warning!** Project is in alpha stage, API's might be changed.

Pack is utility that aims to simplify user interaction with pacman, automate some operations and provide additional functionality for software delivery.

Pack can be used to push your packages to [pacman resgistries](https://fmnx.su/core/registry) and install software from them. You can test it on our [public gitea instance](https://fmnx.su/core/-/packages).

Pack API sligtly differs from pacman, to make some operations faster. For example, flag `-q`, or `--quick`, can be used as shortcut alternative to `--noconfirm`. Some other flags also might be changed, run `pack -Sh`, `pack -Rh`, `pack -Ph`, `pack -Oh` to get full description for pack commands.

If pack does not cover your needs, leave an issue in github/codeberg/fmnx repository.

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

- Push package:

```sh
pack -P example.com/group/package
```
