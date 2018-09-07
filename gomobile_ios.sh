#!/usr/bin/env bash
#Only works with patched gomobile version
CGO_LDFLAGS_ALLOW="-fobjc-arc" \
gomobile bind -target=ios/arm64 -x -v github.com/mysteriumnetwork/openvpnv3-go-bindings/openvpn3