# %NAME%

Command line application written in **GO** 

[![Apache 2.0 License](https://img.shields.io/badge/License-Apache%202.0-blue)](./LICENSE)
![Go Version](https://img.shields.io/github/go-mod/go-version/%REPOSITORY%)
![Go](https://github.com/%REPOSITORY%/workflows/Build/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/%REPOSITORY%)](https://goreportcard.com/report/github.com/%REPOSITORY%)
![Docker Pulls](https://img.shields.io/docker/pulls/%REPOSITORY%)
<!-- [![codecov](https://codecov.io/gh/%REPOSITORY%/branch/master/graph/badge.svg)](https://codecov.io/gh/%REPOSITORY%) --> 

## Installation

Latest binary releases are available via [GitHub Releases](https://github.com/%REPOSITORY%/releases).

## Build

Build a local distribution for evaluation using goreleaser.

```bash
goreleaser release --skip-publish --snapshot --rm-dist
```

This will create an executable application for your os/architecture under `dist`:

```
dist
├── <app>_darwin_amd64
│   └── <app>
├── <app>_linux_386
│   └── <app>
├── <app>_linux_amd64
│   └── <app>
├── <app>_linux_arm64
│   └── <app>
├── <app>_linux_arm_6
│   └── <app>
└── <app>_windows_amd64
    └── <app>.exe
```

Build and execute locally:

* Get dependencies
```shell
go get -d ./...
```
* Build
```shell
go build cmd/main.go
```
* Run
```shell
./main
```

## License

This project is [licensed](./LICENSE) under Apache 2.0.
