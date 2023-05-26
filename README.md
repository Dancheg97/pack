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
[![Build Status](https://ci.fmnx.su/api/badges/core/pack/status.svg)](https://ci.fmnx.su/core/pack)

> **Warning!** Project is in alpha stage, API's are very likely to be changed.

Decentralized package manager based on `pacman` and `git`. Provides the ability to install packages using web links, aims to provide the easiest way for arch package distribution.

You can use pack to:

- install packages from compatible repos (pack-repo)
- install packages from compatible git repositories
- generate PKGBUILD's for your applications
- set up personal pack repository to provide your packages for others

Single line installation script (for all arch based distributions):

```sh
git clone https://fmnx.su/pkg/pack && cd pack && makepkg --noconfirm -sfri
```

<!--
Add emoji to help command, mb add pull req to cobra.
-->
