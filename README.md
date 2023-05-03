<p align="center">
<img style="align: center; padding-left: 10px; padding-right: 10px; padding-bottom: 10px;" width="238px" height="238px" src="./logo.png" />
</p>

<h2 align="center">Pack - package manager</h2>

[![Generic badge](https://img.shields.io/badge/LICENSE-GPL-orange.svg)](https://fmnx.io/dev/pack/src/branch/main/LICENSE)
[![Generic badge](https://img.shields.io/badge/GITEA-REPO-blue.svg)](https://fmnx.io/dev/pack)
[![Generic badge](https://img.shields.io/badge/GITHUB-REPO-white.svg)](https://github.com/fmnx-io/pack)
[![Build Status](https://ci.fmnx.io/api/badges/dev/repo/status.svg)](https://ci.fmnx.io/dev/pack)

Distributed package manager, based on git and pacman. Accumulates power of both `git` and `pacman` to provide easier way to create arch-compatible packages and distribute them using git links.

💾 Single line installation script:

```sh
git clone https://fmnx.io/dev/pack && cd pack && makepkg --noconfirm -sfri
```

---

### 🚀 Features:

- Install and update packages using git links
- Compatability with all arch based distros
- Simple to write configuration - `.pack.yml`

Configuration [example](add_fl_tmp_link) for flutter project:

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

Configuration [example](add_fl_tmp_link) for go cli tool:

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

```sh
pack get link.sh/owner/pkg1@v0.12 link.sh/owner/pkg2@latest
```

- `update` - if no arguements provided, all packages would be updated, starting from pacman packages

```sh
pack update
```

```sh
pack update link.sh/owner/pkg1 link.sh/owner/pkg2
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
