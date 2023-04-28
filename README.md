<p align="center">
<img style="align: center; padding-left: 10px; padding-right: 10px; padding-bottom: 10px;" width="238px" height="238px" src="./logo.png" />
</p>

<h2 align="center">pack - package manager</h2>

[![Generic badge](https://img.shields.io/badge/LICENSE-GPL-orange.svg)](https://fmnx.io/dev/pack/src/branch/main/LICENSE)
[![Generic badge](https://img.shields.io/badge/GITEA-REPO-blue.svg)](https://fmnx.io/dev/pack)
[![Generic badge](https://img.shields.io/badge/GITHUB-REPO-white.svg)](https://github.com/fmnx-io/pack)
[![Build Status](https://ci.fmnx.io/api/badges/dev/repo/status.svg)](https://ci.fmnx.io/dev/pack)

Git-based pacman-compatible package manager. Since `go` creators started reusing `git` for in go package management system, the value of decentralized systems shined from another perspective. This package manager is trying to reuse the power of both `git` and `pacman` to become new age way of arch package distribution.

---

## 🚀 Features:

- Install and update packages using git repositories
- Create packages compatible with all arch-based distros
- One simple to write file to adapt git repo - `pack.yml`

---

## 💾 Installationion

You can install `pack` on any arch-based distribution using go.

- With go

```sh
sudo pacman -S go
go install fmnx.io/dev/pack@latest
```

---

## 📄 Usage

You can find all commands and description by running `pack -h`.

- `get` - get and install package from repo

```sh
pack get link.io/owner/pkg
```

```sh
pack get link.io/owner/pkg1@v0.12 link.io/owner/pkg2@latest
```

- `update` - if no arguements provided, all packages would be updated, starting from pacman packages

```sh
pack update
```

```sh
pack update link.io/owner/pkg1 link.io/owner/pkg2
```

- `remove` - remove package from system

```sh
pack remove link.io/owner/pkg
```

- `list` - list all packages installed with `pack`

```sh
pack list
```

- `gen` - generate `pack.yml` template
