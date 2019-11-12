#!/bin/bash

export CGO_CXXFLAGS="-v"
export CGO_CCFLAGS="-v"
gomobile bind -target=android/arm64 github.com/mysteriumnetwork/go-openvpn/test
