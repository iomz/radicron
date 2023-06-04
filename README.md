**_radicron_**

[![GoDoc](https://godoc.org/github.com/iomz/radicron?status.svg)](https://godoc.org/github.com/iomz/radicron)
[![Go Report Card](https://goreportcard.com/badge/github.com/iomz/radicron)](https://goreportcard.com/report/github.com/iomz/radicron)
[![codecov](https://codecov.io/gh/iomz/radicron/branch/main/graph/badge.svg?token=fjhUp7BLPB)](https://codecov.io/gh/iomz/radicron)
[![Docker](https://github.com/iomz/radicron/actions/workflows/docker.yml/badge.svg)](https://github.com/iomz/radicron/actions/workflows/docker.yml)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

Sometimes we miss our favorite programs on [radiko](https://radiko.jp/) and they get vanished from http://radiko.jp/#!/timeshift – let's just keep them automatically saved locally.

**Disclaimer**:

- Only works from an IP address within Japan (currently).
- Do not use this program for commercial purposes.

---

<!--toc:start-->

- [Configuration](#configuration)
- [Try with Docker](#try-with-docker)
  - [Build the image yourself](#build-the-image-yourself)
- [Credit](#credit)
<!--toc:end-->

# Configuration

You first need to create a configuration file (`config.yml`) to list programs to look for:

```yaml
area-id: JP13 # if unset, default to "your" region
extra-stations:
  - ALPHA-STATION # include stations not in your region
interval: 168h # fetch every 7 days (Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h" and must be positive)

rules:
  airship: # name your rule as you like
    station-id: FMT # (optional) the staion_id
    title: "GOODYEAR MUSIC AIRSHIP～シティポップ レイディオ～" # this can be a partial match
  citypop:
    keyword: "シティポップ" # search by keyword
  hiccorohee:
    pfm: "ヒコロヒー" # search by pfm
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
