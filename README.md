<p align="center">
<img style="align: center; padding-left: 10px; padding-right: 10px; padding-bottom: 10px;" width="238px" height="238px" src="./logo.png" />
</p>

<h2 align="center">Pack - package manager</h2>

[![Generic badge](https://img.shields.io/badge/LICENSE-GPL-orange.svg)](https://fmnx.io/core/pack/src/branch/main/LICENSE)
[![Generic badge](https://img.shields.io/badge/FMNX-REPO-006db0.svg)](https://fmnx.io/core/pack)
[![Generic badge](https://img.shields.io/badge/CODEBERG-REPO-45a3fb.svg)](https://codeberg.org/fmnx/pack)
[![Generic badge](https://img.shields.io/badge/GITHUB-REPO-white.svg)](https://github.com/fmnx-io/pack)
[![Generic badge](https://img.shields.io/badge/DOCKER-REGISTRY-blue.svg)](https://fmnx.io/core/-/packages/container/pack/latest)
[![Build Status](https://ci.fmnx.io/api/badges/core/pack/status.svg)](https://ci.fmnx.io/core/pack)

Decentralized package manager based on git and pacman. Accumulates power of both `git` and `pacman` to provide easier way to create arch packages and distribute them using git links.

💿 Single line installation script:

```sh
git clone https://fmnx.io/core/pack && cd pack && makepkg --noconfirm -sfri
```

🐋 You can use `pack` with docker:

```
docker run --rm -it fmnx.io/core/pack --help
```

You can use env variables to configure pack behaviour in docker:

- `PACK_ALLOW_AUR` - automatically try downloading packages from AUR if they are not found in pacman repositories
- `PACK_REMOVE_GIT_REPOS` - remove git repositories after installation
- `PACK_REMOVE_BUILT_PACKAGES` - don't cache built `.pkg.tar.zst` packages
- `PACK_DISABLE_PRETTYPRINT` - disable colors in CLI output
- `PACK_DEBUG_MODE` - watch every system call execution

Also you can modify `pack` configuration in `~/.pack/config.yml`.

---

### 🚀 Features:

- Install and update packages using git links (also allows to install AUR packages)
- Compatability with all arch based distros
- Simple to write configuration - `.pack.yml`

Configuration example for flutter project:

```yml
# Dependencies, that are required for project at runtime.
run-deps:
  - vlc
# Dependencies, that are required to build project.
build-deps:
  - flutter
  - clang
  - cmake
# Scripts, that would be executed in root directory to get build files.
scripts:
  - flutter build linux
# File mapping for resulting build files and directories from project root
# to resulting file system.
# Each file or folder will be installed as it is mapped in this file.
mapping:
  assets/logo.png: /usr/share/icons/hicolor/512x512/apps/flutter-fmnx-package.png
  flutter-fmnx-package.sh: /usr/bin/flutter-fmnx-package
  flutter_fmnx_package.desktop: /usr/share/applications/flutter-fmnx-package.desktop
  build/linux/x64/release/bundle: /usr/share/flutter-fmnx-package
```

Configuration example for go cli tool:

```yml
build-deps:
  - go
scripts:
  - go build -o example .
mapping:
  example: /usr/bin/example
```

---

### 📄 Usage

You can find all commands and description by running `pack -h`.

- `get` - get and install package from repo

```sh
pack get link.sh/owner/pkg
```

- `update` - if no arguements provided, all packages would be updated, starting from pacman packages

```sh
pack update
```

- `remove` - remove package from system

```sh
pack remove link.sh/owner/pkg
```

- `list` - list all packages installed with `pack`

```sh
pack list
```

- `gen` - generate `.pack.yml` and update `.gitignore` and `README.md`

<!--
Group pacman packages before installation.
-->
