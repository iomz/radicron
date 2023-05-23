# radiko-auto-downloader

[![GoDoc](https://godoc.org/github.com/iomz/radiko-auto-downloader?status.svg)](https://godoc.org/github.com/iomz/radiko-auto-downloader)
[![Go Report Card](https://goreportcard.com/badge/github.com/iomz/radiko-auto-downloader)](https://goreportcard.com/report/github.com/iomz/radiko-auto-downloader)
[![Docker](https://github.com/iomz/radiko-auto-downloader/actions/workflows/docker.yml/badge.svg)](https://github.com/iomz/radiko-auto-downloader/actions/workflows/docker.yml)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

Sometimes we miss our favorite programs on [radiko](https://radiko.jp/) and they get vanished from http://radiko.jp/#!/timeshift – let's just keep them automatically downloaded locally.

# Configuration

You first need to create a configuration file (`config.toml`) to list programs to look for:

```toml
area-id = 'JP13'  # default to the "kanto" region
interval = '168h' # fetch every 7 days

[programs]
[programs.fmt.airship]
# the 2nd key `fmt` is the staion_id
# the 3rd key `airship` can be anything of your likes
title = 'GOODYEAR MUSIC AIRSHIP～シティポップ レイディオ～' # this can be a partial match
```

# Try with Docker

## Build

```
go mod vendor
docker compose build
```

## Run

```
docker-compose up -d
```

# Credit

This project is heavily based on [yyoshiki41/go-radiko](https://github.com/yyoshiki41/go-radiko) and [yyoshiki41/radigo](https://github.com/yyoshiki41/radigo), and therefore follows the [GPL-3.0 license](https://github.com/yyoshiki41/radigo/blob/main/LICENSE).
