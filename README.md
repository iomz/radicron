![radicron](https://i.imgur.com/yVc93uh.png)

[![build status](https://github.com/iomz/radicron/workflows/build/badge.svg)](https://github.com/iomz/radicron/actions?query=workflow%3Abuild)
[![docker status](https://github.com/iomz/radicron/actions/workflows/docker.yml/badge.svg)](https://github.com/iomz/radicron/actions/workflows/docker.yml)

[![docker image size](https://ghcr-badge.egpl.dev/iomz/radicron/size)](https://github.com/iomz/radicron/pkgs/container/radicron)
[![godoc](https://godoc.org/github.com/iomz/radicron?status.svg)](https://godoc.org/github.com/iomz/radicron)
[![codecov](https://codecov.io/gh/iomz/radicron/branch/main/graph/badge.svg?token=fjhUp7BLPB)](https://codecov.io/gh/iomz/radicron)
[![go report](https://goreportcard.com/badge/github.com/iomz/radicron)](https://goreportcard.com/report/github.com/iomz/radicron)
[![license: GPL v3](https://img.shields.io/badge/license-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

Sometimes we miss our favorite shows on [radiko](https://radiko.jp/) and they get vanished from http://radiko.jp/#!/timeshift – let's just keep them automatically saved locally, from AoE.

**Disclaimer**:

- Never use this program for commercial purposes.

---

<!-- vim-markdown-toc GFM -->

- [Requirements](#requirements)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
  - [Try with Docker](#try-with-docker)
- [Build the image yourself](#build-the-image-yourself)
- [Credit](#credit)

<!-- vim-markdown-toc -->

## Requirements

radicron requires [FFmpeg](https://ffmpeg.org/download.html) to combine m3u8 chunks to a single aac file (or convert to mp3).

Make sure `ffmpeg` exists in your `$PATH`.

The [docker image](#try-with-docker) already contains all the requirements including ffmpeg.

## Installation

```bash
go install github.com/iomz/radicron/cmd/radicron@latest
```

## Configuration

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
    dow: # filter by day of the week (e.g, Mon, tue, WED)
      - fri
    pfm: "宇多丸"
    station-id: TBS
```

In addition, set `${RADIGO_HOME}` to set the download directory.

## Usage

```bash
mkdir -p ./downloads && RADIGO_HOME=./downloads radicron -c config.yml
```

### Try with Docker

By default, it mounts `./config.yml` and `./downloads` to the container.

```console
docker compose up
```

To set the ownership of the downloaded files right, run with `$UID` and `$GID` environment variables:

```console
UID=$(id -u) GID=$(id -g) docker compose up -d
```

## Build the image yourself

In case the [image](https://hub.docker.com/r/iomz/radicron/tags) is not available for your platform:

```console
docker compose build
```

## Credit

This project is heavily based on [yyoshiki41/go-radiko](https://github.com/yyoshiki41/go-radiko) and [yyoshiki41/radigo](https://github.com/yyoshiki41/radigo), and therefore follows the [GPLv3 License](https://github.com/yyoshiki41/radigo/blob/main/LICENSE).
