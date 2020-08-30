# delete-artifacts

Command line application to delete artifacts from a GitHub Workflow.

[![Apache 2.0 License](https://img.shields.io/badge/License-Apache%202.0-blue)](./LICENSE)
![Go Version](https://img.shields.io/github/go-mod/go-version/jimschubert/delete-artifacts)
![Go](https://github.com/jimschubert/delete-artifacts/workflows/Build/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/jimschubert/delete-artifacts)](https://goreportcard.com/report/github.com/jimschubert/delete-artifacts)
[![Docker Pulls](https://img.shields.io/docker/pulls/jimschubert/delete-artifacts)](https://hub.docker.com/r/jimschubert/delete-artifacts)
<!-- [![codecov](https://codecov.io/gh/jimschubert/delete-artifacts/branch/master/graph/badge.svg)](https://codecov.io/gh/jimschubert/delete-artifacts) --> 


From [GitHub Workflow's Documentation](https://docs.github.com/en/actions/configuring-and-managing-workflows/persisting-workflow-data-using-artifacts)

> Artifacts automatically expire after 90 days, but you can always reclaim used GitHub Actions storage by deleting artifacts before they expire on GitHub.

Wouldn't it be better if you could do this without manually going through every damn workflow to click delete? With `delete-artifacts`, you can:

* Delete your largest artifacts only
* Delete your artifacts for a given run
* Automate deleting your artifacts according to a schedule
* Avoid unnecessary costs for your private repositories

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
  -p, --pattern= Regex pattern (POSIX) for matching artifact name to be deleted
  -a, --active=  Consider artifacts as 'active' within this time frame, and avoid deletion. Duration formatted such as 23h59m.
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

## Via Docker

Pass the required environment variable(s) to Docker and run like so:

```bash
docker run -e GITHUB_TOKEN jimschubert/delete-artifacts:latest \
    --dry-run --owner=jimschubert --repo=delete-artifacts-test --pattern='\.bin' --min=0
```

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

## Logging

Having issues? Set `LOG_LEVEL` environment variable to one of `debug`, `info`, `warn`, or `error`.

Log outputs with messages and structured fields. For example:

```text
INFO[0000] delete-artifacts is checking the repo         owner=jimschubert repo=delete-artifacts-test
DEBU[0000] Querying artifacts across all workflows.     
DEBU[0000] Iterating artifact.                           name=artifact.bin size=1048576
DEBU[0000] MinBytes filter has matched.                  MinBytes=0
DEBU[0000] Found a set of artifacts for slated deletion.  count=1
DEBU[0000] Querying artifacts across all workflows.     
DEBU[0000] Zero artifacts remaining for query.          
DEBU[0000] Total number of artifacts to delete.          count=1
INFO[0000] Deleting artifact                             name=artifact.bin size=1048576
INFO[0000] Run complete.                    
```

## License

This project is [licensed](./LICENSE) under Apache 2.0.
