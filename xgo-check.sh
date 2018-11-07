#!/usr/bin/env bash

set -e

[[ -z `which xgo` ]] || go get -u github.com/karalabe/xgo

PATH=$PATH:$GOPATH/bin xgo --image=mysteriumnetwork/xgo:1.11 $@