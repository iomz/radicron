**_radicron_**

[![Build Status](https://github.com/iomz/radicron/workflows/build/badge.svg)](https://github.com/iomz/radicron/actions?query=workflow%3Abuild)
[![Docker](https://github.com/iomz/radicron/actions/workflows/docker.yml/badge.svg)](https://github.com/iomz/radicron/actions/workflows/docker.yml)

[![Docker Image Size](https://ghcr-badge.egpl.dev/iomz/radicron/size?label=Image%20Size)](https://github.com/iomz/radicron/pkgs/container/radicron)
[![GoDoc](https://godoc.org/github.com/iomz/radicron?status.svg)](https://godoc.org/github.com/iomz/radicron)
[![Codecov](https://codecov.io/gh/iomz/radicron/branch/main/graph/badge.svg?token=fjhUp7BLPB)](https://codecov.io/gh/iomz/radicron)
[![Go Report Card](https://goreportcard.com/badge/github.com/iomz/radicron)](https://goreportcard.com/report/github.com/iomz/radicron)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

Sometimes we miss our favorite shows on [radiko](https://radiko.jp/) and they get vanished from http://radiko.jp/#!/timeshift – let's just keep them automatically saved locally, from AoE.

**Disclaimer**:

- Never use this program for commercial purposes.

---

<!--toc:start-->

- [Installation](#installation)
- [Configuration](#configuration)
- [Try with Docker](#try-with-docker)
  - [Build the image yourself](#build-the-image-yourself)
- [Credit](#credit)

<!--toc:end-->

# Installation

```bash
go install github.com/iomz/radicron/cmd/radicron@latest
```

# Configuration

Create a configuration file (`config.yml`) to define rules for recording:

```yaml
area-id: JP13 # if unset, default to "your" region
extra-stations:
  - ALPHA-STATION # include stations not in your region
rules:
  airship: # name your rule as you like
    station-id: FMT # (optional) the staion_id, if not available by default, automatically add this station to the watch list
    title: "GOODYEAR MUSIC AIRSHIP～シティポップ レイディオ～" # this can be a partial match
  citypop:
    keyword: "シティポップ" # search by keyword (also a partial match)
    window: 48h # only within the past window from the current time
  hiccorohee:
    pfm: "ヒコロヒー" # search by pfm
  watchman:
    station-id: TBS
    pfm: "宇多丸"
    dow: # filter by day of the week (e.g, Mon, tue, WED)
      - fri
```

In addition, set `${RADIGO_HOME}` to set the download directory.

# Try with Docker

By default, it mounts `./config.yml` and `./downloads` to the container.

```console
docker compose up
```

To set the ownership of the downloaded files, run it with `$UID` and `$GID` environment variables:

```console
UID=$(id -u) GID=$(id -g) docker compose up -d
```

## Build the image yourself

In case the [image](https://hub.docker.com/r/iomz/radicron/tags) is not available for your platform:

```console
docker compose build
```

# Credit

This project is heavily based on [yyoshiki41/go-radiko](https://github.com/yyoshiki41/go-radiko) and [yyoshiki41/radigo](https://github.com/yyoshiki41/radigo), and therefore follows the [GPL-3.0 license](https://github.com/yyoshiki41/radigo/blob/main/LICENSE).
