#!/bin/sh

go clean
rm -r build
./scripts/build.sh

env GOOS=aix        GOARCH=ppc64    go build -o build/gipns-aix-ppc64
env GOOS=android    GOARCH=386      go build -o build/gipns-android-386
env GOOS=android    GOARCH=amd64    go build -o build/gipns-android-amd64
env GOOS=android    GOARCH=arm      go build -o build/gipns-android-arm
env GOOS=android    GOARCH=arm64    go build -o build/gipns-android-arm64
env GOOS=darwin     GOARCH=386      go build -o build/gipns-darwin-386
env GOOS=darwin     GOARCH=amd64    go build -o build/gipns-darwin-amd64
env GOOS=darwin     GOARCH=arm      go build -o build/gipns-darwin-arm
env GOOS=darwin     GOARCH=arm64    go build -o build/gipns-darwin-arm64
env GOOS=dragonfly  GOARCH=amd64    go build -o build/gipns-dragonfly-amd64
env GOOS=freebsd    GOARCH=386      go build -o build/gipns-freebsd-386
env GOOS=freebsd    GOARCH=amd64    go build -o build/gipns-freebsd-amd64
env GOOS=freebsd    GOARCH=arm      go build -o build/gipns-freebsd-arm
env GOOS=js         GOARCH=wasm     go build -o build/gipns-js-wasm
env GOOS=linux      GOARCH=386      go build -o build/gipns-linux-386
env GOOS=linux      GOARCH=amd64    go build -o build/gipns-linux-amd64
env GOOS=linux      GOARCH=arm      go build -o build/gipns-linux-arm
env GOOS=linux      GOARCH=arm64    go build -o build/gipns-linux-arm64
env GOOS=linux      GOARCH=mips     go build -o build/gipns-linux-mips
env GOOS=linux      GOARCH=mips64   go build -o build/gipns-linux-mips64
env GOOS=linux      GOARCH=mips64le go build -o build/gipns-linux-mips64le
env GOOS=linux      GOARCH=mipsle   go build -o build/gipns-linux-mipsle
env GOOS=linux      GOARCH=ppc64    go build -o build/gipns-linux-ppc64
env GOOS=linux      GOARCH=ppc64le  go build -o build/gipns-linux-ppc64le
env GOOS=linux      GOARCH=s390x    go build -o build/gipns-linux-s390x
env GOOS=nacl       GOARCH=386      go build -o build/gipns-nacl-386
env GOOS=nacl       GOARCH=amd64p32 go build -o build/gipns-nacl-amd64p32
env GOOS=nacl       GOARCH=arm      go build -o build/gipns-nacl-arm
env GOOS=netbsd     GOARCH=386      go build -o build/gipns-netbsd-386
env GOOS=netbsd     GOARCH=amd64    go build -o build/gipns-netbsd-amd64
env GOOS=netbsd     GOARCH=arm      go build -o build/gipns-netbsd-arm
env GOOS=openbsd    GOARCH=386      go build -o build/gipns-openbsd-386
env GOOS=openbsd    GOARCH=amd64    go build -o build/gipns-openbsd-amd64
env GOOS=openbsd    GOARCH=arm      go build -o build/gipns-openbsd-arm
env GOOS=plan9      GOARCH=386      go build -o build/gipns-plan9-386
env GOOS=plan9      GOARCH=amd64    go build -o build/gipns-plan9-amd64
env GOOS=plan9      GOARCH=arm      go build -o build/gipns-plan9-arm
env GOOS=solaris    GOARCH=amd64    go build -o build/gipns-solaris-amd64
env GOOS=windows    GOARCH=386      go build -o build/gipns-windows-386.exe
env GOOS=windows    GOARCH=amd64    go build -o build/gipns-windows-amd64.exe
env GOOS=windows    GOARCH=arm      go build -o build/gipns-windows-arm.exe