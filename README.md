# delete-artifacts

Command line application to delete artifacts from a GitHub Workflow

[![Apache 2.0 License](https://img.shields.io/badge/License-Apache%202.0-blue)](./LICENSE)
![Go Version](https://img.shields.io/github/go-mod/go-version/jimschubert/delete-artifacts)
![Go](https://github.com/jimschubert/delete-artifacts/workflows/Build/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/jimschubert/delete-artifacts)](https://goreportcard.com/report/github.com/jimschubert/delete-artifacts)
![Docker Pulls](https://img.shields.io/docker/pulls/jimschubert/delete-artifacts)
<!-- [![codecov](https://codecov.io/gh/jimschubert/delete-artifacts/branch/master/graph/badge.svg)](https://codecov.io/gh/jimschubert/delete-artifacts) --> 

## Usage

```
Usage:
  delete-artifacts [OPTIONS]

Application Options:
  -o, --owner=   GitHub Owner/Org name [$GITHUB_ACTOR]
  -r, --repo=    GitHub Repo name [$GITHUB_REPO]
  -i, --run-id=  The workflow run id from which to delete artifacts
      --min=     Minimum size in bytes. Artifacts greater than this size will be deleted. (default: 50000000)
      --max=     Maximum size in bytes. Artifacts less than this size will be deleted
  -n, --name=    Artifact name to be deleted
  -p, --pattern= Regex pattern for matching artifact name to be deleted
      --dry-run  Dry-run that does not perform deletions
  -v, --version  Display version information

Help Options:
  -h, --help     Show this help message

```

### Examples

First, export `GITHUB_TOKEN`, then…

```
# Delete artifacts between 1B and 3MB for Run 229589570 in jimschubert/delete-artifacts-test
delete-artifacts --dry-run --owner=jimschubert --repo=delete-artifacts-test --min=1 --max=3000000 --run-id=229589570
```

```
# Delete all artifacts in jimschubert/delete-artifacts-test matching name "delete_me"
delete-artifacts --dry-run --owner=jimschubert --repo=delete-artifacts-test --name delete_me
```

```
# Delete all artifacts in jimschubert/delete-artifacts-test matching a specific regex pattern (ending in .bin), note the escape character
delete-artifacts --dry-run --owner=jimschubert --repo=delete-artifacts-test --pattern='\.bin'
```

*Remove `--dry-run` from examples to perform your delete*

## Installation

Latest binary releases are available via [GitHub Releases](https://github.com/jimschubert/delete-artifacts/releases).

## Build

Build a local distribution for evaluation using goreleaser.

```bash
goreleaser release --skip-publish --snapshot --rm-dist
```

This will create an executable application for your os/architecture under `dist`:

```
dist
├── checksums.txt
├── config.yaml
├── delete-artifacts_darwin_amd64
│   └── delete-artifacts
├── delete-artifacts_linux_386
│   └── delete-artifacts
├── delete-artifacts_linux_amd64
│   └── delete-artifacts
├── delete-artifacts_linux_arm64
│   └── delete-artifacts
├── delete-artifacts_linux_arm_6
│   └── delete-artifacts
├── delete-artifacts_v0.0.0-next_Darwin_x86_64.tar.gz
├── delete-artifacts_v0.0.0-next_Linux_arm64.tar.gz
├── delete-artifacts_v0.0.0-next_Linux_armv6.tar.gz
├── delete-artifacts_v0.0.0-next_Linux_i386.tar.gz
├── delete-artifacts_v0.0.0-next_Linux_x86_64.tar.gz
├── delete-artifacts_v0.0.0-next_Windows_x86_64.zip
├── delete-artifacts_windows_amd64
│   └── delete-artifacts.exe
└── goreleaserdocker356495573
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
