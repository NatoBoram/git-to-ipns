# Git to IPNS

[![Pipeline Status](https://gitlab.com/NatoBoram/git-to-ipfs/badges/master/pipeline.svg)](https://gitlab.com/NatoBoram/git-to-ipfs/commits/master)
[![Go Report Card](https://goreportcard.com/badge/gitlab.com/NatoBoram/git-to-ipfs)](https://goreportcard.com/report/gitlab.com/NatoBoram/git-to-ipfs)
[![GoDoc](https://godoc.org/gitlab.com/NatoBoram/git-to-ipfs?status.svg)](https://godoc.org/gitlab.com/NatoBoram/git-to-ipfs)

Takes a Git repository and throws it on IPNS.

## Getting started

### Dependencies

Building the program requires `go`, `npm` and `ipfs`. Don't forget to add `~/go/bin` to your `$PATH`.

This project uses [Hogan.js](https://twitter.github.io/hogan.js/) and [go.rice](https://github.com/GeertJohan/go.rice).
Their installation is scripted in [tools/dependencies.sh](tools/dependencies.sh).

```bash
./tools/dependencies.sh
```

Please note that `npm` requires root to install packages globally.

Running the program requires `ipfs-cluster-ctl`. Make sure `ipfs-cluster-service` is running, the cluster is healthy, and this peer is trusted by other peers.

### Installation

To install this project, run [tools/install.sh](tools/install.sh).

```bash
./tools/install.sh
```

It will publish the necessary files to IPFS, generate the templates and embed the web interface inside the binary.

### Running

To run this project without installing it, run [tools/run.sh](tools/run.sh).

```bash
./tools/run.sh
```

It will publish the necessary files to IPFS, generate the templates and embed the web interface inside the binary.
