# Git to IPNS

[![Pipeline Status](https://gitlab.com/NatoBoram/git-to-ipns/badges/master/pipeline.svg)](https://gitlab.com/NatoBoram/git-to-ipns/commits/master)
[![Go Report Card](https://goreportcard.com/badge/gitlab.com/NatoBoram/git-to-ipns)](https://goreportcard.com/report/gitlab.com/NatoBoram/git-to-ipns)
[![GoDoc](https://godoc.org/gitlab.com/NatoBoram/git-to-ipns?status.svg)](https://godoc.org/gitlab.com/NatoBoram/git-to-ipns)

Takes a Git repository and throws it on IPNS.

## Getting started

### Dependencies

Building the program requires `go`, `npm` and `ipfs`. Don't forget to add `~/go/bin` to your `$PATH`.

This project uses [Hogan.js](https://twitter.github.io/hogan.js/) and [go.rice](https://github.com/GeertJohan/go.rice).
Their installation is scripted in [scripts/dependencies.sh](scripts/dependencies.sh).

```bash
./scripts/dependencies.sh
```

Please note that `npm` requires root to install packages globally.

Running the program requires `ipfs-cluster-ctl`. Make sure `ipfs-cluster-service` is running, the cluster is healthy, and this peer is trusted by other peers.

### Installation

To install this project, run [scripts/install.sh](scripts/install.sh).

```bash
./scripts/install.sh
```

It will publish the necessary files to IPFS, generate the templates and embed the web interface inside the binary.

### Running

To run this project without installing it, run [scripts/run.sh](scripts/run.sh).

```bash
./scripts/run.sh
```

It will publish the necessary files to IPFS, generate the templates and embed the web interface inside the binary.

### Building

To build this project, run [scripts/cross-compile.sh](scripts/cross-compile.sh).

```bash
./scripts/cross-compile.sh
```

It will publish the necessary files to IPFS, generate the templates, embed the web interface inside the binary, then compile the program for every single operating system and architecture supported by Go. Lots of them will fail, but the result is ready to be published.
