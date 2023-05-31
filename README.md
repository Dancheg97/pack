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

Pack is a tool that aims to simplify process of pacman package creation and allows to access them in docker-like manner from compatible pack repositories. Pack has docker version, that allows you to serve your packages and embed pack server into another go projects.

Pack can be used to:

- install packages from compatible pack repositories
- build your own pack packages and push them to repositories
- start server with to provide access with your packages to others
