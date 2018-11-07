#!/usr/bin/env bash
set -e

./xgo-check.sh --targets=linux/amd64,darwin/amd64,windows/amd64 --out=build/desktop $GOPATH/src/github.com/mysteriumnetwork/go-openvpn/examples/desktop
./xgo-check.sh --targets=ios/*,android/* --out=build/mobile $GOPATH/src/github.com/mysteriumnetwork/go-openvpn/examples/mobile